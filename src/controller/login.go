package controller

import (
	"controller/login"
	"fmt"
	"framework"
	"framework/response"
	"framework/server"
	"info"
	"model"
	"net/http"
	"time"
)

type LoginController struct {
	server.SessionController
	appKey         string
	appSecret      string
	longConnectMap map[string]chan bool
}

func NewLoginController() *LoginController {
	ret := &LoginController{}
	ret.init()
	return ret
}

func (l *LoginController) init() {

}

func (i *LoginController) Path() interface{} {
	return []string{"/login", "/connect"}
}

func (l *LoginController) SessionPath() string {
	return "/"
}

func (l *LoginController) writeLoginInfo(from string, userInfo *info.UserInfo) {
	fmt.Println("writeLoginInfo: ", userInfo.UserID)
	l.WebSession.Set("from", from)
	l.WebSession.Set("id", userInfo.UserID)
	l.WebSession.Set("status", "login")
}

func (l *LoginController) handleLoginConnect(w http.ResponseWriter, successChan chan bool) {
	var index int = 0
	var maxIndex = 10
	var stop bool = false
	go func(stop *bool) {
		for {
			if *stop {
				break
			}
			heart := []byte(`{"code": 0}`)
			w.Write(heart)
			w.(http.Flusher).Flush()
			time.Sleep(10 * time.Second)
			index++
			if index >= maxIndex {
				res := []byte(`{"code": 2}`)
				w.Write(res)
				delete(l.longConnectMap, l.WebSession.SessionID())
				return
			}
		}
	}(&stop)
	<-successChan
	stop = true
	res := []byte(`{"code": 1}`)
	w.Write(res)
	delete(l.longConnectMap, l.WebSession.SessionID())
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
			if err != nil {
				response.JsonResponseWithMsg(w, framework.ErrorRunTimeError, err.Error())
			} else {
				err = model.ShareUserModel().Login(info.AccountTypeQQ, userInfo)
				if err != nil {
					response.JsonResponseWithMsg(w, framework.ErrorRunTimeError, err.Error())
					return
				}
				l.writeLoginInfo("qq", userInfo)
				if _, ok := l.longConnectMap[l.WebSession.SessionID()]; ok {
					l.longConnectMap[l.WebSession.SessionID()] <- true
					response.JsonResponse(w, framework.ErrorOK)
				} else {
					response.JsonResponseWithMsg(w, framework.ErrorRunTimeError, "login timeout")
				}
			}
		}
	case "weibo":
		fmt.Println("login from weibo")
	case "connect":
		l.longConnectMap[l.WebSession.SessionID()] = make(chan bool)
		l.handleLoginConnect(w, l.longConnectMap[l.WebSession.SessionID()])
		fmt.Println("login connect")
	}
}
