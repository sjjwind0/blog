package controller

import (
	"net/http"
)

type BlogController struct {
	server.SessionController
	blogContentMap map[string]*[]byte
}

func NewBlogController() *BlogController {
	controller := &BlogController{}
	controller.blogContentMap = make(map[string]*[]byte)
	return controller
}

func (b *BlogController) Path() interface{} {
	return []string{"/blog", "/article", "/img"}
}

func (b *BlogController) HandlerRequest(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	if r.URL.Path == "/article" {
		id, err := strconv.Atoi(r.Form["id"][0])
		if err != nil {
			response.JsonResponseWithMsg(w, framework.ErrorParamError, "param error")
			return
		}
		b.readBlog(w, id)
	}
}
