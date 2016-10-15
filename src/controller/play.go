package controller

import (
	"fmt"
	"framework"
	"framework/base/config"
	"framework/base/json"
	"framework/response"
	"framework/server"
	"html/template"
	"info"
	"model"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

type playRender struct {
	PluginID          int
	PluginName        string
	PluginDescription string
}

type pluginListRender struct {
	PluginList []*playRender
	Host       *hostRender
}

type PlayController struct {
}

func NewPlayController() *PlayController {
	return &PlayController{}
}

func (a *PlayController) Path() interface{} {
	return []string{"/play", "/big_cover", "/small_cover"}
}

func (a *PlayController) readFileContent(path string) *[]byte {
	fileInfo, err := os.Stat(path)
	if err == nil {
		content := make([]byte, fileInfo.Size())
		file, _ := os.Open(path)
		defer file.Close()
		file.Read(content)
		return &content
	}
	return nil
}

func (a *PlayController) readRes(w http.ResponseWriter, path string) {
	ext := filepath.Ext(path)
	imgContent := a.readFileContent(path)
	if imgContent == nil {
		return
	}
	w.Header().Set("Accept", "*/*")
	w.Header().Set("Content-Length", strconv.Itoa(len(*imgContent)))
	w.Header().Set("Content-Type", server.QueryContentTypeByExt(ext))
	w.Write(*imgContent)
}

func (a *PlayController) HandlerRequest(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/play" {
		playRenderList := &pluginListRender{}
		allPlugins, err := model.SharePluginModel().FetchAllPlugin()
		if err != nil {
			fmt.Println("get plugin failed")
		} else {
			pluginRootPath := config.GetDefaultConfigJsonReader().GetString("storage.file.plugin")
			for iter := allPlugins.Front(); iter != nil; iter = iter.Next() {
				info := iter.Value.(info.PluginInfo)
				pluginInfoPath := filepath.Join(pluginRootPath, info.PluginUUID, "plugin.info")
				description := json.NewJsonReaderFromFile(pluginInfoPath).GetString("description")
				playRender := &playRender{}
				playRender.PluginName = info.PluginName
				playRender.PluginID = info.PluginID
				playRender.PluginDescription = description
				playRenderList.PluginList = append(playRenderList.PluginList, playRender)
			}
		}
		playRenderList.Host = buildHostRender()
		t, err := template.ParseFiles("./src/view/html/play.html")
		if err != nil {
			fmt.Println("parse file error: ", err.Error())
		}
		t.Execute(w, &playRenderList)
	} else if r.URL.Path == "/big_cover" || r.URL.Path == "/small_cover" {
		r.ParseForm()
		id, err := strconv.Atoi(r.Form.Get("id"))
		if err != nil {
			response.JsonResponseWithMsg(w, framework.ErrorParamError, err.Error())
			return
		}
		uuid, err := model.SharePluginModel().GetPluginUUIDByPluginID(id)
		if err != nil || uuid == "" {
			response.JsonResponseWithMsg(w, framework.ErrorParamError, err.Error())
			return
		}
		blogPath := config.GetDefaultConfigJsonReader().GetString("storage.file.plugin")
		var imgPath string = ""
		if r.URL.Path == "/big_cover" {
			imgPath = filepath.Join(blogPath, uuid, "big_cover.jpg")
		} else {
			imgPath = filepath.Join(blogPath, uuid, "small_cover.jpg")
		}
		a.readRes(w, imgPath)
	}
}
