package server

import (
	"fmt"
	"framework/config"
	"framework/server/session"
	"framework/server/session/memory"
	"framework/server/session/redis"
	"net/http"
)

// session两天过期
const (
	kSessionMaxAge = 60 * 60 * 24 * 2
)

var sessionMgrInstance *session.SessoinMgr = nil

type Controller interface {
	HandlerRequest(w http.ResponseWriter, r *http.Request)
}

type NormalController interface {
	Controller
	Path() interface{}
}

type ChildHandlerController interface {
	Controller
	Path() (interface{}, bool)
}

type SessionControllerInterface interface {
	SessionPath() string
}

type SessionController struct {
	WebSession session.Session
}

func (s *SessionController) HandlerRequest(controller SessionControllerInterface,
	w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("s")
	cookiePath := controller.SessionPath()
	if err != nil {
		fmt.Println("get session id error: ", err)
		newSession := s.newSession()
		s.GetSessionMgr().AddSession(newSession)
		c := http.Cookie{Name: "s", Value: newSession.SessionID(), Path: cookiePath, MaxAge: kSessionMaxAge}
		http.SetCookie(w, &c)
		s.WebSession = newSession
	} else {
		c, err := s.GetSessionMgr().QuerySessionById(cookie.Value)
		if err != nil {
			// 找不到session，可能是过期了，或者重启导致session丢失了，new一个session
			newSession := s.newSession()
			s.GetSessionMgr().AddSession(newSession)
			cc := http.Cookie{Name: "s", Value: newSession.SessionID(), Path: cookiePath, MaxAge: kSessionMaxAge}
			http.SetCookie(w, &cc)
			s.WebSession = newSession
		} else {
			s.WebSession = c
			if s.WebSession.IsExpired() {
				fmt.Println(s.WebSession.SessionID() + " is expired")
				// 已经过期，分配一个新的sid
				s.GetSessionMgr().DeleteSession(s.WebSession.SessionID())
				newSession := s.newSession()
				s.GetSessionMgr().AddSession(newSession)
				cc := http.Cookie{Name: "s", Value: newSession.SessionID(), Path: cookiePath, MaxAge: kSessionMaxAge}
				http.SetCookie(w, &cc)
				s.WebSession = newSession
			} else {
				cc := http.Cookie{Name: "s", Value: c.SessionID(), Path: cookiePath, MaxAge: s.WebSession.MaxAge()}
				http.SetCookie(w, &cc)
			}
		}
	}
}

func (s *SessionController) GetSessionMgr() *session.SessoinMgr {
	if sessionMgrInstance == nil {
		sessionMgrInstance = session.NewSessionManager(s.newSessionStorage())
	}
	return sessionMgrInstance
}

func (s *SessionController) ResetSessionDuration() {
	s.WebSession.ResetDuration(kSessionMaxAge)
}

func (s *SessionController) newSession() session.Session {
	sessionType := config.GetDefaultConfigFileManager().ReadConfig("blog.storage.session.type").(string)
	switch sessionType {
	case "redis":
		ss := redis.NewRedisSession()
		ss.InitBaseSession(kSessionMaxAge)
		return ss
	case "memory":
		ss := s.newSession()
		ss.InitBaseSession(kSessionMaxAge)
		return ss
	default:
		panic("unsupport session type")
	}
	return nil
}

func (s *SessionController) newSessionStorage() session.SessionStorage {
	defaultConfig := config.GetDefaultConfigFileManager()
	sessionType := defaultConfig.ReadConfig("blog.storage.session.type").(string)
	switch sessionType {
	case "redis":
		host := defaultConfig.ReadConfig("blog.storage.session.host").(string)
		port := defaultConfig.ReadConfig("blog.storage.session.port").(string)
		password := defaultConfig.ReadConfig("blog.storage.session.password").(string)
		db := int(defaultConfig.ReadConfig("blog.storage.session.db").(int64))
		storage := redis.NewRedisStorage(host, port, password, db)
		storage.SetSessionName("WebSession")
		return storage
	case "memory":
		return memory.NewMemoryStorage()
	default:
		panic("unsupport session type")
	}
	return nil
}
