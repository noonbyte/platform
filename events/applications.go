package events

import (
	"encoding/json"
	"fmt"

	"github.com/noonbyte/platform/infrastructure/nats"

	"github.com/rs/zerolog/log"
)

var (
	ApplicationCreatedSubject = "applications.created"
	ApplicationUpdatedSubject = "applications.updated"
	ApplicationDeletedSubject = "applications.deleted"
)

func SetupApplicationsStreams(n nats.NATS) error {
	streams := map[string][]string{
		"APPLICATIONS": {ApplicationCreatedSubject, ApplicationUpdatedSubject, ApplicationDeletedSubject},
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

type ApplicationCreatedEventData struct {
	ID string `json:"id"`

	Name string `json:"name"`
	Type string `json:"type"`

	BundleIdentifier string `json:"bundle_identifier"`
}

type ApplicationUpdatedEventData struct {
	ID   string `json:"id"`
	Name string `json:"name"`

	Type string `json:"type"`

	BundleIdentifier string `json:"bundle_identifier"`
}

type ApplicationDeletedEventData struct {
	ID string `json:"id"`
}

func PublishApplicationCreated(n nats.NATS, data ApplicationCreatedEventData) error {
	encoded, err := json.Marshal(data)
	if err != nil {
		log.Error().Err(err).Msg("failed to marshal application created event data")
		return err
	}

	if err := n.Publish(ApplicationCreatedSubject, encoded); err != nil {
		log.Error().Err(err).Msg("failed to publish application created event")
		return err
	}
	log.Info().Msg("published application created event")
	return nil
}

func PublishApplicationUpdated(n nats.NATS, data ApplicationUpdatedEventData) error {
	encoded, err := json.Marshal(data)
	if err != nil {
		log.Error().Err(err).Msg("failed to marshal application updated event data")
		return err
	}

	if err := n.Publish(ApplicationUpdatedSubject, encoded); err != nil {
		log.Error().Err(err).Msg("failed to publish application created event")
		return err
	}
	log.Info().Msg("published application created event")
	return nil
}

func PublishApplicationDeleted(n nats.NATS, data ApplicationDeletedEventData) error {
	encoded, err := json.Marshal(data)
	if err != nil {
		log.Error().Err(err).Msg("failed to marshal application deleted event data")
		return err
	}

	if err := n.Publish(ApplicationDeletedSubject, encoded); err != nil {
		log.Error().Err(err).Msg("failed to publish application deleted event")
		return err
	}
	log.Info().Msg("published application deleted event")
	return nil
}
