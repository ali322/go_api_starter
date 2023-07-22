package ws

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/gorilla/websocket"
)

type Client struct {
	UID  string
	SID  string
	Hub  *hub
	Conn *websocket.Conn
	Send chan *Message
}

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 512
)

func (c *Client) ReadPump(handler func(msg *Message) error) {
	defer func() {
		c.Hub.unregister <- c
		c.Conn.Close()
	}()
	// c.Conn.SetReadLimit(maxMessageSize)
	c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	c.Conn.SetPongHandler(func(string) error { c.Conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	body := &Message{}
	for {
		_, raw, err := c.Conn.ReadMessage()
		if err != nil || len(raw) == 0 {
			fmt.Println("read message err", err)
			break
		}
		if err := json.Unmarshal(raw, &body); err != nil {
			fmt.Println("unmarshal message err", err)
			continue
		}
		err = handler(body)
		if err != nil {
			break
		}
	}
}

func (c *Client) Key() string {
	return c.UID + "_" + c.SID
}

func (c *Client) WritePump() {
	ticker := time.NewTicker(3 * time.Second)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.Send:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			w, err := c.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			msg, _ := json.Marshal(message)
			w.Write(msg)
			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
