package html

import (
	"framework/server"
	"path/filepath"
	"framework/base/config"
	"model"
	"fmt"
)

type HtmlPluginRun struct {
	webPath   string
	localPath string
}

func (h *HtmlPluginRun) Run(pluginId int) error {
	pluginPath := config.GetDefaultConfigJsonReader().GetString("storage.file.plugin")
	pluginInfo, err := model.SharePluginModel().FetchPluginByPluginID(pluginId)
	if err != nil {
		return err
	}
	pluginPath = filepath.Join(pluginPath, pluginInfo.PluginUUID)
	configPath := filepath.Join(pluginPath, "plugin.conf")
	fmt.Println("configPath: ", configPath)
	h.webPath = filepath.Join("plugin", pluginInfo.PluginUUID)
	h.localPath = filepath.Join(pluginPath, "run")
	fmt.Println("web: ", h.webPath)
	fmt.Println("local: ", h.localPath)
	server.ShareServerMgrInstance().RegisterStaticFile(h.webPath, h.localPath)
	return nil
}

func (h *HtmlPluginRun) Stop() error {
	server.ShareServerMgrInstance().UnRegisterStaticFile(h.webPath, h.localPath)
	return nil
}
