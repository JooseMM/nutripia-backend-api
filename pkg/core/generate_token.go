package core

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"math/big"
)

func GenerateAZToken(size uint) (string, error) {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ123456789"
	
	token := make([]byte, size)
	for i := range token {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil { 
			return "", err
		}
		token[i] = charset[num.Int64()]
	}

	return string(token), nil
}

func GenerateToken(size int) (*string, error) {
	b := make([]byte, size)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}

	token := base64.RawURLEncoding.EncodeToString(b)
	return &token, nil
}


func HashToken(token *string) string {
	hash := sha256.Sum256([]byte(*token))
	return hex.EncodeToString(hash[:])
}
