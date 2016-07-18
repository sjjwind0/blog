package session

import "time"

type Session interface {
	Set(key string, value interface{}) error
	Get(key string) (interface{}, error)
	Delete(key string) error
	SessionID() string
	InitSessionExpire(expireTime int64)
	CreateTime() int64
	MaxAge() int
	IsExpired() bool
}

type SessionExpire struct {
	createTime int64
	expireTime int64
}

func (s *SessionExpire) InitSessionExpire(expireTime int64) {
	s.createTime = time.Now().Unix()
	s.expireTime = expireTime
}

func (s *SessionExpire) CreateTime() int64 {
	return s.createTime
}

func (s *SessionExpire) MaxAge() int {
	if s.IsExpired() {
		return 0
	}
	return int((s.createTime + s.expireTime - time.Now().Unix()))
}

func (s *SessionExpire) IsExpired() bool {
	return time.Now().Unix()-s.createTime-s.expireTime >= 0
}
