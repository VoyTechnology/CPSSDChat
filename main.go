// This code as it currently stands is based off
// https://github.com/garyburd/gary.burd.info/tree/master/go-websocket-chat
//
// Naturally, this will change as CPSSDChat becomes more advanced.

package main

import (
	"github.com/gorilla/websocket"
	"log"
)

type room struct {
	// All active connections
	connections map[*websocket.Conn]bool

	// A buffered channel of incoming messages
	messages chan string

	// Requests to join the room
	register chan *websocket.Conn

	// Requests to leave the room
	unregister chan *websocket.Conn
}

func newRoom() *room {
	return &room{
		connections: make(map[*websocket.Conn]bool),
		messages:    make(chan string, 20),
		register:    make(chan *websocket.Conn),
		unregister:  make(chan *websocket.Conn),
	}
}

func (r *room) run() {
	for {
		select {
		case conn := <-register:
			r.connections[conn] = true
		case conn := <-unregister:
			if _, ok := r.connections[c]; ok {
				delete(r.connections, c)
			}
		// When receiving a message, send it down all active connections
		case message := <-messages:
			for conn := range r.connections {
				err := conn.WriteMessage(websocket.TextMessage, []byte(message))
				if err != nil {
					log.Printf("ERROR: Could not write message: %s\n", err)
				}
			}
		}
	}
}

func main() {
	log.Fatalln("I do nothing!")
}
