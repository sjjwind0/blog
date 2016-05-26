package controller

import (
	"html/template"
	"log"
	"net/http"
)

type IndexController struct {
}

func NewIndexController() *IndexController {
	return &IndexController{}
}

func (this *IndexController) HandlerAction(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("./src/view/html/index.html")
	if err != nil {
		log.Println(err)
	}
	t.Execute(w, nil)
}
