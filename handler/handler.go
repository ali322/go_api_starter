package handler

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin:     func(r *http.Request) bool { return true },
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var pongWait = 60 * time.Second
var pingPeriod = (pongWait * 9) / 10

func withCors(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
		w.Header().Set("Access-Control-Allow-Headers", "Access-Control-Allow-Origin, Origin, Content-Type, Content-Length, Accept-Encoding, X-Requested-With, X-CSRF-Token, Authorization, X-NT-Captcha")
		w.Header().Set("Access-Control-Expose-Headers", "Content-Length")
		if r.Method == http.MethodOptions {
			return
		}
		handler(w, r)
	}
}

func pong(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "pong")
}

func ApplyRoutes() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/ping", pong).Methods(http.MethodGet)
	// r.HandleFunc("/record", record).Methods(http.MethodGet)
	r.Use(mux.CORSMethodMiddleware(r))
	r.HandleFunc("/message", message).Methods(http.MethodGet)
	// mux.HandleFunc("/meeting/positive", sfu.MeetingPositive)
	// mux.HandleFunc("/live", sfu.LiveStream)
	// mux.HandleFunc("/sig", withCors(sig.Handler))
	return r
}
