package main

import (
	"github.com/gorilla/websocket"
	"net/http"
)

type room struct {
	// forward is a channel that holds incoming messages
	// that should be forwarded to the other clients
	forward chan []byte

	// join is a channel for clients wishing to join the room
	join chan *client

	// leave is a channel for clients wishing to leave the room
	leave chan *client
	// clients holds all current cliens in this room
	clients map[*client]bool
}

// newRoom makes a new room that is ready to go
func newRoom() *room {
	return &room{
		forward: make(chan []byte),
		join:    make(chan *client),
		leave:   make(chan *client),
		clients: make(map[*client]bool),
	}
}

func (r *room) run() {
	println("let's run this room!")
	for {
		select {
		case client := <-r.join:
			println("wtf?")
			r.clients[client] = true
		case client := <-r.leave:
			println("leaving?")
			delete(r.clients, client)
			close(client.send)
		case msg := <-r.forward:
			println("msg?")
			for client := range r.clients {
				select {
				case client.send <- msg:
				default:
					// failed to send
					delete(r.clients, client)
					close(client.send)
				}
			}
		}
	}
}

const (
	socketBufferSize  = 1024
	messageBufferSize = 256
)

var upgrader = &websocket.Upgrader{ReadBufferSize: socketBufferSize, WriteBufferSize: socketBufferSize}

func (r *room) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	println("Upgrading!``")
	socket, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		return
	}

	client := &client{socket: socket, send: make(chan []byte, messageBufferSize), room: r}
	r.join <- client
	defer func() { r.leave <- client }()
	go client.write()
	client.read()
}
