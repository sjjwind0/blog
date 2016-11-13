package ipc

import (
	"errors"
	"fmt"
	"framework"
	"framework/base/ipc/golang"
	"framework/base/timer"
	"model"
	"sync"
	"time"
)

const kWaitTime = 5 * time.Second

var pluginIPCManagerInstance *pluginIPCManager = nil
var pluginIPCManagerOnce sync.Once

type ReadyCallback func(bool)

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
	delegate         PluginDelegate
	callbackList     map[int][]ReadyCallback
	isStartingMap    map[int]bool
}

func SharePluginIPCManager() *pluginIPCManager {
	pluginIPCManagerOnce.Do(func() {
		pluginIPCManagerInstance = &pluginIPCManager{}
	})
	return pluginIPCManagerInstance
}

func (p *pluginIPCManager) OnAcceptNewClient(manager *golang.IPCManager, ipcID int) {
	fmt.Println("OnAcceptNewClient, ipcID: ", ipcID)
	if _, ok := p.pluginInfoMap[ipcID]; ok {
		fmt.Println("pluginIPCManager start success")
		p.pluginInfoMap[ipcID].ready <- true
	} else {
		fmt.Println("plugin start timeout, need to close")
		p.ipcMgr.StopClient(ipcID)
	}
}

func (p *pluginIPCManager) OnClientClose(manager *golang.IPCManager, ipcID int) {
	if v, ok := p.pluginInfoMap[ipcID]; ok {
		if p.delegate != nil {
			p.delegate.OnPluginShutdown(v.pluginID)
		}
		delete(p.pluginInfoMap, ipcID)
	}
}

func (p *pluginIPCManager) OnPluginReady(pluginId int) {
	fmt.Println("plugin: ", pluginId, " start success")
	p.isStartingMap[pluginId] = false
	if callbackList, ok := p.callbackList[pluginId]; ok {
		for _, callback := range callbackList {
			callback(true)
		}
		delete(p.callbackList, pluginId)
	}
}

func (p *pluginIPCManager) OnPluginStartFailed(pluginId int) {
	fmt.Println("plugin: ", pluginId, " start failed")
	if callbackList, ok := p.callbackList[pluginId]; ok {
		for _, callback := range callbackList {
			callback(false)
		}
		delete(p.callbackList, pluginId)
	}
}

func (p *pluginIPCManager) StartListener() {
	p.ipcMgr = golang.NewIPCManager()
	go func() {
		fmt.Println("pluginIPCManager StartListener")
		p.ipcMgr.StartListener()
	}()
}

func (p *pluginIPCManager) SetDelegate(delegate PluginDelegate) {
	p.delegate = delegate
}

func (p *pluginIPCManager) OpenPluginChannel(pluginId int) (int, error) {
	if p.pluginInfoMap == nil {
		p.pluginInfoMap = make(map[int]*pluginInfo)
	}
	if p.pluginIDIPCIDMap == nil {
		p.pluginIDIPCIDMap = make(map[int]int)
	}

	loadPluginInfo, err := model.SharePluginModel().FetchPluginByPluginID(pluginId)
	if err != nil {
		return 0, err
	}
	ipcID := p.ipcMgr.CreateServer(loadPluginInfo.PluginName, p)
	p.pluginInfoMap[ipcID] = &pluginInfo{
		pluginID: pluginId,
		name:     loadPluginInfo.PluginName,
		ipcID:    ipcID,
		ready:    make(chan bool),
	}
	p.pluginIDIPCIDMap[pluginId] = ipcID

	if p.isStartingMap == nil {
		p.isStartingMap = make(map[int]bool)
	}
	p.isStartingMap[pluginId] = true
	return ipcID, nil
}

func (p *pluginIPCManager) WaitingForPluginStart(ipcID int) error {
	// waiting for plugin ready
	pluginId := p.pluginInfoMap[ipcID].pluginID
	t := timer.NewOneShotTimer()
	t.Start(kWaitTime, func() {
		// waiting 5s, if no response, then plugin load failed
		p.pluginInfoMap[ipcID].ready <- false
		if p.delegate != nil {
			p.OnPluginStartFailed(pluginId)
		}
		delete(p.pluginIDIPCIDMap, pluginId)
		delete(p.pluginInfoMap, ipcID)
	})

	isSuccess := <-p.pluginInfoMap[ipcID].ready
	t.Stop()
	if !isSuccess {
		if p.delegate != nil {
			p.OnPluginStartFailed(pluginId)
		}
		return errors.New("load plugin failed!")
	}
	if p.delegate != nil {
		p.OnPluginReady(pluginId)
	}
	fmt.Println("open plugin channel success")
	return nil
}

func (p *pluginIPCManager) ClosePluginChannel(pluginId int) {
	if ipcID, ok := p.pluginIDIPCIDMap[pluginId]; ok {
		fmt.Println("close plugin: ", ipcID)
		p.ipcMgr.StopClient(ipcID)
		delete(p.pluginIDIPCIDMap, pluginId)
		delete(p.pluginInfoMap, ipcID)
		delete(p.isStartingMap, pluginId)
		delete(p.callbackList, pluginId)
	}
}

func (p *pluginIPCManager) CallMethod(pluginId int, request string, callback MethodCallback) {
	if _, ok := p.pluginIDIPCIDMap[pluginId]; !ok {
		if p.callbackList == nil {
			p.callbackList = make(map[int][]ReadyCallback)
		}
		p.callbackList[pluginId] = append(p.callbackList[pluginId], func(isSucccess bool) {
			if isSucccess {
				p.CallMethodInternal(pluginId, request, callback)
			} else {
				callback(framework.ErrorPluginLoadError, "")
			}
		})
		p.delegate.OnPluginNeedStart(pluginId)
	} else {
		if v, ok := p.isStartingMap[pluginId]; ok && v == false {
			p.CallMethodInternal(pluginId, request, callback)
		} else {
			p.callbackList[pluginId] = append(p.callbackList[pluginId], func(isSucccess bool) {
				if isSucccess {
					p.CallMethodInternal(pluginId, request, callback)
				} else {
					callback(framework.ErrorPluginLoadError, "")
				}
			})
		}
	}
}

func (p *pluginIPCManager) CallMethodInternal(pluginId int, request string, callback MethodCallback) {
	if ipcID, ok := p.pluginIDIPCIDMap[pluginId]; ok {
		p.ipcMgr.CallMethod(ipcID, "HttpRequest", request, func(code int, response string) {
			callback(code, response)
		})
	} else {
		fmt.Println("plugin is not running")
	}
}
