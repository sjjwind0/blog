package golang

import (
	"framework/base/config"
	"framework/server"
	"model"
	"os/exec"
	"path/filepath"
)

type GolangPluginRun struct {
}

func (h *GolangPluginRun) Run(pluginId int) error {
	// pluginPath := config.GetDefaultConfigJsonReader().GetString("storage.file.plugin")
	// pluginInfo, err := model.SharePluginModel().FetchPluginByPluginID(pluginId)
	// if err != nil {
	// 	return err
	// }
	// FIXME: test, hardcode pluginPath
	pluginPath := "/home/wind/data/go-plugin-demo"
	binaryPath := filepath.Join(pluginPath, "plugin")
	exec.Command(binaryPath).Start()
	return nil
}

func (h *GolangPluginRun) Stop() error {
	server.ShareServerMgrInstance().UnRegisterStaticFile(h.webPath, h.localPath)
	return nil
}
