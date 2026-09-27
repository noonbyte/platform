package nats

import (
	"fmt"
	"time"

	"github.com/nats-io/nats.go"
)

func (n *natsClient) Publish(subject string, data []byte) error {
	ack, err := n.js.Publish(subject, data, nats.MsgId(fmt.Sprintf("%d", time.Now().UnixNano())))
	if err != nil {
		return fmt.Errorf("failed to publish message: %w", err)
	}

	fmt.Printf("Published message to %s, seq: %d\n", subject, ack.Sequence)
	return nil
}
