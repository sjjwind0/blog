package server

import (
	"container/list"
	"fmt"
	"net/http"
	"sync"
)

const defaultServerPort = 8080

type controllerElement struct {
	webPath    string
	controller Controller
}

type staticFileElement struct {
	webPath   string
	localPath string
}

type serverMgr struct {
	controllerList *list.List
	staticFileList *list.List
	port           int
}

var serverMgrInstance *serverMgr = nil
var serverMgrOnce sync.Once

func ShareServerMgrInstance() *serverMgr {
	serverMgrOnce.Do(func() {
		serverMgrInstance = &serverMgr{}
		serverMgrInstance.controllerList = nil
		serverMgrInstance.staticFileList = nil
		serverMgrInstance.port = defaultServerPort
	})
	return serverMgrInstance
}

func (s *serverMgr) RegisterController(controller Controller) {
	if s.controllerList == nil {
		s.controllerList = list.New()
	}
	s.controllerList.PushBack(controller)
}

func (s *serverMgr) RegisterStaticFile(webPath string, localPath string) {
	if s.staticFileList == nil {
		s.staticFileList = list.New()
	}
	s.staticFileList.PushBack(&staticFileElement{webPath, localPath})
}

func (s *serverMgr) SetServerPort(port int) {
	s.port = port
}

func (s *serverMgr) StartServer() {
	// register controller
	if s.controllerList != nil {
		for controller := s.controllerList.Front(); controller != nil; controller = controller.Next() {
			element := controller.Value.(Controller)
			path := element.Path()
			switch path.(type) {
			case string:
				http.HandleFunc(path.(string), element.HandlerRequest)
			case []string:
				pathList := path.([]string)
				for _, p := range pathList {
					http.HandleFunc(p, element.HandlerRequest)
				}
			}
		}
	}

	// register static file
	if s.staticFileList != nil {
		for file := s.staticFileList.Front(); file != nil; file = file.Next() {
			element := file.Value.(*staticFileElement)
			http.Handle(element.webPath, http.FileServer(http.Dir(element.localPath)))
		}
	}

	fmt.Println("server at port: ", s.port)
	http.ListenAndServe(fmt.Sprintf(":%d", s.port), nil)
}
