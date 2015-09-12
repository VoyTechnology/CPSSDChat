// This code as it currently stands is based off
// https://github.com/garyburd/gary.burd.info/tree/master/go-websocket-chat
//
// Naturally, this will change as CPSSDChat becomes more advanced.

package chat

import (
	"github.com/gorilla/websocket"
	"log"
)

type Room struct {
	// All active connections
	connections map[*Connection]bool

	// A buffered channel of incoming messages. []byte and not string
	// becaus the websockets use []bytes.
	messages chan []byte

	// Requests to join the room
	Register chan *Connection

	// Requests to leave the room
	Unregister chan *Connection
}

func NewRoom() *Room {
	return &Room{
		connections: make(map[*Connection]bool),
		messages:    make(chan []byte, 20),
		Register:    make(chan *Connection),
		Unregister:  make(chan *Connection),
	}
}

func (r *Room) Run() {
	for {
		select {
		case conn := <-r.Register:
			r.connections[conn] = true
		case conn := <-r.Unregister:
			if _, ok := r.connections[conn]; ok {
				delete(r.connections, conn)
			}
		// When receiving a message, send it down all active connections
		case message := <-r.messages:
			for conn := range r.connections {
				conn.mess <- message
			}
		}
	}
}

type Connection struct {
	ws *websocket.Conn
	r  *Room
	// Buffered channel of messages to be sent to the client
	mess chan []byte
}

func NewConnection(ws *websocket.Conn, r *Room) *Connection {
	return &Connection{
		ws:   ws,
		r:    r,
		mess: make(chan []byte, 256),
	}
}

func (conn *Connection) Reader() {
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

func (conn *Connection) Writer() {
	for message := range conn.mess {
		err := conn.ws.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			log.Println("ERROR: Could not write to websocket: ", err)
			break
		}
	}
	conn.ws.Close()
}
