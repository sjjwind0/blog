package plugin

import (
	"errors"
	"fmt"
	"framework"
	"framework/response"
	"model"
	"net/http"
	"plugin/build"
	"plugin/ipc"
	"plugin/run"
	"plugin/storage"
	"sync"
)

var pluginMgrInstance *pluginMgr = nil
var pluginMgrOnce sync.Once

type pluginMgr struct {
	pluginRunnerMap map[int]run.PluginRun
}

func SharePluginMgrInstance() *pluginMgr {
	pluginMgrOnce.Do(func() {
		pluginMgrInstance = &pluginMgr{}
	})
	return pluginMgrInstance
}

func (p *pluginMgr) OnPluginNeedStart(pluginId int) {
	p.LoadPlugin(pluginId)
}

func (p *pluginMgr) OnPluginShutdown(pluginId int) {
	if _, ok := p.pluginRunnerMap[pluginId]; ok {
		delete(p.pluginRunnerMap, pluginId)
	}
}

func (p *pluginMgr) Initialize() {
	ipc.SharePluginIPCManager().SetDelegate(p)
	ipc.SharePluginIPCManager().StartListener()
}

func (p *pluginMgr) AddNewPlugin(rawPluginPath string, callback build.ProgressCallback) error {
	storage := storage.NewPluginStorage(rawPluginPath)
	err := storage.Run()
	if err != nil {
		return err
	}
	pluginId := storage.GetPluginID()
	buildMgr, err := build.NewBuilderMgr(pluginId)
	if err != nil {
		fmt.Println("get build failed: ", err)
		return err
	}
	buildMgr.Run(callback)
	return nil
}

func (p *pluginMgr) LoadPlugin(pluginId int) error {
	fmt.Println("pluginMgr LoadPlugin: ", pluginId)
	if p.pluginRunnerMap == nil {
		p.pluginRunnerMap = make(map[int]run.PluginRun)
	}
	if _, ok := p.pluginRunnerMap[pluginId]; ok {
		fmt.Println("plugin is running")
		return nil
	}
	loadPluginInfo, err := model.SharePluginModel().FetchPluginByPluginID(pluginId)
	if err != nil {
		fmt.Println("fetch plugin info failed: ", err)
		return err
	}

	runner := run.NewPluginRunner(loadPluginInfo.PluginType, pluginId)
	if err != nil {
		return err
	}
	p.pluginRunnerMap[pluginId] = runner
	return runner.Run()
}

func (p *pluginMgr) StopPlugin(pluginId int) error {
	fmt.Println("pluginMgr StopPlugin: ", pluginId)
	if runner, ok := p.pluginRunnerMap[pluginId]; ok {
		runner.Stop()
		return nil
	}
	fmt.Println("plugin is not running")
	return errors.New("plugin is not runner")
}

func (p *pluginMgr) HandleRequest(pluginId int, w http.ResponseWriter, r *http.Request) {
	if runner, ok := p.pluginRunnerMap[pluginId]; ok {
		runner.HandlePluginRequest(pluginId, w, r)
	} else {
		if p.LoadPlugin(pluginId) == nil {
			fmt.Println("load success: ", pluginId)
			p.HandleRequest(pluginId, w, r)
		} else {
			response.JsonResponseWithMsg(w, framework.ErrorPluginNotExist, "plugin not exist")
		}
	}
}
