package main

import (
	"github.com/gorilla/websocket"
)

// Client represents a single chatting user
type client struct {
	socket *websocket.Conn
	// send is a channel on which messages are sent
	send chan []byte

	// room is the room this client is chatting in
	room *room
}

func (c *client) read() {
	for {
		if _, msg, err := c.socket.ReadMessage(); err == nil {
			c.room.forward <- msg
		} else {
			println("Breaking!")
			break
		}
	}
	println("closing!")
	c.socket.Close()
}

func (c *client) write() {
	for msg := range c.send {
		if err := c.socket.WriteMessage(websocket.TextMessage, msg); err != nil {
			break
		}
	}
	c.socket.Close()
}
