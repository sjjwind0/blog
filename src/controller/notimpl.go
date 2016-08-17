package controller

import (
	"net/http"
)

type NotImplController struct {
}

func NewNotImplController() *NotImplController {
	return &NotImplController{}
}

func (i *NotImplController) Path() interface{} {
	return []string{"/high"}
}

func (i *NotImplController) HandlerRequest(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("正在开发中。。。"))
}
