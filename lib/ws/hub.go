package ws

import (
	"app/util"

	"github.com/gorilla/websocket"
)

var Hub = &hub{
	clients:    make(map[string]*Client),
	register:   make(chan *Client),
	unregister: make(chan *Client),
	broadcast:  make(chan *broadcast),
}

type hub struct {
	clients    map[string]*Client
	register   chan *Client
	unregister chan *Client
	broadcast  chan *broadcast
}

func (h *hub) Join(uid string, sid string, conn *websocket.Conn) *Client {
	c := &Client{UID: uid, SID: sid, Hub: h, Conn: conn, Send: make(chan *Message)}
	h.register <- c
	return c
}

func (h *hub) Leave(c *Client) {
	h.unregister <- c
}

func (h *hub) Broadcast(uid string, sid string, event string, data any, only []string, from string) error {
	msg := &Message{Event: event, From: from, Data: data, Only: only}
	p := &broadcast{uid: uid, sid: sid, message: msg}
	h.broadcast <- p
	return nil
}

func (h *hub) Notify(uid string, sid string, event string, data any, only []string) error {
	msg := &Message{Event: event, From: "", Data: data, Only: only}
	p := &broadcast{uid: uid, sid: sid, message: msg}
	h.broadcast <- p
	return nil
}

func (h *hub) Run() {
	for {
		select {
		case c := <-h.register:
			h.clients[c.Key()] = c
		case c := <-h.unregister:
			if _, ok := h.clients[c.Key()]; ok {
				delete(h.clients, c.Key())
			}
		case p := <-h.broadcast:
			for id, c := range h.clients {
				if c.SID == p.sid {
					if len(p.message.Only) > 0 {
						if !util.IsExistsInStrSlice(c.UID, p.message.Only) {
							continue
						}
					} else {
						if c.UID == p.uid {
							continue
						}
					}
					select {
					case c.Send <- p.message:
					default:
						close(c.Send)
						delete(h.clients, id)
					}
				}
			}
		}
	}
}
