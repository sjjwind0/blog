package config

import (
	"encoding/json"
	"github.com/bitly/go-simplejson"
	"os"
	"strings"
)

const (
	ConfigType_FileConfig = iota
	ConfigType_ContentConfig
)

type configManager struct {
	config     string
	configType int
}

func IsFileExixt(path string) bool {
	_, err := os.Stat(path)
	return err == nil || os.IsExist(err)
}

func NewConfigFileManager(config string) *configManager {
	return &configManager{config, ConfigType_FileConfig}
}

func NewConfigContentManager(content string) *configManager {
	return &configManager{content, ConfigType_ContentConfig}
}

func (c *configManager) loadConfigFileContent() string {
	if IsFileExixt(c.config) == false {
		file, _ := os.Create(c.config)
		defer file.Close()
		return ""
	} else {
		fileInfo, err := os.Stat(c.config)
		if err != nil {
			return ""
		}
		file, err := os.Open(c.config)
		defer file.Close()
		if err != nil {
			return ""
		}
		content := make([]byte, fileInfo.Size())
		file.Read(content)
		return string(content)
	}
}

func (c *configManager) ReadConfig(key string) interface{} {
	content := ""
	if c.configType == ConfigType_FileConfig {
		content = c.loadConfigFileContent()
	} else {
		content = c.config
	}
	if content != "" {
		js, err := simplejson.NewJson([]byte(content))
		if err == nil {
			root, err := js.Map()
			if err == nil {
				keyList := strings.Split(key, ".")
				value := root
				for i := range keyList {
					if i == len(keyList)-1 {
						retValue := value[keyList[i]]
						if realValue, ok := retValue.(json.Number); ok {
							if v, o := realValue.Int64(); o == nil {
								return v
							}
							if v, o := realValue.Float64(); o == nil {
								return v
							}
							return realValue.String()
						}
						return retValue
					}
					var ok bool = true
					if value, ok = value[keyList[i]].(map[string]interface{}); !ok {
						break
					}
				}
			}
		}
	}
	return nil
}
