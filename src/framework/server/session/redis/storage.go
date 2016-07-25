package redis

import (
	"errors"
	"framework/server/session"
	"gopkg.in/redis.v4"
)

type redisStorage struct {
	client      *redis.Client
	sessionName string
}

func NewRedisStorage(host, port, password string, db int) *redisStorage {
	instance := &redisStorage{}
	addr := host + ":" + port
	instance.client = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})
	return instance
}

func (r *redisStorage) Add(sessionId string, s session.Session) error {
	return r.addSession(sessionId, s)
}

func (r *redisStorage) Get(sessionId string) (session.Session, error) {
	s, err := r.getSession(sessionId)
	return s, err
}

func (r *redisStorage) Delete(sessionId string) error {
	return errors.New("404 not found")
}

func (r *redisStorage) SetSessionName(name string) {
	r.sessionName = name
}

func (r *redisStorage) addSession(sessionId string, s session.Session) error {
	key := "com.session.object." + r.sessionName + "." + sessionId
	// add value
	value := "value.null"
	err := r.client.Set(key, value, 0).Err()
	if err != nil {
		return err
	}
	// add createTime
	key = "com.session.create." + r.sessionName + "." + sessionId
	err = r.client.Set(key, s.CreateTime(), 0).Err()
	if err != nil {
		return err
	}
	// add expireTime
	key = "com.session.expire." + r.sessionName + "." + sessionId
	err = r.client.Set(key, s.ExpireTime(), 0).Err()
	return err
}

func (r *redisStorage) getSession(sessionId string) (session.Session, error) {
	key := "com.session.object." + r.sessionName + "." + sessionId
	_, err := r.client.Get(key).Result()
	if err == nil {
		key = "com.session.create." + r.sessionName + "." + sessionId
		createTime, err := r.client.Get(key).Int64()
		if err != nil {
			return nil, err
		}
		key = "com.session.expire." + r.sessionName + "." + sessionId
		expireTime, err := r.client.Get(key).Int64()
		if err != nil {
			return nil, err
		}
		return restoreRedisSessionFromRedis(sessionId, createTime, expireTime), nil
	}
	return nil, err
}

func (r *redisStorage) deleteSession(sessionId string) error {
	return nil
}

func (r *redisStorage) setSessionContent(s session.Session, key string, value interface{}) error {
	sessionId := s.SessionID()
	insertKey := "com.session.object." + r.sessionName + "." + sessionId + "." + key
	return r.client.Set(insertKey, value, 0).Err()
}

func (r *redisStorage) querySessionContent(s session.Session, key string) (string, error) {
	sessionId := s.SessionID()
	queryKey := "com.session.object." + r.sessionName + "." + sessionId + "." + key
	value, err := r.client.Get(queryKey).Result()
	return value, err
}

func (r *redisStorage) deleteSessionContent(s *redisSession, key string) error {
	sessionId := s.SessionID()
	deleteKey := "com.session.object." + r.sessionName + "." + sessionId + "." + key
	return r.client.Del(deleteKey).Err()
}
