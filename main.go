// This code as it currently stands is based off
// https://github.com/garyburd/gary.burd.info/tree/master/go-websocket-chat
//
// Naturally, this will change as CPSSDChat becomes more advanced.

package main

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type room struct {
	// All active connections
	connections map[*connection]bool

	// A buffered channel of incoming messages. []byte and not string
	// becaus the websockets use []bytes.
	messages chan []byte

	// Requests to join the room
	register chan *connection

	// Requests to leave the room
	unregister chan *connection
}

func newRoom() *room {
	return &room{
		connections: make(map[*connection]bool),
		messages:    make(chan []byte, 20),
		register:    make(chan *connection),
		unregister:  make(chan *connection),
	}
}

func (r *room) run() {
	for {
		select {
		case conn := <-r.register:
			r.connections[conn] = true
		case conn := <-r.unregister:
			if _, ok := r.connections[conn]; ok {
				delete(r.connections, conn)
			}
		// When receiving a message, send it down all active connections
		case message := <-r.messages:
			for conn := range r.connections {
				err := conn.ws.WriteMessage(websocket.TextMessage, message)
				if err != nil {
					log.Printf("ERROR: Could not write message: %s\n", err)
				}
			}
		}
	}
}

type connection struct {
	ws *websocket.Conn
	r  *room
}

func (conn *connection) reader() {
	for {
		_, message, err := conn.ws.ReadMessage()
		if err != nil {
			log.Printf("ERROR: Could not read message: %s\n", err)
			break
		}
		conn.r.messages <- message
	}
	conn.ws.Close()
}

type websocketHandler struct {
	r *room
}

func (handler websocketHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("ERROR: Could not upgrade websocket: %s\n", err)
	}
	conn := &connection{ws: ws, r: handler.r}
	handler.r.register <- conn
	defer func() { handler.r.unregister <- conn }()
}

func main() {
	r := newRoom()
	go r.run()
	http.Handle("/connect", websocketHandler{r: r})
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalln("FATAL: ", err)
	}
}
