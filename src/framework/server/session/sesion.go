package session

import "time"

type Session interface {
	Set(key string, value interface{}) error
	Get(key string) (interface{}, error)
	Delete(key string) error
	SessionID() string
	InitBaseSession(expireTime int64)
	CreateTime() int64
	ExpireTime() int64
	MaxAge() int
	MaxDuration() time.Duration
	IsExpired() bool
	setSessionStorage(storage SessionStorage)
	GetStorage() SessionStorage
}

type BaseSession struct {
	createTime int64
	expireTime int64
	storage    SessionStorage
	SessionId  string
}

func (b *BaseSession) SessionID() string {
	return b.SessionId
}

func (b *BaseSession) InitBaseSession(expireTime int64) {
	b.createTime = time.Now().Unix()
	b.expireTime = expireTime
}

func (b *BaseSession) InitBaseSessionWithCreateTime(createTime, expireTime int64) {
	b.createTime = createTime
	b.expireTime = expireTime
}

func (b *BaseSession) CreateTime() int64 {
	return b.createTime
}

func (b *BaseSession) ExpireTime() int64 {
	return b.expireTime
}

func (b *BaseSession) MaxAge() int {
	if b.IsExpired() {
		return 0
	}
	return int((b.createTime + b.expireTime - time.Now().Unix()))
}

func (b *BaseSession) MaxDuration() time.Duration {
	return time.Duration(b.MaxAge()) * time.Second
}

func (b *BaseSession) IsExpired() bool {
	return time.Now().Unix()-b.createTime-b.expireTime >= 0
}

func (b *BaseSession) setSessionStorage(storage SessionStorage) {
	b.storage = storage
}

func (b *BaseSession) GetStorage() SessionStorage {
	return b.storage
}
