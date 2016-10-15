package info

const (
	// none
	PluginType_None = iota

	// html
	PluginType_H5 = iota

	// native code
	PluginType_CPP    = iota
	PluginType_Java   = iota
	PluginType_Golang = iota

	// script
	PluginType_Node   = iota
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
