package core

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"math/big"
)



func HashToken(token *string) string {
	hash := sha256.Sum256([]byte(*token))
	return hex.EncodeToString(hash[:])
}
