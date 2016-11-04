package ipc

import (
	"errors"
	"fmt"
	"framework/base/ipc/golang"
	"framework/base/timer"
	"model"
	"sync"
	"time"
)

const kWaitTime = 5 * time.Second

var pluginIPCManagerInstance *pluginIPCManager = nil
var pluginIPCManagerOnce sync.Once

type MethodCallback func(code int, response string)

type pluginInfo struct {
	pluginID int
	name     string
	ipcID    int
	ready    chan bool
}

type pluginIPCManager struct {
	ipcMgr           *golang.IPCManager
	pluginInfoMap    map[int]*pluginInfo
	pluginIDIPCIDMap map[int]int
}

func SharePluginIPCManager() *pluginIPCManager {
	pluginIPCManagerOnce.Do(func() {
		pluginIPCManagerInstance = &pluginIPCManager{}
	})
	return pluginIPCManagerInstance
}

func (p *pluginIPCManager) OnAcceptNewClient(manager *golang.IPCManager, ipcID int) {
	if _, ok := p.pluginInfoMap[ipcID]; ok {
		p.pluginInfoMap[ipcID].ready <- true
	}
}

func (p *pluginIPCManager) OnClientClose(manager *golang.IPCManager, ipcID int) {
	if _, ok := p.pluginInfoMap[ipcID]; ok {
		delete(p.pluginInfoMap, ipcID)
	}
}

func (p *pluginIPCManager) StartListener() {
	p.ipcMgr = golang.NewIPCManager()
	go func() {
		fmt.Println("pluginIPCManager StartListener")
		p.ipcMgr.StartListener()
	}()
}

func (p *pluginIPCManager) OpenPluginChannel(pluginId int) error {
	if p.pluginInfoMap == nil {
		p.pluginInfoMap = make(map[int]*pluginInfo)
	}
	if p.pluginIDIPCIDMap == nil {
		p.pluginIDIPCIDMap = make(map[int]int)
	}

	loadPluginInfo, err := model.SharePluginModel().FetchPluginByPluginID(pluginId)
	if err != nil {
		return err
	}
	ipcID := p.ipcMgr.CreateServer(loadPluginInfo.PluginName, p)
	p.pluginInfoMap[ipcID] = &pluginInfo{
		pluginID: pluginId,
		name:     loadPluginInfo.PluginName,
		ipcID:    ipcID,
		ready:    make(chan bool),
	}

	// waiting for plugin ready
	t := timer.NewOneShotTimer()
	t.Start(kWaitTime, func() {
		// waiting 5s, if no response, then plugin load failed
		delete(p.pluginIDIPCIDMap, pluginId)
		delete(p.pluginInfoMap, ipcID)
		p.pluginInfoMap[ipcID].ready <- false
	})

	isSuccess := <-p.pluginInfoMap[ipcID].ready
	t.Stop()
	if !isSuccess {
		return errors.New("load plugin failed!")
	}
	fmt.Println("open plugin channel success")
	return nil
}

func (p *pluginIPCManager) ClosePluginChannel(pluginId int) {
	if ipcID, ok := p.pluginIDIPCIDMap[pluginId]; ok {
		fmt.Println("close plugin: ", ipcID)
		// TODO:
		// p.ipcMgr.Stop(ipcID)
	}
}

func (p *pluginIPCManager) CallMethod(pluginId int, request string, callback MethodCallback) {
	if ipcID, ok := p.pluginIDIPCIDMap[pluginId]; ok {
		p.ipcMgr.CallMethod(ipcID, "HttpRequest", request, func(code int, response string) {
			callback(code, response)
		})
	} else {
		fmt.Println("plugin is not running")
	}
}
