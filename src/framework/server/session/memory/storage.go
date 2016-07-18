package memory

import (
	"errors"
	"framework/server/session"
)

type memoryStorage struct {
	sessionMap map[string]*memorySession
}

func NewMemoryStorage() *memoryStorage {
	return &memoryStorage{}
}

func (m *memoryStorage) Add(sessionId string, s session.Session) error {
	if m.sessionMap == nil {
		m.sessionMap = make(map[string]*memorySession)
	}
	if _, ok := m.sessionMap[sessionId]; !ok {
		m.sessionMap[sessionId] = s.(*memorySession)
	}
	return nil
}

func (m *memoryStorage) Get(sessionId string) (session.Session, error) {
	if v, ok := m.sessionMap[sessionId]; ok {
		return v, nil
	}
	return nil, errors.New("404 not found")
}

func (m *memoryStorage) Delete(sessionId string) error {
	if _, ok := m.sessionMap[sessionId]; ok {
		delete(m.sessionMap, sessionId)
		return nil
	}
	return errors.New("404 not found")
}
