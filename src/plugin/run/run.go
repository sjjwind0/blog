package run

type PluginRun interface {
	Run(pluginId int) error
	Stop() error
}
