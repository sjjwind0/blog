package impl

import (
	"fmt"
	"golang.org/x/net/websocket"
)

type MessageController struct {
	ws *websocket.Conn
}

func NewMessageController() *MessageController {
	return &MessageController{}
}

func (m *MessageController) SendMessage(message string) {
	websocket.Message.Send(m.ws, message)
}

func (m *MessageController) Path() interface{} {
	return "/message"
}

func (m *MessageController) HandlerRequest(ws *websocket.Conn) {
	m.ws = ws
	msgManager := NewMessageManager(m)
	msgManager.StartNewChess()
	var err error
	for {
		var reply string
		err = websocket.Message.Receive(ws, &reply)
		if err != nil {
			fmt.Println("Error in receive!", err)
			// 连接断掉
			msgManager.MessageProc("MessageCut", "")
			return
		}
		msgManager.MessageProc("RecvData", reply)
		if err != nil {
			fmt.Println("Can't send")
			continue
		}
	}
}
