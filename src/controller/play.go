package controller

import (
	//"fmt"
	//"html/template"
	"net/http"
	//"plugin"
)

type pluginRender struct {
	PluginTitle       string
	PluginDescription string
	PluginCoverURL    string
	PluginURL         string
	PluginDownloadURL string
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
	//pluginsInfo := plugin.GetDefaultPluginManager().GetAllPluginInfo()
	//pluginList := &pluginListRender{}
	//for _, plugin := range pluginsInfo {
	//	pluginRender := &pluginRender{}
	//	pluginRender.PluginTitle = plugin.GetPluginName()
	//	pluginRender.PluginDescription = plugin.GetPluginDescription()
	//	pluginRender.PluginCoverURL = plugin.GetPluginCoverURL()
	//	pluginRender.PluginURL = plugin.GetPluginDisplayPath()
	//	pluginRender.PluginDownloadURL = plugin.GetPluginDownloadURL()
	//	pluginList.PluginList = append(pluginList.PluginList, pluginRender)
	//}
	//pluginList.Host = buildHostRender()
	//t, err := template.ParseFiles("./src/view/html/play.html")
	//if err != nil {
	//	fmt.Println("parse file error: ", err.Error())
	//}
	//t.Execute(w, &pluginList)
}
