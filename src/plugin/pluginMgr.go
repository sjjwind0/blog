package plugin

import (
	"errors"
	"fmt"
	"framework/base/ipc/golang"
	"framework/base/timer"
	"model"
	"plugin/run/html"
	"plugin/storage"
	"sync"
	"time"
)

var pluginMgrInstance *pluginMgr = nil
var pluginMgrOnce sync.Once

type pluginInfo struct {
	pluginID int
	name     string
	ipcID    int
	ready    chan bool
}

type pluginMgr struct {
	ipcMgr        *golang.IPCManager
	pluginInfoMap map[int]*pluginInfo
}

func SharePluginMgrInstance() *pluginMgr {
	pluginMgrOnce.Do(func() {
		pluginMgrInstance = &pluginMgr{}
		pluginMgrInstance.ipcMgr = golang.NewIPCManager()
		go func() {
			pluginMgrInstance.ipcMgr.StartListener()
		}()
	})
	return pluginMgrInstance
}

func (p *pluginMgr) OnAcceptNewClient(manager *golang.IPCManager, ipcID int) {
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

func (p *pluginMgr) AddNewPlugin(rawPluginPath string) error {
	storage := storage.NewPluginStorage(rawPluginPath)
	err := storage.Run()
	if err != nil {
		return err
	}
	return nil
}

func (p *pluginMgr) LoadPlugin(pluginId int) error {
	info, err := model.SharePluginModel().FetchPluginByPluginID(pluginId)
	if err != nil {
		fmt.Println("fetch plugin info failed: ", err)
		return err
	}

	if p.pluginInfoMap == nil {
		p.pluginInfoMap = make(map[int]*pluginInfo)
	}

	// start plugin ipc
	ipcID := p.ipcMgr.CreateServer(info.PluginName, p)
	p.pluginInfoMap[ipcID] = &pluginInfo{
		pluginID: pluginId,
		name:     info.PluginName,
		ipcID:    ipcID,
		ready:    make(chan bool),
	}

	var runner html.HtmlPluginRun = html.HtmlPluginRun{}
	err = runner.Run(pluginId)
	if err != nil {
		fmt.Println("LoadPlugin: ", err)
		return err
	}
	// waiting for plugin ready
	ready := p.pluginInfoMap[pluginId].ready

	t := timer.NewOneShotTimer()
	t.Start(time.Second*5, func() {
		// waiting 5s, if no response, then plugin load failed
		delete(p.pluginInfoMap, pluginId)
		ready <- false
	})

	isSuccess := <-ready
	if !isSuccess {
		return errors.New("load plugin failed!")
	}
	return nil
}
