package build

import (
	"errors"
	"fmt"
	"framework/base/config"
	"info"
	"model"
	"path/filepath"
	"plugin/build/golang"
	"plugin/build/html"
	"plugin/build/step"
)

type ProgressCallback func(info string, err string, isComplete bool)

type Builder interface {
	BuildStep() []step.BuildStep
}

type BuilderMgr struct {
	builder Builder
}

func newBuildMgrFromBuilder(builder Builder) *BuilderMgr {
	return &BuilderMgr{builder: builder}
}

func NewBuilderMgr(pluginId int) (*BuilderMgr, error) {
	pluginInfo, err := model.SharePluginModel().FetchPluginByPluginID(pluginId)
	if err != nil {
		fmt.Println("fetch plugin info error: ", err)
		return nil, err
	}
	pluginPath := config.GetDefaultConfigJsonReader().GetString("storage.file.plugin")
	pluginPath = filepath.Join(pluginPath, pluginInfo.PluginUUID, "code")

	var builder Builder = nil
	switch pluginInfo.PluginType {
	case info.PluginType_None:
		fmt.Println("none plugin type, do nothing")
	case info.PluginType_H5:
		builder = html.NewHtmlBuilder()
	case info.PluginType_CPP:
	case info.PluginType_Java:
	case info.PluginType_Golang:
		builder = golang.NewBuilder(pluginPath)
	case info.PluginType_Node:
	case info.PluginType_Python:
	default:
		fmt.Println("unsport plugin type")
		return nil, errors.New("unsport plugin type")
	}
	mgr := newBuildMgrFromBuilder(builder)
	return mgr, nil
}

func (b *BuilderMgr) Run(callback ProgressCallback) {
	go func() {
		var outputStr string = ""
		var errorStr string = ""
		buildStepList := b.builder.BuildStep()
		for index, step := range buildStepList {
			fmt.Println(index, ". build")
			description := fmt.Sprintf("%d. %s\n", index+1, step.Description())
			callback(description, "", false)
			outString, errString, err := step.Run()
			if err != nil {
				errString += err.Error()
				fmt.Println("err: ", err)
				return
			}
			outputStr += outString
			errorStr += errString
			callback(outString, errString, false)
		}
		fmt.Println("complete")
		callback(outputStr, errorStr, true)
	}()
}
