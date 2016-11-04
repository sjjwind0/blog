package handler

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"run"
	"sync"
	"time"
)

type GolangIPCDelegate struct {
}

var ipcMainRunnerOnce sync.Once
var ipcMainRunnerInstance *ipcMainRunner = nil

type ipcMainRunner struct {
	run.BasicRunnerManager
	pluginName    string
	clientManager *IPCManager
	clientIPCID   int
	serverManager *IPCManager
	serverIPCID   int
}

func GetMainRunner(name string) run.RunnerMgr {
	return &ipcMainRunner{pluginName: name}
}

func (i *ipcMainRunner) OnConnect(manager *IPCManager, ipcID int) {
	fmt.Println("connect server success!")
}

func (i *ipcMainRunner) OnAcceptNewClient(manager *IPCManager, ipcID int) {
	fmt.Println("accept new client!")
}

func (i *ipcMainRunner) OnClientClose(manager *IPCManager, ipcID int) {
}

func (i *ipcMainRunner) OnServerClose(manager *IPCManager) {
}

func (i *ipcMainRunner) HandleIPCRequest(request string, response *string) {
	i.BasicRunnerManager.HandleRequest(run.NormalRequest(request), response)
}

func (i *ipcMainRunner) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	mapStringInterfaceToMapStringString := func(value map[string]interface{}) map[string]string {
		var retMap map[string]string = make(map[string]string)
		for k, v := range value {
			retMap[k] = v.(string)
		}
		return retMap
	}
	// build normal request
	var stringArrayMerge = func(header map[string][]string) map[string]string {
		var ret map[string]string = make(map[string]string)
		for key, value := range header {
			ret[key] = value[0]
		}
		return ret
	}
	var httpRequest map[string]interface{} = map[string]interface{}{
		"method": r.Method,
		"url":    r.URL.String(),
		"header": stringArrayMerge(r.Header),
		"length": r.ContentLength,
		"host":   r.Host,
		"addr":   r.RemoteAddr,
		"form":   stringArrayMerge(r.Form),
	}
	httpRequestBytes, _ := json.Marshal(httpRequest)
	var req map[string]interface{} = map[string]interface{}{
		"type":    "com.request.http",
		"request": string(httpRequestBytes),
	}
	singal := make(chan bool)
	requestBytes, _ := json.Marshal(req)
	callback := func(code int, response string) {
		if code == 0 {
			// to http response
			var js interface{} = nil
			err := json.Unmarshal([]byte(response), &js)
			if err != nil {
				fmt.Println("ServerHTTP json.Unmarshal error: ", err)
				singal <- true
				return
			}
			if v, ok := js.(map[string]interface{}); ok {
				if header, ok := v["header"]; ok {
					nativeHeader := mapStringInterfaceToMapStringString(header.(map[string]interface{}))
					for key, value := range nativeHeader {
						w.Header().Set(key, value)
					}
				}
				if length, ok := v["length"]; ok {
					w.Header().Set("Content-Length", fmt.Sprintf("%d", int64(length.(float64))))
				}
				if code, ok := v["code"]; ok {
					w.WriteHeader(int(code.(float64)))
				}
				if body, ok := v["body"]; ok {
					nativeBody := body.(string)
					data, _ := base64.StdEncoding.DecodeString(nativeBody)
					w.Write(data)
				}
			}
		} else {
			fmt.Println("callback error: ", response)
		}
		singal <- true
	}
	i.serverManager.CallMethod(i.serverIPCID, "HttpRequest", string(requestBytes), callback)
	<-singal
}

func (i *ipcMainRunner) StartServer() {
	client := NewIPCManager()
	ipcID := client.OpenClient(i.pluginName, i)
	i.clientManager = client
	i.clientIPCID = ipcID
	client.RegisterMethod(ipcID, "HttpRequest", func(request string, response *string) {
		i.HandleIPCRequest(request, response)
	})
	client.StartListener()
	fmt.Println("client over")
	for {
		time.Sleep(time.Second * 1)
	}
}
