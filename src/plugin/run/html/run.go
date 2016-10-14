package html

import (
	"framework/base/config"
	"framework/server"
	"model"
	"path/filepath"
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
	h.webPath = filepath.Join("plugin", pluginInfo.PluginUUID)
	h.localPath = filepath.Join(pluginPath, "run")
	server.ShareServerMgrInstance().RegisterStaticFile(h.webPath, h.localPath)
	return nil
}

func (h *HtmlPluginRun) Stop() error {
	server.ShareServerMgrInstance().UnRegisterStaticFile(h.webPath, h.localPath)
	return nil
}
