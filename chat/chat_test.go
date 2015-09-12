package chat

import (
	"testing"
)

func TestConnection(t *testing.T) {
	r := NewRoom()
	go r.Run()
	conn := &Connection{
		ws:   nil,
		r:    r,
		mess: make(chan []byte, 256),
	}
	r.Register <- conn
	r.messages <- []byte("test")
	if m := <-conn.mess; string(m) != "test" {
		t.Error("Message not received by connection from room.")
	}
	r.Unregister <- conn
}
