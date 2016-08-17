package controller

import (
	"fmt"
	"html/template"
	"net/http"
)

type AboutController struct {
}

func NewAboutController() *AboutController {
	return &AboutController{}
}

func (a *AboutController) Path() interface{} {
	return "/about"
}

func (a *AboutController) HandlerRequest(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("./src/view/html/about.html")
	if err != nil {
		fmt.Println("parse file error: ", err.Error())
	}
	t.Execute(w, nil)
}
