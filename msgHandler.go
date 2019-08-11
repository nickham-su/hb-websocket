package WS

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"github.com/nickham-su/go-logger"
)

// 心跳包处理
func pingHandler(msg []byte, ws *WS) bool {
	var pi ping
	if err := json.Unmarshal(msg, &pi); err != nil || pi.Ping == 0 {
		return false
	}

	po := &pong{pi.Ping}
	if data, err := json.Marshal(po); err == nil {
		err = ws.connect.WriteMessage(websocket.TextMessage, data)
		if err != nil {
			logger.Error.Println(err)
		}
	}
	return true
}

// 订阅数据
func subHandler(msg []byte, ws *WS) bool {
	var data subMessage
	if err := json.Unmarshal(msg, &data); err != nil || data.Ch == "" {
		return false
	}
	ck := ws.subChannel[data.Ch]
	ck(data.Tick)
	return true
}

type ping struct {
	Ping int `json:"ping"`
}
type pong struct {
	Pong int `json:"pong"`
}
type subMessage struct {
	Ch   string          `json:"ch"`
	Ts   int             `json:"ts"`
	Tick json.RawMessage `json:"tick"`
}
