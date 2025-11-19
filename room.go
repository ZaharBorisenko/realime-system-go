package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

type Room struct {

	// holds all current clients in the room
	Clients map[*Client]bool

	// join is a channel for all clients wishing to join the room
	Join chan *Client

	// leave is a channel for all clients wishing to leave the room
	Leave chan *Client

	// forward is a channel that holds incoming messages that should be forwarded to the other clients.

	Forward chan []byte
}

func NewRoom() *Room {
	return &Room{
		Forward: make(chan []byte),
		Join:    make(chan *Client),
		Leave:   make(chan *Client),
		Clients: make(map[*Client]bool),
	}
}

func (r *Room) Run() {
	for {
		select {
		case client := <-r.Join:
			r.Clients[client] = true
		case client := <-r.Leave:
			delete(r.Clients, client)
			close(client.Receive)
		case msg := <-r.Forward:
			for client := range r.Clients {
				client.Receive <- msg
			}
		}
	}
}

const (
	socketBufferSize  = 1024
	messageBufferSize = 256
)

var upgrader = &websocket.Upgrader{ReadBufferSize: socketBufferSize, WriteBufferSize: socketBufferSize}

var rooms = make(map[string]*Room)
var mu sync.Mutex

func getRoom(name string) *Room {
	mu.Lock()
	defer mu.Unlock()
	if r, ok := rooms[name]; ok {
		return r
	}
	r := NewRoom()
	rooms[name] = r
	go r.Run()
	return r
}

func (r *Room) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	roomName := req.URL.Query().Get("room")
	userName := req.URL.Query().Get("username")
	if roomName == "" {
		http.Error(w, "Room name required", http.StatusBadRequest)
		return
	}
	if userName == "" {
		userName = fmt.Sprintf("user_%d", rand.Intn(1000))
		log.Println("given random name")
	}

	realRoom := getRoom(roomName)

	socket, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}

	client := &Client{
		Socket:  socket,
		Receive: make(chan []byte, messageBufferSize),
		Room:    realRoom,
		Name:    userName,
	}

	realRoom.Join <- client
	defer func() { realRoom.Leave <- client }()
	go client.Write()
	client.Read()
}
