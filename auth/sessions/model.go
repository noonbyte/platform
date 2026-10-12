package sessions

import "time"

type SessionStatus string

var (
	SessionStatusExpired SessionStatus = "expired"
	SessionStatusRevoked SessionStatus = "revoked"
)

type Session struct {
	ID string `json:"string"`

	Status SessionStatus `json:"status"`

	UserID string `json:"user_id"`

	ExpiresAt time.Time `json:"expires_at"`
}
