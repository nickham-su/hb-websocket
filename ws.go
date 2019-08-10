package WS

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"time"
)

func New(url string) *WS {
	ws := &WS{
		url,
		nil,
		make(map[string]func(json.RawMessage)),
		[]func([]byte, *WS) bool{},
	}
	ws.newConnect()
	// 注册消息处理方法
	ws.msgHandlers = append(ws.msgHandlers,
		pingHandler,
		subHandler,
	)
	ws.readMessage()
	return ws
}

type WS struct {
	url         string
	connect     *websocket.Conn
	subChannel  map[string]func(json.RawMessage)
	msgHandlers []func([]byte, *WS) bool
}

func (ws *WS) readMessage() {
	go func() {
		for {
			_, message, err := ws.connect.ReadMessage()
			if err != nil {
				fmt.Println("read:", err)
				// 重连
				ws.reconnect()
			}
			go ws.handle(message)
		}
	}()
}

func (ws *WS) handle(message []byte) {
	if msg, err := gzipDecode(message); err == nil {
		for _, handler := range ws.msgHandlers {
			if handler(msg, ws) { // 任务链中返回true,表示msg被正确处理
				break
			}
		}
	}
}

// 创建连接
func (ws *WS) newConnect() {
	fmt.Println("newConnect")
	if ws.connect != nil {
		ws.connect.Close()
		ws.connect = nil
	}

	c, _, err := websocket.DefaultDialer.Dial(ws.url, nil)
	if err == nil {
		ws.connect = c
	} else {
		time.Sleep(time.Second)
		ws.newConnect()
	}
}

// 重连
func (ws *WS) reconnect() {
	fmt.Println("reconnect")
	ws.newConnect()
	fmt.Println("重新订阅")
	for channel, callback := range ws.subChannel {
		ws.Sub(channel, callback)
	}
}

// 订阅
func (ws *WS) Sub(ch string, callback func(json.RawMessage)) {
	ws.subChannel[ch] = callback
	subMsg := &struct {
		Sub string `json:"sub"`
	}{ch}
	if b, err := json.Marshal(subMsg); err == nil {
		fmt.Println(string(b))
		ws.connect.WriteMessage(websocket.TextMessage, b)
	}
}
