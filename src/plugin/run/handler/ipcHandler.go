package handler

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"plugin/ipc"
)

type IPCRequestHandler struct {
}

func (i *IPCRequestHandler) HandlePluginRequest(pluginId int, w http.ResponseWriter, r *http.Request) {
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
	ipc.SharePluginIPCManager().CallMethod(pluginId, string(requestBytes), callback)
	<-singal
}
