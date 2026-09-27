package nats

import (
	"fmt"

	"github.com/nats-io/nats.go"
)

func (n *natsClient) StreamAdd(streamName string, subjects ...string) error {
	_, err := n.js.AddStream(&nats.StreamConfig{
		Name:     streamName,
		Subjects: subjects,
		Storage:  nats.FileStorage,
	})
	if err != nil && err != nats.ErrStreamNameAlreadyInUse {
		return fmt.Errorf("failed to create stream: %w", err)
	}
	return nil
}
