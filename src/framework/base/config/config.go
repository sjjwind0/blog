package config

import (
	"framework/base/json"
	"sync"
)

const (
	ConfigType_FileConfig = iota
	ConfigType_ContentConfig
)

const kDefaultConfigName = "default.conf"

var configMgrMap map[string]*json.JsonReader = make(map[string]*json.JsonReader)
var configMapListLock sync.Mutex

func GetDefaultConfigJsonReader() *json.JsonReader {
	return GetGlobalConfigFileManager(kDefaultConfigName)
}

func GetGlobalConfigFileManager(configPath string) *json.JsonReader {
	configMapListLock.Lock()
	defer configMapListLock.Unlock()
	if value, ok := configMgrMap[configPath]; ok {
		return value
	}
	configMgrMap[configPath] = json.NewJsonReaderFromFile(configPath)
	return configMgrMap[configPath]
}
