package handler

import (
	"app/middleware"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"net/http"
	"time"
)

var upgrader = websocket.Upgrader{
	CheckOrigin:     func(r *http.Request) bool { return true },
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var pongWait = 60 * time.Second
var pingPeriod = (pongWait * 9) / 10

func pong(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "pong")
}

func ApplyRoutes() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/ping", pong).Methods(http.MethodGet)
	// r.HandleFunc("/record", record).Methods(http.MethodGet)
	r.Use(middleware.Logger, middleware.Recover, middleware.Cors)
	r.HandleFunc("/message", message).Methods(http.MethodGet)
	// mux.HandleFunc("/meeting/positive", sfu.MeetingPositive)
	// mux.HandleFunc("/live", sfu.LiveStream)
	// mux.HandleFunc("/sig", withCors(sig.Handler))
	return r
}
