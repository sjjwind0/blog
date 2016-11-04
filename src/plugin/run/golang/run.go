package golang

import (
	"fmt"
	"framework/base/config"
	"framework/base/shell"
	"model"
	"path/filepath"
	"plugin/ipc"
)

type golangPluginRun struct {
	pluginId int
}

func NewGolangPluginRunner(pluginId int) *golangPluginRun {
	return &golangPluginRun{pluginId: pluginId}
}

func (p *golangPluginRun) Run() error {
	err := ipc.SharePluginIPCManager().OpenPluginChannel(p.pluginId)
	if err != nil {
		fmt.Println("open plugin channel error: ", err)
		return err
	}

	loadPluginInfo, err := model.SharePluginModel().FetchPluginByPluginID(p.pluginId)
	if err != nil {
		p.Stop()
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
			p.Stop()
		}
	}()
	return nil
}

func (p *golangPluginRun) Stop() error {
	ipc.SharePluginIPCManager().ClosePluginChannel(p.pluginId)
	return nil
}
