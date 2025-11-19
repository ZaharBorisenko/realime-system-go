package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
)

type Client struct {
	//A websocket for this user
	Socket *websocket.Conn

	//receive is a channel to receive messages from other clients
	Receive chan []byte

	Room *Room
	Name string
}

func (c *Client) Read() {
	defer c.Socket.Close()

	for {
		_, msg, err := c.Socket.ReadMessage()

		if err != nil {
			fmt.Errorf("couldn't READ the message client %w", err)
			return
		}

		outgo := map[string]string{
			"name":    c.Name,
			"message": string(msg),
		}

		jsMessage, err := json.Marshal(outgo)
		if err != nil {
			fmt.Errorf("encoding failed %w", err)
			continue
		}

		c.Room.Forward <- jsMessage
	}
}

func (c *Client) Write() {
	defer c.Socket.Close()

	for msg := range c.Receive {
		err := c.Socket.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			fmt.Errorf("couldn't WRITE the message client %w", err)
			return
		}
	}

}
