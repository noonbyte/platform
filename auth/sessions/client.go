package sessions

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type SessionClient struct {
	client *http.Client

	url string
}

func NewSessionClient(url string) *SessionClient {
	return &SessionClient{
		client: &http.Client{},

		url: url,
	}
}

func (c *SessionClient) Get(
	ctx context.Context,
	sessionID string,
) (*Session, error) {
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		fmt.Sprintf("%s/%s", c.url, sessionID),
		nil,
	)

	if err != nil {
		return nil, err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, ErrInvalidSession
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(
			"session service returned status %d",
			resp.StatusCode,
		)
	}

	var session Session

	if err := json.NewDecoder(resp.Body).Decode(&session); err != nil {
		return nil, fmt.Errorf("decode session: %w", err)
	}

	return &session, nil
}
