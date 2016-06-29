package session

type SessionStorage interface {
	Add(sessionId string, s *Session) error
	Get(sessionId string) (*Session, error)
	Delete(sessionId string) error
}
