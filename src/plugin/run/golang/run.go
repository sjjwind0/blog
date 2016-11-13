package golang

import (
	"fmt"
	"framework/base/config"
	"framework/base/shell"
	"model"
	"os"
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
	progress *os.Process
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

	p.progress, err = shell.RunShellAsync(func(output string, errOutput string) {
		fmt.Println("shell out: ", output)
		fmt.Println("shell err: ", errOutput)
	}, binaryDir, "./plugin")
	return ipc.SharePluginIPCManager().WaitingForPluginStart(ipcId)
}

func (p *golangPluginRun) Stop() error {
	p.stopType = StopBySelf
	ipc.SharePluginIPCManager().ClosePluginChannel(p.pluginId)
	p.progress.Kill()
	return nil
}
