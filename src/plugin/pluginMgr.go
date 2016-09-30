package plugin

import (
	"sync"
	"plugin/storage"
	"plugin/run/html"
	"fmt"
)

var pluginMgrInstance *pluginMgr = nil
var pluginMgrOnce sync.Once

type pluginMgr struct {

}

func SharePluginMgrInstance() *pluginMgr {
	pluginMgrOnce.Do(func() {
		pluginMgrInstance = &pluginMgr{}
	})
	return pluginMgrInstance
}

func (p *pluginMgr) AddNewPlugin(rawPluginPath string) error {
	storage := storage.NewPluginStorage(rawPluginPath)
	err := storage.Run()
	if err != nil {
		return err
	}
	return nil
}

func (p *pluginMgr) LoadPlugin(pluginId int) error  {
	var runner html.HtmlPluginRun = html.HtmlPluginRun{}
	err := runner.Run(pluginId)
	if err != nil {
		fmt.Println("LoadPlugin: ", err)
	}
	return nil
}