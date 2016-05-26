package server

import (
	"net/http"
)

type Controller interface {
	HandlerAction(w http.ResponseWriter, r *http.Request)
}
