package token

import (
	"time"

	"github.com/google/uuid"
)

type Payload struct {
	ID        uuid.UUID `json:"id"` // this field will be use to validate leaked tokens
	Username  string    `json:"username"`
	IssuedAt  time.Time `json:"issued_at"`
	ExpiresAt time.Time `json:"expires_at"`
}

func NewPayload(username string, duration time.Duration) (*Payload, error) {
	tokenId, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	payload := &Payload{
		ID:        tokenId,
		Username:  username,
		IssuedAt:  time.Now(),
		ExpiresAt: time.Now().Add(duration),
	}

	return payload, nil
}
