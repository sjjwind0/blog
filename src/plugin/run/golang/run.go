package golang

import (
	"fmt"
	"framework/base/config"
	"framework/base/shell"
	"model"
	"path/filepath"
	"plugin/ipc"
	"plugin/run/handler"
)

const (
	StopByKnown = iota
	StopBySelf  = iota
)

type golangPluginRun struct {
	pluginId int
	stopType int
	handler.IPCRequestHandler
}

func NewGolangPluginRunner(pluginId int) *golangPluginRun {
	return &golangPluginRun{pluginId: pluginId, stopType: StopByKnown}
}

func (p *golangPluginRun) Run() error {
	loadPluginInfo, err := model.SharePluginModel().FetchPluginByPluginID(p.pluginId)
	if err != nil {
		p.Stop()
		fmt.Println("fetch plugin error: ", err)
		return err
	}

	ipcId, err := ipc.SharePluginIPCManager().OpenPluginChannel(p.pluginId)
	if err != nil {
		fmt.Println("open plugin channel error: ", err)
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
		if p.stopType == StopBySelf {
			fmt.Println("stop success")
		} else {
			fmt.Println("plugin was killed by some reason")
			p.Stop()
		}
	}()
	return ipc.SharePluginIPCManager().WaitingForPluginStart(ipcId)
}

func (p *golangPluginRun) Stop() error {
	p.stopType = StopBySelf
	ipc.SharePluginIPCManager().ClosePluginChannel(p.pluginId)
	return nil
}
