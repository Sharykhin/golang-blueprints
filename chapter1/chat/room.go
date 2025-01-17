package main

import (
	"github.com/Sharykhin/blueprint/chapter1/trace"
	"github.com/gorilla/websocket"
	"github.com/stretchr/objx"
	"log"
	"net/http"
)

const (
	socketBufferSize  = 1024
	messageBufferSize = 256
)

var upgrader = &websocket.Upgrader{ReadBufferSize: socketBufferSize, WriteBufferSize: socketBufferSize}

type (
	room struct {
		// forward is a channel that holds incoming messages
		// that should be forwarded to the other clients.
		forward chan *message
		// join is a channel for clients wishing to join the room.
		join chan *client
		// leave is a channel for clients wishing to leave the room.
		leave chan *client
		// clients holds all current clients in this room.
		clients map[*client]bool
		// tracer will receive trace information of activity
		// in the room.
		tracer trace.Tracer
	}
)

func (r *room) run() {
	for {
		select {
		case client := <-r.join:
			r.tracer.Trace("New client joined")
			r.clients[client] = true
		case client := <-r.leave:
			delete(r.clients, client)
			close(client.send)
			r.tracer.Trace("Client left")
		case msg := <-r.forward:
			r.tracer.Trace("Message received: ", msg.Message)
			for client := range r.clients {
				client.send <- msg
				r.tracer.Trace(" -- sent to client")
			}
		}
	}
}

func (r *room) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	socket, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		log.Fatal("ServerHTTP", err)
		return
	}

	authCookie, err := req.Cookie("auth")
	if err != nil {
		log.Fatal("Failed to get auth cookie", err)
	}

	client := &client{
		socket:   socket,
		send:     make(chan *message, messageBufferSize),
		room:     r,
		userData: objx.MustFromBase64(authCookie.Value),
	}

	r.join <- client
	defer func() {
		r.leave <- client
	}()
	go client.write()
	client.read()
}

func newRoom() *room {
	return &room{
		forward: make(chan *message),
		join:    make(chan *client),
		leave:   make(chan *client),
		clients: make(map[*client]bool),
		tracer:  trace.Off(),
	}
}
