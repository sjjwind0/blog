package plugin

import (
	"errors"
	"fmt"
	"framework/base/ipc/golang"
	"framework/base/timer"
	"info"
	"model"
	"plugin/run"
	golangRunner "plugin/run/golang"
	h5Runner "plugin/run/html"
	"plugin/storage"
	"sync"
	"time"
)

var pluginMgrInstance *pluginMgr = nil
var pluginMgrOnce sync.Once

type MethodCallback func(code int, response string)

type pluginInfo struct {
	pluginID int
	name     string
	ipcID    int
	ready    chan bool
}

type pluginMgr struct {
	ipcMgr           *golang.IPCManager
	pluginInfoMap    map[int]*pluginInfo
	pluginIDIPCIDMap map[int]int
}

func SharePluginMgrInstance() *pluginMgr {
	pluginMgrOnce.Do(func() {
		pluginMgrInstance = &pluginMgr{}
	})
	return pluginMgrInstance
}

func (p *pluginMgr) OnAcceptNewClient(manager *golang.IPCManager, ipcID int) {
	fmt.Println("new plugin run: ", ipcID)
	if _, ok := p.pluginInfoMap[ipcID]; ok {
		p.pluginInfoMap[ipcID].ready <- true
	}
}

func (p *pluginMgr) OnClientClose(manager *golang.IPCManager, ipcID int) {
	if _, ok := p.pluginInfoMap[ipcID]; ok {
		fmt.Println("plugin uninstall: ", p.pluginInfoMap[ipcID].name)
		delete(p.pluginInfoMap, ipcID)
	}
}

func (p *pluginMgr) Initialize() {
	p.ipcMgr = golang.NewIPCManager()
	go func() {
		p.ipcMgr.StartListener()
	}()
}

func (p *pluginMgr) AddNewPlugin(rawPluginPath string) error {
	storage := storage.NewPluginStorage(rawPluginPath)
	err := storage.Run()
	if err != nil {
		return err
	}
	return nil
}

func (p *pluginMgr) LoadPlugin(pluginId int) error {
	if _, ok := p.pluginIDIPCIDMap[pluginId]; ok {
		fmt.Println("plugin is running")
		return nil
	}
	loadPluginInfo, err := model.SharePluginModel().FetchPluginByPluginID(pluginId)
	if err != nil {
		fmt.Println("fetch plugin info failed: ", err)
		return err
	}

	if p.pluginInfoMap == nil {
		p.pluginInfoMap = make(map[int]*pluginInfo)
	}
	if p.pluginIDIPCIDMap == nil {
		p.pluginIDIPCIDMap = make(map[int]int)
	}

	// start plugin ipc
	fmt.Println("start: ", loadPluginInfo.PluginName)
	ipcID := p.ipcMgr.CreateServer(loadPluginInfo.PluginName, p)
	p.pluginInfoMap[ipcID] = &pluginInfo{
		pluginID: pluginId,
		name:     loadPluginInfo.PluginName,
		ipcID:    ipcID,
		ready:    make(chan bool),
	}
	p.pluginIDIPCIDMap[pluginId] = ipcID

	var runner run.PluginRun
	switch loadPluginInfo.PluginType {
	case info.PluginType_Golang:
		runner = &golangRunner.GolangPluginRun{}
	case info.PluginType_H5:
		runner = &h5Runner.HtmlPluginRun{}
	default:
		fmt.Println("not support now")
	}
	err = runner.Run(pluginId)
	if err != nil {
		fmt.Println("LoadPlugin: ", err)
		return err
	}
	// waiting for plugin ready
	ready := p.pluginInfoMap[ipcID].ready

	t := timer.NewOneShotTimer()
	t.Start(time.Second*10, func() {
		// waiting 5s, if no response, then plugin load failed
		delete(p.pluginInfoMap, pluginId)
		ready <- false
	})

	isSuccess := <-ready
	t.Stop()
	if !isSuccess {
		return errors.New("load plugin failed!")
	}
	return nil
}

func (p *pluginMgr) CallMethod(pluginId int, request string, callback MethodCallback) {
	if ipcID, ok := p.pluginIDIPCIDMap[pluginId]; ok {
		p.ipcMgr.CallMethod(ipcID, "HttpRequest", request, func(code int, response string) {
			callback(code, response)
		})
	} else {
		p.LoadPlugin(pluginId)
		p.CallMethod(pluginId, request, callback)
	}
}
