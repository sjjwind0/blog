package controller

/*
 * blog需要的公有api
 */

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"framework"
	"framework/response"
	"framework/server"
	"html/template"
	"info"
	"io/ioutil"
	"model"
	"net/http"
	"strconv"
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
		comment, err := model.ShareCommentModel().FetchCommentByCommentId(info.CommentType_Blog, commentId)
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

func (a *APIController) handlePublicCommentAction(w http.ResponseWriter, inf map[string]interface{}) {
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
	userId, err := strconv.Atoi(uid.(string))
	if err != nil {
		response.JsonResponseWithMsg(w, framework.ErrorAccountNotLogin, err.Error())
		return
	}
	parseInt := func(name string, retValue *int) bool {
		var ok bool
		if _, ok = inf[name]; ok {
			switch inf[name].(type) {
			case int, int32, int64:
				*retValue = inf[name].(int)
			case float32:
				*retValue = int(inf[name].(float32))
			case float64:
				*retValue = int(inf[name].(float64))
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
		if _, ok := inf["content"]; ok {
			switch inf["content"].(type) {
			case string:
				content = inf["content"].(string)
				commentId, err := model.ShareCommentModel().AddComment(info.CommentType_Blog, userId, blogId, commentId, content)
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

func (a *APIController) handleGetUserInfoRequest(w http.ResponseWriter) {
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
	userId, err := strconv.Atoi(uid.(string))
	if err != nil {
		response.JsonResponseWithMsg(w, framework.ErrorAccountNotLogin, err.Error())
		return
	}
	userInfo, err := model.ShareUserModel().GetUserInfoById(int64(userId))
	if err != nil {
		response.JsonResponseWithMsg(w, framework.ErrorRunTimeError, err.Error())
		return
	}
	response.JsonResponseWithData(w, framework.ErrorOK, "", map[string]interface{}{
		"name": userInfo.UserName,
		"pic":  userInfo.SmallFigureurl,
	})
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
				case "getUserInfo":
					fmt.Println("getUserInfo")
					a.SessionController.HandlerRequest(a, w, r)
					a.handleGetUserInfoRequest(w)
					return
				}
			}
		}
	}
	response.JsonResponse(w, framework.ErrorParamError)
}
