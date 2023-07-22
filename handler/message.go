package handler

import (
	"app/lib/ws"
	"net/http"
)

func message(w http.ResponseWriter, r *http.Request) {
	// logger := util.Logger
	q := r.URL.Query()
	roomID := q.Get("room")
	nodeID := q.Get("node")
	unsafeConn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// defer unsafeConn.Close()
	// room := pkg.RoomManager.GetRoom(roomID)
	// conn := &ws.ThreadSafeWriter{unsafeConn, sync.Mutex{}}
	c := ws.Hub.Join(nodeID, roomID, unsafeConn)
	go c.ReadPump(func(msg *ws.Message) error {
		switch msg.Event {
		case "broadcast":
			c.Hub.Broadcast(c.UID, c.SID, "broadcast", msg.Data, msg.Only, c.UID)
		}
		return nil
	})
	go c.WritePump()
}
