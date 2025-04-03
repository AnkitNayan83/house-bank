package token

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const minSecretKeySize = 32

type JWTMaker struct {
	secretKey string
}

func NewJwtMaker(secretKey string) (Maker, error) {
	if len(secretKey) < minSecretKeySize {
		return nil, fmt.Errorf("invalid key size: must be at least %d characters", minSecretKeySize)
	}

	return &JWTMaker{secretKey}, nil
}

func (maker *JWTMaker) CreateToken(username string, duration time.Duration) (string, error) {
	// create the payload
	payload, err := NewPayload(username, duration)
	if err != nil {
		return "", err
	}

	// create the token
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":         payload.ID,
		"username":   payload.Username,
		"issued_at":  payload.IssuedAt,
		"expires_at": payload.ExpiresAt,
	})

	// sign the token
	token, err := jwtToken.SignedString([]byte(maker.secretKey))
	if err != nil {
		return "", err
	}

	return token, nil
}

func (maker *JWTMaker) VerifyToken(token string) (*Payload, error) {

}
