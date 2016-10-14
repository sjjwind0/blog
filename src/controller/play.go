package controller

import (
	"fmt"
	"html/template"
	"info"
	"model"
	"net/http"
)

type pluginRender struct {
	PluginID          int
	PluginName        string
	PluginDescription string
}

type pluginListRender struct {
	PluginList []*pluginRender
	Host       *hostRender
}

type PlayController struct {
}

func NewPlayController() *PlayController {
	return &PlayController{}
}

func (a *PlayController) Path() interface{} {
	return "/play"
}

func (a *PlayController) HandlerRequest(w http.ResponseWriter, r *http.Request) {
	pluginRenderList := &pluginListRender{}
	allPlugins, err := model.SharePluginModel().FetchAllPlugin()
	if err != nil {
		fmt.Println("get plugin failed")
	} else {
		for iter := allPlugins.Front(); iter != nil; iter = iter.Next() {
			info := iter.Value.(info.PluginInfo)
			pluginRender := &pluginRender{}
			pluginRender.PluginName = info.PluginName
			pluginRender.PluginID = info.PluginID
			pluginRender.PluginDescription = "test description"
			pluginRenderList.PluginList = append(pluginRenderList.PluginList, pluginRender)
		}
	}
	pluginRenderList.Host = buildHostRender()
	t, err := template.ParseFiles("./src/view/html/play.html")
	if err != nil {
		fmt.Println("parse file error: ", err.Error())
	}
	t.Execute(w, &pluginRenderList)
}
