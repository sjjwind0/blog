package storage

import (
	"errors"
	"fmt"
	"framework/base/archive"
	"framework/base/config"
	"framework/base/json"
	"info"
	"model"
	"os"
	"path/filepath"
)

type pluginStorage struct {
	rawPluginPath string
}

func NewPluginStorage(path string) *pluginStorage {
	return &pluginStorage{rawPluginPath: path}
}

/*
plugin 文件结构
		- run
			二进制可执行文件
			config
			view
				html
				js
				css

		- plugin.info
			{
				"name": "chess",
				"uuid": "uuid",
				"version": "1.0.0",
				"type": "golang|html|python|node|C++",
				"description": "plugin的文件描述",
				"command": {
					"work_dir": "./run",
					"env": "/bin/python",
					"file": "main.py"
				},
			}

		- big_cover.jpg

		- small_cover.jpg

		- raw.zip

		- code.zip
*/
func (p *pluginStorage) Run() error {
	c := config.GetDefaultConfigJsonReader()
	tmpRawPath := c.GetString("storage.file.tmp")

	err := archive.UnZipToPath(p.rawPluginPath, tmpRawPath)
	if err != nil {
		return err
	}

	// 1. save raw.zip
	rawPath := c.GetString("storage.file.raw")
	os.MkdirAll(rawPath, os.ModePerm)

	// 2. save plugin
	// 2.1 read config
	tmpConfigPath := filepath.Join(tmpRawPath, "raw", "plugin.info")
	fmt.Println("tmpConfigPath: ", tmpConfigPath)
	jsonReader := json.NewJsonReaderFromFile(tmpConfigPath)
	uuid := jsonReader.GetString("uuid")
	pluginIsExist, err := model.SharePluginModel().PluginIsExistByUUID(uuid)
	if err != nil {
		return err
	}
	pluginName := jsonReader.GetString("name")
	version := jsonReader.GetString("version")
	pluginType := p.languageToPluginType(jsonReader.GetString("type"))
	if pluginType == info.PluginType_None {
		return errors.New("error plugin")
	}

	// 2.2 save db
	fmt.Println(pluginIsExist)
	fmt.Println(uuid)
	fmt.Println(pluginName)
	fmt.Println(pluginType)
	fmt.Println(version)
	if pluginIsExist {
		// update plugin
		err = model.SharePluginModel().UpdatePlugin(uuid, pluginName, pluginType, version)
	} else {
		err = model.SharePluginModel().InsertPlugin(uuid, pluginName, pluginType, version)
	}
	if err != nil {
		return err
	}

	// 2.3 save file
	pluginRootPath := filepath.Join(c.GetString("storage.file.plugin"), uuid)
	os.MkdirAll(pluginRootPath, os.ModePerm)

	// 2.4 save code.zip
	tmpCodePath := filepath.Join(tmpRawPath, "raw", "code.zip")
	saveCodePath := filepath.Join(pluginRootPath, "code.zip")
	err = os.Rename(tmpCodePath, saveCodePath)
	if err != nil {
		return err
	}

	// 2.5 Extract code.zip to current path
	extractPath := filepath.Join(pluginRootPath, "run")
	fmt.Println("extractPath: ", extractPath)
	err = archive.UnZipToPath(saveCodePath, tmpRawPath)
	if err != nil {
		return err
	}
	err = os.Rename(filepath.Join(tmpRawPath, "code"), extractPath)
	if err != nil {
		return err
	}

	// 2.6 copy cover img, config and so on to here
	tmpBigCoverImgPath := filepath.Join(tmpRawPath, "raw", "big_cover.jpg")
	tmpSmallCoverImgPath := filepath.Join(tmpRawPath, "raw", "small_cover.jpg")
	saveBigCoverImgPath := filepath.Join(pluginRootPath, "big_cover.jpg")
	saveSmallCoverImgPath := filepath.Join(pluginRootPath, "small_cover.jpg")
	saveConfigPath := filepath.Join(pluginRootPath, "plugin.conf")

	err = os.Rename(tmpBigCoverImgPath, saveBigCoverImgPath)
	if err != nil {
		return err
	}

	err = os.Rename(tmpSmallCoverImgPath, saveSmallCoverImgPath)
	if err != nil {
		return err
	}

	err = os.Rename(tmpConfigPath, saveConfigPath)
	if err != nil {
		return err
	}

	// 2.7 copy raw.zip
	saveRawPath := filepath.Join(rawPath, uuid+".zip")
	fmt.Println("raw: ", saveRawPath)
	err = os.Rename(p.rawPluginPath, saveRawPath)
	if err != nil {
		return err
	}

	return nil
}

func (p *pluginStorage) languageToPluginType(language string) int {
	switch language {
	case "golang":
		return info.PluginType_Golang
	case "cpp":
		return info.PluginType_CPP
	case "html":
		return info.PluginType_H5
	case "python":
		return info.PluginType_Python
	case "java":
		return info.PluginType_Java
	case "node":
		return info.PluginType_Node
	default:
		return info.PluginType_None
	}
}
