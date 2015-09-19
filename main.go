// This code as it currently stands is based off
// https://github.com/garyburd/gary.burd.info/tree/master/go-websocket-chat
//
// Naturally, this will change as CPSSDChat becomes more advanced.

package main

import (
	"github.com/gorilla/websocket"
	"github.com/voytechnology/cpssdchat/chat"
	"log"
	"net/http"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type websocketHandler struct {
	r *chat.Room
}

func (handler websocketHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("ERROR: Could not upgrade websocket: %s\n", err)
		return
	}
	conn := chat.NewConnection(ws, handler.r)
	handler.r.Register <- conn
	defer func() { handler.r.Unregister <- conn }()
	go conn.Writer()
	conn.Reader()
}

func main() {
	r := chat.NewRoom()
	go r.Run()
	http.Handle("/connect", websocketHandler{r: r})
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalln("FATAL: ", err)
	}
}
