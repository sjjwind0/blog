package golang

import (
	"fmt"
	"framework/base/config"
	"framework/base/shell"
	"model"
	"path/filepath"
)

type GolangPluginRun struct {
}

func (h *GolangPluginRun) Run(pluginId int) error {
	loadPluginInfo, err := model.SharePluginModel().FetchPluginByPluginID(pluginId)
	if err != nil {
		return err
	}

	pluginRootPath := config.GetDefaultConfigJsonReader().GetString("storage.file.plugin")
	pluginRootPath = filepath.Join(pluginRootPath, loadPluginInfo.PluginUUID)
	binaryDir := filepath.Join(pluginRootPath, "run")
	fmt.Println("binaryDir: ", binaryDir)
	go func() {
		output, errOutput, err := shell.RunShell(binaryDir, "./plugin")
		if err == nil {
			fmt.Println("shell out: ", output)
			fmt.Println("shell err: ", errOutput)
		} else {
			fmt.Println("run shell error: ", err)
		}
	}()
	return nil
}

func (h *GolangPluginRun) Stop() error {
	return nil
}
