package golang

import (
	"fmt"
	"testing"
	"time"
)

type TestIPCDelegate struct {
}

func (t *TestIPCDelegate) OnAcceptNewClient(manager *IPCManager, ipcID int) {
	manager.CallMethod(ipcID, "TestClientFunc1", "我有一头小毛驴我从来也不骑", func(code int, response string) {
		if code == 0 {
			fmt.Println("TestClientFunc1 response: ", response)
		} else {
			fmt.Println("callback failed")
		}
	})
}

func (t *TestIPCDelegate) OnClientClose(manager *IPCManager, ipcID int) {
	fmt.Println("OnClientClose: ", ipcID)
}

func (t *TestIPCDelegate) OnConnect(manager *IPCManager, ipcID int) {
	manager.CallMethod(ipcID, "TestServerFunc", "我有一头小毛驴我从来也不骑", func(code int, response string) {
		if code == 0 {
			fmt.Println("TestServerFunc response: ", response)
		} else {
			fmt.Println("callback failed")
		}
	})
	manager.CallMethod(ipcID, "TestServerFunc1", "我有一头小毛驴我从来也不骑", func(code int, response string) {
		if code == 0 {
			fmt.Println("TestServerFunc1 response: ", response)
		} else {
			fmt.Println("callback failed")
		}
	})
}

func (t *TestIPCDelegate) OnServerClose(manager *IPCManager) {
}

func Test_IPC(t *testing.T) {
	go func() {
		server := NewIPCManager()
		d := new(TestIPCDelegate)
		ipcID := server.CreateServer("test", d)
		server.RegisterMethod(ipcID, "TestServerFunc", func(request string, response *string) {
			*response = request
		})
		server.RegisterMethod(ipcID, "TestServerFunc1", func(request string, response *string) {
			*response = "我就是服务端返回的结果: " + request
		})
		server.StartListener()
	}()
	// waiting for server start
	time.Sleep(100 * time.Millisecond)
	go func() {
		client := NewIPCManager()
		d := new(TestIPCDelegate)
		ipcID := client.OpenClient("test", d)
		client.RegisterMethod(ipcID, "TestClientFunc1", func(request string, response *string) {
			*response = "我就是Client返回的结果: " + request
		})
		client.StartListener()
	}()
	for {
		time.Sleep(1 * time.Second)
	}
}
