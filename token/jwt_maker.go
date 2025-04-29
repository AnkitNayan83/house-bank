package token

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

const minSecretKeySize = 32

type JWTMaker struct {
	secretKey string
}

var (
	ErrExpiredToken = errors.New("token has expired")
	ErrInvalidToken = errors.New("token is invalid")
)

func NewJwtMaker(secretKey string) (Maker, error) {
	if len(secretKey) < minSecretKeySize {
		return nil, fmt.Errorf("invalid key size: must be at least %d characters", minSecretKeySize)
	}

	return &JWTMaker{secretKey}, nil
}

func (maker *JWTMaker) CreateToken(username string, duration time.Duration) (string, *Payload, error) {
	// create the payload
	payload, err := NewPayload(username, duration)
	if err != nil {
		return "", nil, err
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
		return "", nil, err
	}

	return token, payload, nil
}
func (maker *JWTMaker) VerifyToken(token string) (*Payload, error) {
	// parse the token
	jwtToken, err := jwt.Parse(token, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return []byte(maker.secretKey), nil
	})
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrExpiredToken
		}
		return nil, ErrExpiredToken
	}

	// check if the token is valid
	claims, ok := jwtToken.Claims.(jwt.MapClaims)
	if !ok || !jwtToken.Valid {
		return nil, ErrInvalidToken
	}

	// Extract and validate the ID
	idStr, ok := claims["id"].(string)
	if !ok {
		return nil, ErrInvalidToken
	}
	id, err := uuid.Parse(idStr)
	if err != nil {
		return nil, ErrInvalidToken
	}

	// Extract other claims with type checking
	username, ok := claims["username"].(string)
	if !ok {
		return nil, ErrInvalidToken
	}

	issuedAtFloat, ok := claims["issued_at"].(float64)
	if !ok {
		return nil, ErrInvalidToken
	}

	expiresAtFloat, ok := claims["expires_at"].(float64)
	if !ok {
		return nil, ErrInvalidToken
	}

	issuedAt := time.Unix(int64(issuedAtFloat), 0)
	expiresAt := time.Unix(int64(expiresAtFloat), 0)

	if time.Now().After(expiresAt) {
		return nil, ErrExpiredToken
	}

	payload := &Payload{
		ID:        id,
		Username:  username,
		IssuedAt:  issuedAt,
		ExpiresAt: expiresAt,
	}

	return payload, nil
}
