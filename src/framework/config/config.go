package config

import (
	"encoding/json"
	"github.com/bitly/go-simplejson"
	"os"
	"strings"
	"sync"
)

const (
	ConfigType_FileConfig = iota
	ConfigType_ContentConfig
)

const kDefaultConfigName = "default.conf"

var configMgrMap map[string]*configManager = make(map[string]*configManager)
var configMapListLock sync.Mutex

type configManager struct {
	config     string
	configType int
}

func IsFileExixt(path string) bool {
	_, err := os.Stat(path)
	return err == nil || os.IsExist(err)
}

func GetDefaultConfigFileManager() *configManager {
	return GetGlobalConfigFileManager(kDefaultConfigName)
}

func GetGlobalConfigFileManager(config string) *configManager {
	configMapListLock.Lock()
	defer configMapListLock.Unlock()
	if value, ok := configMgrMap[config]; ok {
		return value
	}
	configMgrMap[config] = &configManager{config, ConfigType_FileConfig}
	return configMgrMap[config]
}

func GetConfigFileManager(config string) *configManager {
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

func transfer(value interface{}) interface{} {
	switch value.(type) {
	case int, int32, string, int64, float32, float64, []string, []int, []int64, []float32:
		return value
	case json.Number:
		realValue := value.(json.Number)
		if v, o := realValue.Int64(); o == nil {
			return v
		}
		if v, o := realValue.Float64(); o == nil {
			return v
		}
		return realValue.String()
	case map[string]interface{}:
		var retMap map[string]interface{} = make(map[string]interface{})
		realMap := value.(map[string]interface{})
		for k, v := range realMap {
			retMap[k] = transfer(v)
		}
		return retMap
	case []interface{}:
		realList := value.([]interface{})
		var retList []interface{}
		for _, v := range realList {
			retList = append(retList, transfer(v))
		}
		return retList
	}
	return nil
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
						return transfer(value[keyList[i]])
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
