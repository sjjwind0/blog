package server

import (
	"fmt"
	"framework/server/session"
	"framework/server/session/memory"
	"net/http"
)

// session两天过期
const (
	kSessionMaxAge = 60 * 24 * 2
)

var isFirst bool = true

type Controller interface {
	HandlerRequest(w http.ResponseWriter, r *http.Request)
	Path() interface{}
}

type SessionControllerInterface interface {
	SessionPath() string
}

type SessionController struct {
	WebSession session.Session
}

func (s *SessionController) init() {
	if isFirst {
		ss := s.GetSessionMgr()
		ss.SetStorage(memory.NewMemoryStorage())
		isFirst = false
	}
}

func (s *SessionController) HandlerRequest(controller SessionControllerInterface,
	w http.ResponseWriter, r *http.Request) {
	s.init()
	cookie, err := r.Cookie("s")
	cookiePath := controller.SessionPath()
	if err != nil {
		fmt.Println("get session id error: ", err)
		sss := memory.NewMemorySession()
		sss.InitSessionExpire(kSessionMaxAge)
		s.GetSessionMgr().AddSession(sss)
		c := http.Cookie{Name: "s", Value: sss.SessionID(), Path: cookiePath, MaxAge: kSessionMaxAge}
		http.SetCookie(w, &c)
		s.WebSession = sss
	} else {
		c, err := s.GetSessionMgr().QuerySessionById(cookie.Value)
		if err != nil {
			// 找不到session，可能是过期了，或者重启导致session丢失了，new一个session
			sss := memory.NewMemorySession()
			sss.InitSessionExpire(kSessionMaxAge)
			s.GetSessionMgr().AddSession(sss)
			cc := http.Cookie{Name: "s", Value: sss.SessionID(), Path: cookiePath, MaxAge: kSessionMaxAge}
			http.SetCookie(w, &cc)
			s.WebSession = sss
		} else {
			s.WebSession = c
			if s.WebSession.IsExpired() {
				fmt.Println(s.WebSession.SessionID() + " is expired")
				// 已经过期，分配一个新的sid
				s.GetSessionMgr().DeleteSession(s.WebSession.SessionID())
				sss := memory.NewMemorySession()
				sss.InitSessionExpire(kSessionMaxAge)
				s.GetSessionMgr().AddSession(sss)
				cc := http.Cookie{Name: "s", Value: sss.SessionID(), Path: cookiePath, MaxAge: kSessionMaxAge}
				http.SetCookie(w, &cc)
				s.WebSession = sss
			} else {
				cc := http.Cookie{Name: "s", Value: c.SessionID(), Path: cookiePath, MaxAge: s.WebSession.MaxAge()}
				http.SetCookie(w, &cc)
			}
		}
	}
}

func (s *SessionController) GetSessionMgr() *session.SessoinMgr {
	return session.GetSessionManager("web")
}
