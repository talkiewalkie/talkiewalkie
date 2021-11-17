package clients

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"log"
	"time"
)

type PubSubClient interface {
	Subscribe(string) (chan *pq.Notification, func(), error)
	Publish(string, proto.Message) error
}

type PgPubSub struct {
	db                   *sqlx.DB
	connInfo             string
	minReconnectInterval time.Duration
	maxReconnectInterval time.Duration
}

var _ PubSubClient = PgPubSub{}

func (ps PgPubSub) Subscribe(topic string) (chan *pq.Notification, func(), error) {
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

	return listener.Notify, func() {
		if _, err := ps.db.Exec(fmt.Sprintf("UNLISTEN %s", topic)); err != nil {
			log.Printf("ERR: failed to cancel subscription to topic[%s]: %+v", topic, err)
		} else if err := listener.Close(); err == nil {
			log.Printf("ERR: failed to close topic[%s] listener: %+v", topic, err)
		}
	}, nil
}

func (ps PgPubSub) Publish(topic string, event proto.Message) error {
	payload, err := protojson.Marshal(event)
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
