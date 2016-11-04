package html

import (
	"framework/base/config"
	"framework/server"
	"model"
	"path/filepath"
)

type htmlPluginRun struct {
	webPath   string
	localPath string
	pluginId  int
}

func NewHtmlPluginRunner(pluginId int) *htmlPluginRun {
	return &htmlPluginRun{pluginId: pluginId}
}

func (h *htmlPluginRun) Run() error {
	pluginPath := config.GetDefaultConfigJsonReader().GetString("storage.file.plugin")
	pluginInfo, err := model.SharePluginModel().FetchPluginByPluginID(h.pluginId)
	if err != nil {
		return err
	}
	pluginPath = filepath.Join(pluginPath, pluginInfo.PluginUUID)
	h.webPath = filepath.Join("plugin", pluginInfo.PluginUUID)
	h.localPath = filepath.Join(pluginPath, "run")
	server.ShareServerMgrInstance().RegisterStaticFile(h.webPath, h.localPath)
	return nil
}

func (h *htmlPluginRun) Stop() error {
	server.ShareServerMgrInstance().UnRegisterStaticFile(h.webPath, h.localPath)
	return nil
}
