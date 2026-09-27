package nats

import (
	"fmt"
	"strings"
	"time"

	"github.com/nats-io/nats.go"
)

func (n *natsClient) Subscribe(subject string) (<-chan *nats.Msg, error) {
	ch := make(chan *nats.Msg)

	sub, err := n.js.PullSubscribe(subject, sanitizeConsumerName(subject))
	if err != nil {
		return nil, fmt.Errorf("failed to subscribe: %w", err)
	}

	go func() {
		defer close(ch)
		for {
			msgs, err := sub.Fetch(10, nats.MaxWait(1*time.Second))
			if err != nil && err != nats.ErrTimeout {
				fmt.Printf("NATS fetch error: %v\n", err)
				time.Sleep(time.Second)
				continue
			}

			for _, msg := range msgs {
				ch <- msg
			}
		}
	}()

	return ch, nil
}

func sanitizeConsumerName(subject string) string {
	name := strings.ReplaceAll(subject, ".", "_")
	return name + "_consumer"
}

func AckMessages(msgs ...*nats.Msg) {
	for _, msg := range msgs {
		_ = msg.Ack()
	}
}
