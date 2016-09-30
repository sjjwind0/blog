package info

const (
	// none
	PluginType_None = iota

	// html
	PluginType_H5 = iota

	// native code
	PluginType_CPP    = iota
	PluginType_JAVA   = iota
	PluginType_GOLANG = iota

	// script
	PluginType_NODE   = iota
	PluginType_Python = iota
)

type PluginInfo struct {
	PluginID           int
	PluginUUID         string
	PluginType         int
	PluginName         string
	PluginVersion      string
	PluginTime         int64
	PluginVisitCount   int
	PluginPraiseCount  int
	PluginDissentCount int
}
