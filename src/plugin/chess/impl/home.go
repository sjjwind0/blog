package impl

import (
	"html/template"
	"log"
	"net/http"
)

type HomeController struct {
}

func NewHomeController() *HomeController {
	return &HomeController{}
}

func (h *HomeController) Path() interface{} {
	return "/plugin/chess"
}

func (h *HomeController) HandlerRequest(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("./src/plugin/chess/res/chess.html")
	if err != nil {
		log.Println(err)
	}

	t.Execute(w, nil)
}
