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
	h.TransmissionRequestHandler.UnRegister()
	return nil
}
