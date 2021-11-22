package clients

import (
	"fmt"
	"log"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"

	"github.com/talkiewalkie/talkiewalkie/pb"
)

type PubSubClient interface {
	Subscribe(string) (chan *pq.Notification, func(), error)
	Publish(string, *pb.Event) error
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

			PubSubLogf(topic, "%s", estr)
			if err != nil {
				PubSubLogf(topic, "could not subscribe: %+v", err)
			}
		},
	)

	if err := listener.Listen(topic); err != nil {
		return nil, nil, err
	}

	return listener.Notify, func() {
		if _, err := ps.db.Exec(fmt.Sprintf("UNLISTEN %s", topic)); err != nil {
			PubSubLogf(topic, "ERR: failed to cancel subscription: %+v", err)
		} else if err := listener.Close(); err == nil {
			PubSubLogf(topic, "ERR: failed to close listener: %+v", err)
		}
	}, nil
}

func (ps PgPubSub) Publish(topic string, event *pb.Event) error {
	PubSubLogf(topic, "pushing [%T] with local uuid set to '%s'", event.Content, event.LocalUuid)
	_, err := ps.db.Exec(fmt.Sprintf("NOTIFY %s, '%s'", topic, event.Uuid))
	if err != nil {
		return fmt.Errorf("could not notify of event in topic '%s': %+v", topic, err)
	} else {
		return err
	}
}

func PubSubLogf(topic, msg string, args ...interface{}) {
	log.Printf("[pubsub:%s] %s", topic, fmt.Sprintf(msg, args...))
}

func NewPgPubSub(db *sqlx.DB, connInfo string) PgPubSub {
	return PgPubSub{
		db:                   db,
		connInfo:             connInfo,
		minReconnectInterval: 5 * time.Second,
		maxReconnectInterval: 10 * time.Second,
	}
}
