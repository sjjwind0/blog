package session

var sessionMgrInstance *SessoinMgr = nil
var sessionMap map[string]*SessoinMgr = make(map[string]*SessoinMgr)

type SessoinMgr struct {
	storage SessionStorage
}

func GetSessionManager(name string) *SessoinMgr {
	if v, ok := sessionMap[name]; ok {
		return v
	}
	sessionMgrInstance = &SessoinMgr{}
	sessionMap[name] = sessionMgrInstance
	return sessionMgrInstance
}

func (s *SessoinMgr) SetStorage(storage SessionStorage) {
	s.storage = storage
}

func (s *SessoinMgr) QuerySessionById(sessionId string) (Session, error) {
	ss, err := s.storage.Get(sessionId)
	return ss, err
}

func (s *SessoinMgr) AddSession(session Session) error {
	return s.storage.Add(session.SessionID(), session)
}

func (s *SessoinMgr) DeleteSession(sessionId string) error {
	return s.storage.Delete(sessionId)
}

func (s *SessoinMgr) ReloadSession() {
	// not implement
}
