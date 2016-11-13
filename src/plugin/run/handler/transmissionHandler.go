package handler

import (
	"errors"
	"fmt"
	"framework"
	"framework/base/config"
	"framework/response"
	"framework/server"
	"model"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type TransmissionRequestHandler struct {
	staticFileMap map[string]string
}

func (t *TransmissionRequestHandler) Register(pluginId int) error {
	pluginPath := config.GetDefaultConfigJsonReader().GetString("storage.file.plugin")
	pluginInfo, err := model.SharePluginModel().FetchPluginByPluginID(pluginId)
	if err != nil {
		return err
	}
	pluginPath = filepath.Join(pluginPath, pluginInfo.PluginUUID)
	webPath := filepath.Join("plugin", strconv.Itoa(pluginInfo.PluginID))
	localPath := filepath.Join(pluginPath, "code")

	if t.staticFileMap == nil {
		t.staticFileMap = make(map[string]string)
	}
	if _, ok := t.staticFileMap[webPath]; ok {
		fmt.Println("static file has beed registered!")
		return errors.New("static file has beed registered!")
	}
	walkPath := localPath
	filepath.Walk(walkPath, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			rel, _ := filepath.Rel(localPath, path)
			webFilePath := rel
			t.staticFileMap[webFilePath], err = filepath.Abs(path)
		}
		return nil
	})
	return nil
}

func (t *TransmissionRequestHandler) HandlePluginRequest(pluginId int, w http.ResponseWriter, r *http.Request) {
	fmt.Println("TransmissionRequestHandler HandlePluginRequest")
	currentPath := r.URL.Path
	if currentPath[0] == '/' {
		currentPath = currentPath[1:]
	}
	if local, ok := t.staticFileMap[currentPath]; ok {
		ext := filepath.Ext(local)
		contentType := server.QueryContentTypeByExt(strings.ToLower(ext))
		file, err := os.Open(local)
		if err != nil {
			response.JsonResponseWithMsg(w, framework.ErrorNoSuchFileOrDirectory, err.Error())
			return
		}
		defer file.Close()
		fileInfo, err := os.Stat(local)
		if err != nil {
			response.JsonResponseWithMsg(w, framework.ErrorNoSuchFileOrDirectory, err.Error())
			return
		}
		content := make([]byte, fileInfo.Size())
		file.Read(content)
		w.Header().Set("Accept", "*/*")
		w.Header().Set("Content-Length", fmt.Sprintf("%d", fileInfo.Size()))
		w.Header().Set("Content-Type", contentType)
		w.Write(content)
	} else {
		response.JsonResponseWithMsg(w, framework.ErrorFileNotExist, "no such file")
	}
}
