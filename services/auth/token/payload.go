package token

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// Different types of error returned by the VerifyToken function
var (
	InvalidTokenError = fmt.Errorf("invalid token")
	ExpiredTokenError = fmt.Errorf("expired token")
)

// Payload is the payload data that is stored in the token
type Payload struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	IssuedAt  time.Time `json:"issued_at"`
	ExpiredAt time.Time `json:"expired_at"`
}

// NewPayload creates a new payload with the given username and duration
func NewPayload(username string, duration time.Duration) (*Payload, error) {
	tokenID, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}
	payload := &Payload{
		ID:        tokenID,
		Username:  username,
		IssuedAt:  time.Now(),
		ExpiredAt: time.Now().Add(duration),
	}
	return payload, nil
}

// Valid checks if the payload is valid
func (payload *Payload) Valid() error {
	if payload.IssuedAt.Before(payload.ExpiredAt) {
		return nil
	}
	return ExpiredTokenError
}
