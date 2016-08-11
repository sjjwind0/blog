package controller

import (
	"controller/login"
	"framework"
	"framework/response"
	"framework/server"
	"html/template"
	"info"
	"model"
	"net/http"
	"strconv"
)

type loginRender struct {
	Code           int
	Msg            string
	IsLoginSuccess bool
}

type LoginController struct {
	server.SessionController
	appKey    string
	appSecret string
}

func NewLoginController() *LoginController {
	ret := &LoginController{}
	ret.init()
	return ret
}

func (l *LoginController) init() {

}

func (i *LoginController) Path() interface{} {
	return "/login"
}

func (l *LoginController) SessionPath() string {
	return "/"
}

func (l *LoginController) writeLoginInfo(from string, userInfo *info.UserInfo) {
	l.WebSession.Set("from", from)
	l.WebSession.Set("id", strconv.Itoa(int(userInfo.UserID)))
	l.WebSession.Set("status", "login")
	l.ResetSessionDuration()
}

func (l *LoginController) handleLoginInfo(w http.ResponseWriter, userInfo *info.UserInfo, err error) {
	var render *loginRender = nil
	if err != nil {
		render = &loginRender{
			Code:           framework.ErrorRunTimeError,
			Msg:            err.Error(),
			IsLoginSuccess: false,
		}
	} else {
		err = model.ShareUserModel().Login(info.AccountTypeQQ, userInfo)
		if err != nil {
			render = &loginRender{
				Code:           framework.ErrorRunTimeError,
				Msg:            err.Error(),
				IsLoginSuccess: false,
			}
		} else {
			l.writeLoginInfo("qq", userInfo)
			render = &loginRender{
				Code:           framework.ErrorOK,
				Msg:            "",
				IsLoginSuccess: true,
			}
		}
	}
	t, err := template.ParseFiles("./src/view/html/login-result.html")
	t.Execute(w, render)
}

func (l *LoginController) handleLogout(w http.ResponseWriter) {
	status, err := l.WebSession.Get("status")
	if err == nil {
		if status == "login" {
			// 已经登录
			err = l.GetSessionMgr().DeleteSession(l.WebSession.SessionID())
			if err == nil {
				response.JsonResponse(w, framework.ErrorOK)
				return
			}
		}
	}
	response.JsonResponseWithMsg(w, framework.ErrorRunTimeError, err.Error())
}

func (l *LoginController) HandlerRequest(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	l.SessionController.HandlerRequest(l, w, r)
	loginType := r.Form.Get("type")
	switch loginType {
	case "qq":
		code := r.Form.Get("code")
		if len(code) != 0 {
			userInfo, err := login.GetQQLoginInstance().Login(code)
			l.handleLoginInfo(w, userInfo, err)
		}
	case "weibo":
		code := r.Form.Get("code")
		if len(code) != 0 {
			userInfo, err := login.GetWebLoginInstance().Login(code)
			l.handleLoginInfo(w, userInfo, err)
		}
	case "logout":
		l.handleLogout(w)
	default:
		response.JsonResponseWithMsg(w, framework.ErrorRunTimeError, "unsupport login type")
	}
}
