package core

import (
	"crypto/rand"
	"encoding/base64"
)

func GenerateToken(size int) (*string, error) {
	b := make([]byte, size)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}

	token := base64.RawURLEncoding.EncodeToString(b)
	return &token, nil
}
