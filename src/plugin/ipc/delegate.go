package ipc

type PluginDelegate interface {
	OnPluginReady(pluginId int)
	OnPluginShutdown(pluginId int)
}
