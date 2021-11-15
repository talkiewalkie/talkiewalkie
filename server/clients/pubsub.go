package clients

import (
	"encoding/json"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	uuid2 "github.com/satori/go.uuid"
	"log"
	"time"
)

type PubSubClient interface {
	Subscribe(string) (chan *pq.Notification, func() error, error)
	Publish(string, IPubSubEvent) error
}

type PgPubSub struct {
	db                   *sqlx.DB
	connInfo             string
	minReconnectInterval time.Duration
	maxReconnectInterval time.Duration
}

var _ PubSubClient = PgPubSub{}

type PubSubEventType int

const (
	PubSubEventTypeNewMessage PubSubEventType = iota // 0
)

type PubSubEvent struct {
	Timestamp time.Time       `json:"timestamp"`
	Type      PubSubEventType `json:"type"`
}

type IPubSubEvent interface {
	Str() string
}

func (p PubSubEvent) Str() string {
	return "base event"
}

type NewMessageEvent struct {
	PubSubEvent
	MessageUuid uuid2.UUID `json:"message_uuid"`
}

func (ps PgPubSub) Subscribe(topic string) (chan *pq.Notification, func() error, error) {
	listener := pq.NewListener(
		ps.connInfo,
		ps.minReconnectInterval,
		ps.maxReconnectInterval,
		func(event pq.ListenerEventType, err error) {
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
		},
	)

	if err := listener.Listen(topic); err != nil {
		return nil, nil, err
	}
	log.Printf("subscribed to topic '%s'", topic)

	return listener.Notify, func() error {
		_, err := ps.db.Exec(fmt.Sprintf("UNLISTEN %s", topic))
		if err2 := listener.Close(); err == nil {
			return err2
		}
		return err
	}, nil
}

func (ps PgPubSub) Publish(topic string, event IPubSubEvent) error {
	payload, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("could not serialize event: %+v", err)
	}
	log.Printf("publishing to '%s': %+v", topic, event)
	_, err = ps.db.Exec(fmt.Sprintf("NOTIFY %s, '%s'", topic, string(payload)))
	if err != nil {
		return fmt.Errorf("could not notify of event in topic '%s': %+v", topic, err)
	} else {
		return err
	}
}

func NewPgPubSub(db *sqlx.DB, connInfo string) PgPubSub {
	return PgPubSub{
		db:                   db,
		connInfo:             connInfo,
		minReconnectInterval: 5 * time.Second,
		maxReconnectInterval: 10 * time.Second,
	}
}