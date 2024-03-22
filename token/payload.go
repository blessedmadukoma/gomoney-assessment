package token

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

var ErrExpiredToken = errors.New("token has expired")

// Payload contains the payload data of the token
type Payload struct {
	ID        uuid.UUID `json:"id"`
	Email     string    `json:"email"`
	IssuedAt  time.Time `json:"issued_at"`
	ExpiredAt time.Time `json:"expired_at"`
}

// NewPayload creates a new token payload with a specific email and duration
func NewTokenPayload(email string, duration time.Duration) (*Payload, error) {
	tokenID, err := uuid.NewRandom()
	if err != nil {
		return nil, fmt.Errorf("Error creating token ID: %s", err)
	}

	return &Payload{
		ID:        tokenID,
		Email:     email,
		IssuedAt:  time.Now(),
		ExpiredAt: time.Now().Add(duration),
	}, nil
}

// Valid checks if the token payload is valid or not
func (p *Payload) Valid() error {
	if time.Now().After(p.ExpiredAt) {
		return ErrExpiredToken
	}
	return nil
}
