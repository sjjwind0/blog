package controller

/*
 * blog需要的公有api
 */

import (
	"encoding/base64"
	"encoding/json"
	"framework"
	"framework/response"
	"framework/server"
	"html/template"
	"info"
	"io/ioutil"
	"model"
	"net/http"
)

type apiCommentRender struct {
	UserID         string
	UserName       string
	Pic            string
	CommentID      string
	CommentContent string
	CommentTime    string
	Floor          int
	User           *info.UserInfo
	ChildContent   template.HTML
}

type APIController struct {
	server.SessionController
}

func NewAPIController() *APIController {
	return &APIController{}
}

func (a *APIController) Path() interface{} {
	return "/api"
}

func (a *APIController) SessionPath() string {
	return "/"
}

func (a *APIController) buildComment(commentId int) (string, error) {
	var commentList []*info.CommentInfo = nil
	for commentId != -1 {
		comment, err := model.ShareCommentModel().FetchCommentByCommentId(commentId)
		if err != nil {
			return "", err
		}
		commentList = append(commentList, comment)
		commentId = comment.ParentCommentID
	}
	// 逆序
	commentListLength := len(commentList)
	for i := 0; i < commentListLength/2; i++ {
		tmp := commentList[i]
		commentList[i] = commentList[commentListLength-i-1]
		commentList[commentListLength-i-1] = tmp
	}
	comment := buildOneCommentFromCommentList(&commentList)
	return comment, nil
}

func (a *APIController) handlePublicCommentAction(w http.ResponseWriter, info map[string]interface{}) {
	status, err := a.WebSession.Get("status")
	if err != nil {
		response.JsonResponseWithMsg(w, framework.ErrorAccountNotLogin, err.Error())
		return
	}
	if status != "login" {
		response.JsonResponseWithMsg(w, framework.ErrorAccountNotLogin, "account not login")
		return
	}
	uid, err := a.WebSession.Get("id")
	if err != nil {
		response.JsonResponseWithMsg(w, framework.ErrorAccountNotLogin, err.Error())
		return
	}
	var userId int = int(uid.(int64))
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
	var blogId, commentId int
	var content string
	if parseInt("blogId", &blogId) && parseInt("commentId", &commentId) {
		if _, ok := info["content"]; ok {
			switch info["content"].(type) {
			case string:
				content = info["content"].(string)
				commentId, err := model.ShareCommentModel().AddComment(userId, blogId, commentId, content)
				if err == nil {
					comment, err := a.buildComment(commentId)
					if err == nil {
						var data map[string]interface{} = make(map[string]interface{})
						data["comment"] = base64.StdEncoding.EncodeToString([]byte(comment))
						response.JsonResponseWithData(w, framework.ErrorOK, "", data)
						return
					}
				}
				response.JsonResponseWithMsg(w, framework.ErrorSQLError, err.Error())
				return
			}
		}

	}
	response.JsonResponse(w, framework.ErrorParamError)
}

func (a *APIController) HandlerRequest(w http.ResponseWriter, r *http.Request) {
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
					a.SessionController.HandlerRequest(a, w, r)
					a.handlePublicCommentAction(w, info)
					return
				case "blog":
				}
			}
		}
	}
	response.JsonResponse(w, framework.ErrorParamError)
}
