package plugin

import (
	"plugin/chess"
	"sync"
)

type PluginInfo interface {
	GetPluginName() string
	GetPluginUUID() string
	GetPluginCoverURL() string
	GetPluginDisplayPath() string
	GetPluginDescription() string
	GetPluginDownloadURL() string
	GetPluginPublicTime() int64
}

type PluginRunner interface {
	ResourceHandler() []string
	NormalHanlder() []interface{}
	WebSocketHandler() []interface{}
}

var pluginManagerInstance *pluginManager = nil
var pluginManagerOnce sync.Once

type pluginManager struct {
	pluginsInfo   []PluginInfo
	PluginsRunner []PluginRunner
}

func GetDefaultPluginManager() *pluginManager {
	pluginManagerOnce.Do(func() {
		pluginManagerInstance = &pluginManager{}
		pluginManagerInstance.startRegister()
	})
	return pluginManagerInstance
}

func (p *pluginManager) startRegister() {
	newPlugin := chess.NewChessPlugin()
	p.pluginsInfo = append(p.pluginsInfo, newPlugin)
	p.PluginsRunner = append(p.PluginsRunner, newPlugin)
}

func (p *pluginManager) GetAllPluginRunner() []PluginRunner {
	return p.PluginsRunner
}

func (p *pluginManager) GetAllPluginInfo() []PluginInfo {
	return p.pluginsInfo
}
