package websockets

import (
	"encoding/json"
	"fmt"
	"github.com/lib/pq"
	"github.com/talkiewalkie/talkiewalkie/common"
	"github.com/talkiewalkie/talkiewalkie/entities"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

type WsMessageOut struct {
	Text         string `json:"text"`
	AuthorHandle string `json:"authorHandle"`
}

func websocketSync(c *websocket.Conn, closing chan bool) {
	defer func() {
		if err := c.Close(); err != nil {
			log.Printf("%s failed to close websocket: %+v", time.Now(), err)
		}
	}()
	clientPingInterval := 5 * time.Second
	c.SetPongHandler(func(string) error { c.SetReadDeadline(time.Now().Add(clientPingInterval)); return nil })

	for {
		_, msg, err := c.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("ws error: %v", err)
			}
			closing <- true
			break
		}
		// TODO: no use case for reading from the ws at the moment, will probly port the /message handler in the
		// 	ws in the near future though
		log.Printf("DEBUG: received ws message: %s", string(msg))
	}
}

func pubsubSync(c *websocket.Conn, topic string, listener *pq.Listener, unlisten func() error, closing chan bool) {
	defer func() {
		if err := unlisten(); err != nil {
			log.Printf("could not stop listening to topic '%s': %+v", topic, err)
		}
	}()

	for {
		select {
		case <-closing:
			log.Printf("websocket closing, stopping pubsub sync goroutine (topic=%s)", topic)
			return

		case psEvent := <-listener.Notify:
			var payload common.PubSubEvent
			if err := json.Unmarshal([]byte(psEvent.Extra), &payload); err != nil {
				log.Printf("failed to parse pubsub payload on topic '%s': %s", topic, psEvent.Extra)
			}

			switch payload.Type {
			case "newmessage":
				var msg common.NewMessageEvent
				_ = json.Unmarshal([]byte(psEvent.Extra), &msg)
				c.WriteMessage(websocket.TextMessage, []byte(psEvent.Extra))
			default:
				log.Printf("received unknown pubsub message on topic '%s': (%T) %+v", topic, payload, payload)
			}
		}
	}
}

func GroupWebsocketHandler(w http.ResponseWriter, r *http.Request) {
	ctx := common.WithAuthedContext(r)

	conn := websocket.Upgrader{HandshakeTimeout: time.Second}
	c, err := conn.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("could not initiate websocket: %+v", err)
		http.Error(w, fmt.Sprintf("could not initiate websocket: %+v", err), http.StatusInternalServerError)
		c.Close()
		return
	}

	topic := entities.UserPubSubTopic(ctx.User)
	log.Printf("established websocket connection [%s]", topic)

	listener, unlisten, err := ctx.Components.PgPubSub.Subscribe(topic)
	if err != nil {
		http.Error(w, fmt.Sprintf("could not subscribe to pubsub topic: %+v", err), http.StatusBadRequest)
		c.Close()
		if err := unlisten(); err != nil {
			log.Printf("could not stop listening to topic '%s': %+v", topic, err)
		}
		return
	}

	closing := make(chan bool)
	go websocketSync(c, closing)
	go pubsubSync(c, topic, listener, unlisten, closing)
}
