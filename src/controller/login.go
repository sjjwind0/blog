package controller

import (
	"controller/login"
	"fmt"
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
	fmt.Println("writeLoginInfo: ", userInfo.UserID)
	l.WebSession.Set("from", from)
	l.WebSession.Set("id", strconv.Itoa(int(userInfo.UserID)))
	l.WebSession.Set("status", "login")
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
			return
		}
	case "weibo":
		fmt.Println("login from weibo")
	}
	response.JsonResponseWithMsg(w, framework.ErrorRunTimeError, "unsupport login type")
}
