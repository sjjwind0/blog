package session

type SessoinMgr struct {
	storage SessionStorage
}

func NewSessionManager(storage SessionStorage) *SessoinMgr {
	sessionMgrInstance := &SessoinMgr{}
	sessionMgrInstance.storage = storage
	return sessionMgrInstance
}

func (s *SessoinMgr) QuerySessionById(sessionId string) (Session, error) {
	ss, err := s.storage.Get(sessionId)
	if err == nil {
		ss.setSessionStorage(s.storage)
	}
	return ss, err
}

func (s *SessoinMgr) AddSession(session Session) error {
	session.setSessionStorage(s.storage)
	return s.storage.Add(session.SessionID(), session)
}

func (s *SessoinMgr) DeleteSession(sessionId string) error {
	return s.storage.Delete(sessionId)
}

func (s *SessoinMgr) ReloadSession() {
	// not implement
}
