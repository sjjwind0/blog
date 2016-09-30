package personal

import (
	"encoding/json"
	"errors"
	"fmt"
	"framework"
	"framework/base/config"
	"framework/response"
	"framework/server"
	"io/ioutil"
	"model"
	"net/http"
	"os"
	"path/filepath"
	"info"
)

type PersonalDeleteController struct {
	server.SessionController
}

func NewPersonalDeleteController() *PersonalDeleteController {
	return &PersonalDeleteController{}
}

func (p *PersonalDeleteController) Path() interface{} {
	return "/personal/delete"
}

func (p *PersonalDeleteController) SessionPath() string {
	return "/"
}

func (p *PersonalDeleteController) HandlerRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		response.JsonResponse(w, framework.ErrorMethodError)
		return
	}
	p.SessionController.HandlerRequest(p, w, r)

	if status, err := p.WebSession.Get("status"); err == nil && status == "auth" {
		response.JsonResponse(w, framework.ErrorOK)
		return
	}

	result, err := ioutil.ReadAll(r.Body)
	r.Body.Close()
	if err != nil {
		response.JsonResponseWithMsg(w, framework.ErrorParamError, err.Error())
		return
	}
	var f interface{}
	json.Unmarshal(result, &f)
	if m, ok := f.(map[string]interface{}); ok {
		if bid, ok := m["id"].(float64); ok {
			blogId := int(bid)
			p.deleteBlog(blogId)
			return
		}
	}
	response.JsonResponseWithMsg(w, framework.ErrorParamError, "no username")
}

func (p *PersonalDeleteController) deleteBlog(blogId int) error {
	// 1. 删除db，包括blog，comment
	fmt.Println("deleteBlog: ", blogId)
	isExist, err := model.ShareBlogModel().BlogIsExistByBlogID(blogId)
	if err != nil {
		return err
	}
	if isExist {
		blogInfo, err := model.ShareBlogModel().FetchBlogByBlogID(blogId)
		if err != nil {
			return err
		}
		err = model.ShareBlogModel().DeleteBlog(blogId)
		if err != nil {
			return err
		}
		err = model.ShareCommentModel().DeleteAllBlogComment(info.CommentType_Blog, blogId)
		if err != nil {
			return err
		}
		// 2. 删除本地文件, raw文件暂时不删
		blogPath := config.GetDefaultConfigJsonReader().Get("blog.storage.file.blog").(string)
		blogPath = filepath.Join(blogPath, blogInfo.BlogUUID)
		return os.RemoveAll(blogPath)
	}
	return errors.New("no such blog")
}
