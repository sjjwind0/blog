package run

import (
	"net/http"
)

type RequestHandler interface {
	HandlePluginRequest(pluginId int, w http.ResponseWriter, r *http.Request)
}
