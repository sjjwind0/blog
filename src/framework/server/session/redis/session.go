package redis

import (
	"errors"
	"framework/server/session"
)

type redisSession struct {
	session.BaseSession
}

func NewRedisSession() *redisSession {
	s := &redisSession{}
	s.SessionId = session.NewSessionID()
	return s
}

func restoreRedisSessionFromRedis(sessionId string, createTime, expireTime int64) *redisSession {
	s := &redisSession{}
	s.SessionId = sessionId
	s.InitBaseSessionWithCreateTime(createTime, expireTime)
	return s
}

func (r *redisSession) Set(key string, value interface{}) error {
	if key == "" {
		return errors.New("key must not be empty")
	}
	return r.GetStorage().(*redisStorage).setSessionContent(r, key, value)
}

func (r *redisSession) Get(key string) (interface{}, error) {
	if key == "" {
		return nil, errors.New("key must not be empty")
	}
	value, err := r.GetStorage().(*redisStorage).querySessionContent(r, key)
	return value, err
}

func (r *redisSession) ResetDuration(duration int64) error {
	r.BaseSession.ResetDuration(duration)
	return r.GetStorage().(*redisStorage).refreshSessionDuration(r, r.MaxDuration())
}

func (r *redisSession) Delete(key string) error {
	if key == "" {
		return errors.New("key must not be empty")
	}
	return r.GetStorage().(*redisStorage).deleteSessionContent(r, key)
}
