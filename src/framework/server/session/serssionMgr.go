package session

import (
	"sync"
)

var sessionMgrInstance *sessoinMgr = nil
var sessionMgrOnce sync.Once

type sessoinMgr struct {
	storage *SessionStorage
}

func GetSessionMgrInstance() *sessoinMgr {
	sessionMgrOnce.Do(func() {
		sessoinMgr = &sessoinMgr{}
	})
	return sessoinMgr
}

func (s *sessoinMgr) SetStorage(storage *SessionStorage) {
	s.storage = storage
}

func (s *sessoinMgr) QuerySessionById(sessionId string) (*Session, error) {
	return s.storage.Get(sessionId)
}

func (s *sessoinMgr) AddSession(session *Session) error {
	return s.storage.Get(session.SessionID(), session)
}

func (s *sessoinMgr) DeleteSession(sessionId string) error {
	s.storage.Delete(sessionId)
}

func (s *sessoinMgr) ReloadSession() {
	// not implement
}
