package ipc

type PluginDelegate interface {
	OnPluginNeedStart(pluginId int)
	OnPluginShutdown(pluginId int)
}
