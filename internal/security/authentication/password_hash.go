package authentication

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
)

type Params struct {
	Memory      uint32
	Iterations  uint32
	Parallelism uint8
	SaltLength  uint32
	KeyLength   uint32
}

func HashPassword(password string) (string, error) {
	p := &Params{
		Memory:      64 * 1024, // 64MB
		Iterations:  3,
		Parallelism: 2,
		SaltLength:  16,
		KeyLength:   32,
	}

	// 1. Generate a random salt
	salt := make([]byte, p.SaltLength)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}

	// 2. Generate the hash
	hash := argon2.IDKey([]byte(password), salt, p.Iterations, p.Memory, p.Parallelism, p.KeyLength)

	// 3. Encode to a format you can store in the DB
	// Standard format: $argon2id$v=19$m=65536,t=3,p=2$salt$hash
	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	encodedHash := fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		argon2.Version, p.Memory, p.Iterations, p.Parallelism, b64Salt, b64Hash)

	return encodedHash, nil
}

func VerifyPassword(password, encodedHash string) bool {
	// 1. Split the string into its component parts
	vals := strings.Split(encodedHash, "$")
	if len(vals) != 6 {
		return false
	}

	// 2. Verify the Argon2 version
	var version int
	_, err := fmt.Sscanf(vals[2], "v=%d", &version)
	if err != nil {
		return false
	}
	if version != argon2.Version {
		return false
	}

	// 3. Extract the parameters (Memory, Iterations, Parallelism)
	p := &Params{}
	_, err = fmt.Sscanf(vals[3], "m=%d,t=%d,p=%d", &p.Memory, &p.Iterations, &p.Parallelism)
	if err != nil {
		return false
	}

	// 4. Decode the Salt
	salt, err := base64.RawStdEncoding.DecodeString(vals[4])
	if err != nil {
		return false
	}

	// 5. Decode the existing Hash
	existingHash, err := base64.RawStdEncoding.DecodeString(vals[5])
	if err != nil {
		return false
	}
	p.KeyLength = uint32(len(existingHash))

	// 6. Hash the "new" password using the extracted parameters and salt
	newHash := argon2.IDKey(
		[]byte(password),
		salt,
		p.Iterations,
		p.Memory,
		p.Parallelism,
		p.KeyLength,
	)

	// 7. Compare them in constant time
	if subtle.ConstantTimeCompare(newHash, existingHash) == 1 {
		return true
	}

	return false
}
