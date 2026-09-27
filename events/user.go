package events

import (
	"encoding/json"
	"fmt"

	"github.com/noonbyte/platform/infrastructure/nats"
	"github.com/rs/zerolog/log"
)

var (
	UserCreatedSubject = "users.created"
	UserUpdatedSubject = "users.updated"
	UserDeletedSubject = "users.deleted"
)

func SetupUserStreams(n nats.NATS) error {
	streams := map[string][]string{
		"USERS": {UserCreatedSubject, UserUpdatedSubject, UserDeletedSubject},
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

type UserCreatedEventData struct {
	ID string `json:"id"`

	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserUpdatedEventData struct {
	ID string `json:"id"`

	Email    *string `json:"email"`
	Username *string `json:"username"`
	Password *string `json:"password"`
}

type UserDeletedEventData struct {
	ID string `json:"id"`
}

func PublishUserCreated(n nats.NATS, data UserCreatedEventData) error {
	encoded, err := json.Marshal(data)
	if err != nil {
		log.Error().Err(err).Msg("failed to marshal user created event data")
		return err
	}

	if err := n.Publish(UserCreatedSubject, encoded); err != nil {
		log.Error().Err(err).Msg("failed to publish user created event")
		return err
	}
	log.Info().Msg("published user created event")
	return nil
}

func PublishUserUpdated(n nats.NATS, data UserUpdatedEventData) error {
	encoded, err := json.Marshal(data)
	if err != nil {
		log.Error().Err(err).Msg("failed to marshal user updated event data")
		return err
	}

	if err := n.Publish(UserUpdatedSubject, encoded); err != nil {
		log.Error().Err(err).Msg("failed to publish user updated event")
		return err
	}
	log.Info().Msg("published user updated event")
	return nil
}

func PublishUserDeleted(n nats.NATS, data UserDeletedEventData) error {
	encoded, err := json.Marshal(data)
	if err != nil {
		log.Error().Err(err).Msg("failed to marshal user deleted event data")
		return err
	}

	if err := n.Publish(UserDeletedSubject, encoded); err != nil {
		log.Error().Err(err).Msg("failed to publish user deleted event")
		return err
	}
	log.Info().Msg("published user deleted event")
	return nil
}
