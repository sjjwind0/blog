package memory

import (
	"errors"
)

type memoryStorage struct {
	sessionMap map[string]*memorySession
}

func NewMemoryStorage() *memorySession {
	return &memoryStorage{}
}

func (m *memoryStorage) Add(sessionId string, s *Session) error {
	if ok, _ := m.sessionMap[sessionId]; !ok {
		m.sessionMap[sessionId] = s
	}
}

func (m *memoryStorage) Get(sessionId string) (*Session, error) {
	if ok, v := m.sessionMap[sessionId]; ok {
		return v, nil
	}
	return nil, errors.New("404 not found")
}

func (m *memoryStorage) Delete(sessionId string) error {
	if ok, _ := m.sessionMap[sessionId]; ok {
		delete(m.sessionMap, sessionId)
		return nil
	}
	return errors.New("404 not found")
}
