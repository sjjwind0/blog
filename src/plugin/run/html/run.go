package html

import (
	"plugin/run/handler"
)

type htmlPluginRun struct {
	handler.TransmissionRequestHandler
}

func NewHtmlPluginRunner(pluginId int) *htmlPluginRun {
	runner := new(htmlPluginRun)
	runner.TransmissionRequestHandler.Register(pluginId)
	return runner
}

func (h *htmlPluginRun) Run() error {
	return nil
}

func (h *htmlPluginRun) Stop() error {
	// server.ShareServerMgrInstance().UnRegisterStaticFile(h.webPath, h.localPath)
	return nil
}
