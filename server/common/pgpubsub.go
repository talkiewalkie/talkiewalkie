package common

import (
	"encoding/json"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"log"
	"time"
)

type PgPubSub struct {
	db                   *sqlx.DB
	connInfo             string
	minReconnectInterval time.Duration
	maxReconnectInterval time.Duration
}

type PubSubEvent struct {
	Timestamp time.Time `json:"timestamp"`
	Type      string    `json:"type"`
}

type IPubSubEvent interface {
	Str() string
}

func (p PubSubEvent) Str() string {
	return "base event"
}

type NewMessageEvent struct {
	PubSubEvent
	Message      string `json:"message"`
	AuthorHandle string `json:"authorHandle"`
}

func (ps PgPubSub) Subscribe(topic string) (*pq.Listener, func() error, error) {
	listener := pq.NewListener(ps.connInfo, ps.minReconnectInterval, ps.maxReconnectInterval, func(event pq.ListenerEventType, err error) {
		var estr string
		switch event {
		case 0:
			estr = "ListenerEventConnected"
		case 1:
			estr = "ListenerEventDisconnected"
		case 2:
			estr = "ListenerEventReconnected"
		case 3:
			estr = "ListenerEventConnectionAttemptFailed"
		default:
			estr = "Unknown listener event"
		}

		log.Printf("received new event in listener created for '%s': %+v", topic, estr)
		if err != nil {
			log.Printf("could not subscribe to topic '%s': %+v", topic, err)
		}
	})

	if err := listener.Listen(topic); err != nil {
		return nil, nil, err
	}
	log.Printf("subscribed to topic '%s'", topic)

	return listener, func() error {
		_, err := ps.db.Exec(fmt.Sprintf("UNLISTEN %s", topic))
		return err
	}, nil
}

func (ps PgPubSub) Publish(topic string, event IPubSubEvent) error {
	payload, err := json.Marshal(event)
	if err != nil {
		return err
	}
	log.Printf("publishing to '%s': %+v", topic, event)
	_, err = ps.db.Exec(fmt.Sprintf("NOTIFY %s, '%s'", topic, string(payload)))
	return err
}

func NewPgPubSub(db *sqlx.DB, connInfo string) PgPubSub {
	return PgPubSub{
		db:                   db,
		connInfo:             connInfo,
		minReconnectInterval: 5 * time.Second,
		maxReconnectInterval: 10 * time.Second,
	}
}
