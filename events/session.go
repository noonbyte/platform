package events

import (
	"encoding/json"
	"fmt"
	"noonbyte/platform/infrastructure/nats"

	"github.com/rs/zerolog/log"
)

var (
	SessionCreatedSubject = "sessions.created"
	SessionDeletedSubject = "sessions.deleted"
)

func SetupSessionStreams(n nats.NATS) error {
	streams := map[string][]string{
		"SESSIONS": {SessionCreatedSubject, SessionDeletedSubject},
	}

	var errors []error
	for stream, subjects := range streams {
		if err := n.StreamAdd(stream, subjects...); err != nil {
			log.Error().Err(err).Str("stream", stream).Msg("failed to declare stream")
			errors = append(errors, err)
		} else {
			log.Info().Str("stream", stream).Msg("successfully declared stream")
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("failed to declare some streams: %v", errors)
	}
	return nil
}

type SessionCreatedEventData struct {
	ID     string `json:"id"`
	Device any    `json:"device"`
}

type SessionDeletedEventData struct {
	ID string `json:"id"`
}

func PublishSessionCreated(n nats.NATS, data SessionCreatedEventData) error {
	encoded, err := json.Marshal(data)
	if err != nil {
		log.Error().Err(err).Msg("failed to marshal session created event data")
		return err
	}

	if err := n.Publish(SessionCreatedSubject, encoded); err != nil {
		log.Error().Err(err).Msg("failed to publish session created event")
		return err
	}
	log.Info().Msg("published session created event")
	return nil
}

func PublishSessionDeleted(n nats.NATS, data SessionDeletedEventData) error {
	encoded, err := json.Marshal(data)
	if err != nil {
		log.Error().Err(err).Msg("failed to marshal session deleted event data")
		return err
	}

	if err := n.Publish(SessionDeletedSubject, encoded); err != nil {
		log.Error().Err(err).Msg("failed to publish session deleted event")
		return err
	}
	log.Info().Msg("published session deleted event")
	return nil
}
