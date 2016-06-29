package personal

import (
	"sync"
)

type authManager struct {
}

var authOnce sync.Once
var authManagerInstance *authManager

func ShareAuthManager() *authManager {
	authOnce.Do(func() {
		authManagerInstance = &authManager{}
	})
	return authManagerInstance
}

func (a *authManager) handle() {

}
