package nats

import (
	"fmt"

	"github.com/nats-io/nats.go"
)

type NATS interface {
	Publish(subject string, data []byte) error
	Subscribe(subject string) (<-chan *nats.Msg, error)
	StreamAdd(streamName string, subjects ...string) error
}

type natsClient struct {
	conn *nats.Conn
	js   nats.JetStreamContext
}

func New(url string) (NATS, error) {
	nc, err := nats.Connect(url)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to NATS: %w", err)
	}

	js, err := nc.JetStream()
	if err != nil {
		return nil, fmt.Errorf("failed to create JetStream context: %w", err)
	}

	return &natsClient{
		conn: nc,
		js:   js,
	}, nil
}
