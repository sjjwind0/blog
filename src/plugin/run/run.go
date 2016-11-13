package run

import (
	"fmt"
	"info"
	"net/http"
	"plugin/run/golang"
	"plugin/run/html"
)

type PluginRun interface {
	Run() error
	Stop() error
	HandlePluginRequest(pluginId int, w http.ResponseWriter, r *http.Request)
}

func NewPluginRunner(pluginType int, pluginId int) PluginRun {
	switch pluginType {
	case info.PluginType_Golang:
		return golang.NewGolangPluginRunner(pluginId)
	case info.PluginType_H5:
		return html.NewHtmlPluginRunner(pluginId)
	default:
		fmt.Println("not support now")
		return nil
	}
}
