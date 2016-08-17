package server

import (
	"fmt"
	"framework"
	"framework/response"
	"net/http"
	"os"
	"path/filepath"
	"strings"
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
	controllerMap             map[string]Controller
	staticFileMap             map[string]string
	childHandlerControllerMap map[string]Controller
	port                      int
}

var serverMgrInstance *serverMgr = nil
var serverMgrOnce sync.Once

func ShareServerMgrInstance() *serverMgr {
	serverMgrOnce.Do(func() {
		serverMgrInstance = &serverMgr{}
		serverMgrInstance.controllerMap = nil
		serverMgrInstance.staticFileMap = nil
		serverMgrInstance.port = defaultServerPort
	})
	return serverMgrInstance
}

func (s *serverMgr) RegisterController(controller Controller) {
	registerController := func(controllerMap *map[string]Controller, path interface{},
		controller Controller) {
		switch path.(type) {
		case string:
			if _, ok := (*controllerMap)[path.(string)]; ok {
				fmt.Println("controller has been registered!")
				return
			}
			(*controllerMap)[path.(string)] = controller
		case []string:
			for _, p := range path.([]string) {
				if _, ok := (*controllerMap)[p]; ok {
					fmt.Println("controller has been registered!")
					return
				}
				(*controllerMap)[p] = controller
			}
		}
	}
	if normalController, ok := controller.(NormalController); ok {
		if s.controllerMap == nil {
			s.controllerMap = make(map[string]Controller)
		}
		registerController(&s.controllerMap, normalController.Path(), normalController)
	} else if childHandlerController, ok := controller.(ChildHandlerController); ok {
		if s.childHandlerControllerMap == nil {
			s.childHandlerControllerMap = make(map[string]Controller)
		}
		path, enableChildPath := childHandlerController.Path()
		registerController(&s.controllerMap, path, childHandlerController)
		if enableChildPath {
			registerController(&s.childHandlerControllerMap, path, childHandlerController)
		}
	}
}

func (s *serverMgr) RegisterStaticFile(webPath string, localPath string) {
	if s.staticFileMap == nil {
		s.staticFileMap = make(map[string]string)
	}
	if _, ok := s.staticFileMap[webPath]; ok {
		fmt.Println("static file has beed registered!")
		return
	}
	walkPath := filepath.Join(localPath, webPath)
	filepath.Walk(walkPath, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			webFilePath := path[len(localPath)+1:]
			s.staticFileMap[webFilePath] = path
		}
		return nil
	})
}

func (s *serverMgr) SetServerPort(port int) {
	s.port = port
}

func (s *serverMgr) handlerStatisFile(w http.ResponseWriter, currentPath string) bool {
	if currentPath[0] == '/' {
		currentPath = currentPath[1:]
	}
	if local, ok := s.staticFileMap[currentPath]; ok {
		ext := filepath.Ext(local)
		contentType := ""
		if v, ok := extContentTypeMap[strings.ToLower(ext)]; ok {
			contentType = v
		} else {
			contentType = "application/octet-stream"
		}
		file, err := os.Open(local)
		if err != nil {
			response.JsonResponseWithMsg(w, framework.ErrorNoSuchFileOrDirectory, err.Error())
			return false
		}
		defer file.Close()
		fileInfo, err := os.Stat(local)
		if err != nil {
			response.JsonResponseWithMsg(w, framework.ErrorNoSuchFileOrDirectory, err.Error())
			return false
		}
		content := make([]byte, fileInfo.Size())
		file.Read(content)
		w.Header().Set("Accept", "*/*")
		w.Header().Set("Content-Length", fmt.Sprintf("%d", fileInfo.Size()))
		w.Header().Set("Content-Type", contentType)
		w.Write(content)
		return true
	}
	return false
}

func (s *serverMgr) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	currentPath := r.URL.Path
	// 1. 首先在controller里面寻找
	if controller, ok := s.controllerMap[currentPath]; ok {
		controller.HandlerRequest(w, r)
		return
	}
	// 2. 在static file 里面寻找
	if s.handlerStatisFile(w, currentPath) {
		return
	}
	// 3. 逐级分解，看是不是某个controller的子集
	for true {
		lastIndex := strings.LastIndex(currentPath, "/")
		if lastIndex != -1 {
			currentPath = currentPath[:lastIndex]
			if controller, ok := s.childHandlerControllerMap[currentPath]; ok {
				controller.HandlerRequest(w, r)
				return
			}
		} else {
			break
		}
	}
	// 4. 404
	fmt.Println("404: ", r.URL.Path)
	w.WriteHeader(http.StatusNotFound)
}

func (s *serverMgr) StartServer() {
	http.HandleFunc("/", s.ServeHTTP)
	http.ListenAndServe(fmt.Sprintf(":%d", s.port), nil)
}
