package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

func read(c *websocket.Conn) {
	defer c.Close()

	for {
		_, msg, err := c.ReadMessage()
		if err != nil {
			log.Println(err)
			break
		}

		log.Println("received", string(msg))
		c.SetWriteDeadline(time.Now().Add(10 * time.Second))
		c.WriteMessage(websocket.TextMessage, []byte("hey you sent a message"))
	}
}

func ws(w http.ResponseWriter, r *http.Request) {
	conn := websocket.Upgrader{
		HandshakeTimeout:  0,
		ReadBufferSize:    0,
		WriteBufferSize:   0,
		WriteBufferPool:   nil,
		Subprotocols:      nil,
		Error:             nil,
		CheckOrigin:       nil,
		EnableCompression: false,
	}
	c, err := conn.Upgrade(w, r, nil)
	if err != nil {
		log.Println("could not initiate websocket", err)
	}

	go read(c)
}
