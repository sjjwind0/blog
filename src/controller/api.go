package controller

/*
 * blog需要的公有api
 */

import (
	"encoding/json"
	"framework"
	"framework/response"
	"io/ioutil"
	"model"
	"net/http"
)

type APIController struct {
}

func NewAPIController() *APIController {
	return &APIController{}
}

func (a *APIController) Path() interface{} {
	return "/api"
}

func (a *APIController) handlePublicCommentAction(w http.ResponseWriter, info map[string]interface{}) {
	parseInt := func(name string, retValue *int) bool {
		var ok bool
		if _, ok = info[name]; ok {
			switch info[name].(type) {
			case int, int32, int64:
				*retValue = info[name].(int)
			case float32:
				*retValue = int(info[name].(float32))
			case float64:
				*retValue = int(info[name].(float64))
			default:
				return false
			}
			return true
		}
		return false
	}
	var blogId, userId, commentId int
	var content string
	if parseInt("blogId", &blogId) && parseInt("userId", &userId) && parseInt("commentId", &commentId) {
		if _, ok := info["content"]; ok {
			switch info["content"].(type) {
			case string:
				content = info["content"].(string)
				err := model.ShareCommentModel().AddComment(userId, blogId, commentId, content)
				if err == nil {
					response.JsonResponse(w, framework.ErrorOK)
				} else {
					response.JsonResponseWithMsg(w, framework.ErrorSQLError, err.Error())
				}
				return
			}
		}

	}
	response.JsonResponse(w, framework.ErrorParamError)
}

func (a *APIController) HandlerAction(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		response.JsonResponse(w, framework.ErrorMethodError)
		return
	}
	result, err := ioutil.ReadAll(r.Body)
	r.Body.Close()
	if err != nil {
		response.JsonResponse(w, framework.ErrorParamError)
		return
	}
	var f interface{}
	json.Unmarshal(result, &f)
	switch f.(type) {
	case map[string]interface{}:
		info := f.(map[string]interface{})
		if api, ok := info["type"]; ok {
			switch api.(type) {
			case string:
				switch api.(string) {
				case "talk":
					a.handlePublicCommentAction(w, info)
					return
				case "blog":
				}
			}
		}
	}
	response.JsonResponse(w, framework.ErrorParamError)
}
