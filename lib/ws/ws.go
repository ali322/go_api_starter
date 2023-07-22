package ws

import (
	"encoding/json"
	"sync"

	"github.com/gorilla/websocket"
)

type ThreadSafeWriter struct {
	*websocket.Conn
	sync.Mutex
}

func (t *ThreadSafeWriter) WriteJSON(v interface{}) error {
	t.Lock()
	defer t.Unlock()

	return t.Conn.WriteJSON(v)
}

func (t *ThreadSafeWriter) WriteReply(method string, params interface{}) error {
	t.Lock()
	defer t.Unlock()

	data, err := json.Marshal(params)
	if err != nil {
		return err
	}
	return t.Conn.WriteJSON(map[string]interface{}{
		"method": method, "params": string(data),
	})
}

func (t *ThreadSafeWriter) WriteNotify(event string, data interface{}) error {
	t.Lock()
	defer t.Unlock()

	return t.Conn.WriteJSON(map[string]interface{}{
		"event": event, "data": data,
	})
}

type Message struct {
	Event string      `json:"event"`
	Data  interface{} `json:"data"`
	From  string      `json:"from"`
	Only  []string    `json:"only"`
}

type Notify struct {
	Method string `json:"method"`
	Params string `json:"params"`
}

type broadcast struct {
	message *Message
	uid     string
	sid     string
}
