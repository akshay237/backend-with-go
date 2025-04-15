package token

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

const minSecretKeySize = 32

type UserClaims struct {
	ID       uuid.UUID `json:"id"`
	Username string    `json:"username"`
	*jwt.RegisteredClaims
}

// JWTMaker is a JSON Web Token maker
type JWTMaker struct {
	secretKey string
}

// NewJWTMaker creates a new JWTMaker
func NewJWTMaker(secretKey string) (Maker, error) {
	if len(secretKey) < minSecretKeySize {
		return nil, fmt.Errorf("invalid key size: must be at least %d characters", minSecretKeySize)
	}

	return &JWTMaker{secretKey: secretKey}, nil
}

// CreateToken creates a new token for a specific username and duration
func (m *JWTMaker) CreateToken(username string, duration time.Duration) (string, error) {

	// 1. Create the token payload
	payload, err := NewPayload(username, duration)
	if err != nil {
		return "", err
	}

	// 2. create the claims and jwt token
	claims := &UserClaims{
		ID:       payload.ID,
		Username: payload.Username,
		RegisteredClaims: &jwt.RegisteredClaims{
			Issuer:    "Simple Bank",
			IssuedAt:  jwt.NewNumericDate(payload.IssuedAt),
			ExpiresAt: jwt.NewNumericDate(payload.ExpiredAT),
		},
	}
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 3. return the signed token
	return jwtToken.SignedString([]byte(m.secretKey))
}

// VerifyToken checks if the token is valid or not
func (m *JWTMaker) VerifyToken(token string) (*Payload, error) {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, ErrInvalidToken
		}
		return []byte(m.secretKey), nil
	}

	jwtToken, err := jwt.ParseWithClaims(token, &UserClaims{}, keyFunc)
	if err != nil {
		return nil, ErrInvalidToken
	}

	claims, ok := jwtToken.Claims.(*UserClaims)
	if !ok {
		return nil, ErrInvalidToken
	}
	return &Payload{
		ID:        claims.ID,
		Username:  claims.Username,
		IssuedAt:  claims.IssuedAt.Local(),
		ExpiredAT: claims.ExpiresAt.Local(),
	}, nil
}
