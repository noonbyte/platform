package nats

import (
	"fmt"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/nats-io/nats.go"
)

func (n *natsClient) Subscribe(
	subject string,
	consumerName string,
) (<-chan *nats.Msg, error) {
	ch := make(chan *nats.Msg)

	sub, err := n.js.PullSubscribe(subject, consumerName)
	if err != nil {
		return nil, fmt.Errorf("failed to subscribe: %w", err)
	}

	go func() {
		defer close(ch)

		for {
			msgs, err := sub.Fetch(10, nats.MaxWait(time.Second))
			if err != nil {
				if err == nats.ErrTimeout {
					continue
				}

				log.Error().
					Err(err).
					Str("subject", subject).
					Str("consumer", consumerName).
					Msg("NATS fetch error")

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
