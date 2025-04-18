package token

import (
	"fmt"
	"time"

	"golang.org/x/crypto/chacha20poly1305"

	pasteo "github.com/o1egl/paseto"
)

// PasteoMaker is a PASTEO token maker.
type PasteoMaker struct {
	pasteo       *pasteo.V2
	symmetricKey []byte
}

// NewPasteoMaker creates a new Pasteo token.
func NewPasteoMaker(symmetricKey string) (Maker, error) {
	if len(symmetricKey) != chacha20poly1305.KeySize {
		return nil, fmt.Errorf("invalid key size: must be exactly %d characters", chacha20poly1305.KeySize)
	}
	return &PasteoMaker{
		pasteo:       pasteo.NewV2(),
		symmetricKey: []byte(symmetricKey),
	}, nil
}

// CreateToken creates a new token for a specific username and duration
func (pt *PasteoMaker) CreateToken(username string, duration time.Duration) (string, error) {

	// 1. create the new payload
	payload, err := NewPayload(username, duration)
	if err != nil {
		return "", err
	}

	return pt.pasteo.Encrypt(pt.symmetricKey, payload, nil)
}

// VerifyToken checks if the token is valid or not
func (pt *PasteoMaker) VerifyToken(token string) (*Payload, error) {

	// 1. define payload var to store payload after decrypt
	payload := &Payload{}

	// 2. call the decrypt func to get the payload
	err := pt.pasteo.Decrypt(token, pt.symmetricKey, payload, nil)
	if err != nil {
		return nil, ErrInvalidToken
	}

	// 3. check if the payload is valid or not
	err = payload.Valid()
	if err != nil {
		return nil, ErrExpiredToken
	}

	return payload, nil
}
