package memory

import (
	"errors"
	"framework/server/session"
	"time"
)

type memorySession struct {
	session.SessionExpire
	sessionId string
	content   map[string]interface{}
}

func NewMemorySession() *memorySession {
	s := &memorySession{}
	s.sessionId = session.NewSessionID()
	return s
}

func (m *memorySession) Set(key string, value interface{}) error {
	if key == "" {
		return errors.New("key must not be empty")
	}
	if m.content == nil {
		m.content = make(map[string]interface{})
	}
	m.content[key] = value
	return nil
}

func (m *memorySession) Get(key string) (interface{}, error) {
	if m.content == nil {
		return nil, errors.New("no session")
	}
	if key == "" {
		return nil, errors.New("key must not be empty")
	}
	if v, ok := m.content[key]; ok {
		return v, nil
	}
	return nil, errors.New("no such data")
}

func (m *memorySession) Delete(key string) error {
	if m.content == nil {
		return errors.New("no session")
	}
	if key == "" {
		return errors.New("key must not be empty")
	}
	delete(m.content, key)
	return nil
}

func (m *memorySession) SessionID() string {
	return m.sessionId
}

func (m *memorySession) CreateTime() int64 {
	return time.Now().Unix()
}
