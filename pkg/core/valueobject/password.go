package valueobject

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"strings"
	"unicode"

	"golang.org/x/crypto/argon2"
)

type params struct {
	memory      uint32
	iterations  uint32
	parallelism uint8
	saltLength  uint32
	keyLength   uint32
}

type Passworder interface {
	String() string
	IsEqual(rawPassword string) bool
}

type Password struct {
	hash string
}

func (p *Password) String() string {
	return p.hash
}

func (p *Password) IsEqual(rawPassword string) bool {
	// 1. Split the string into its component parts
	vals := strings.Split(p.hash, "$")
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
	params := &params{}
	_, err = fmt.Sscanf(
		vals[3],
		"m=%d,t=%d,p=%d",
		&params.memory,
		&params.iterations,
		&params.parallelism,
	)
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
	params.keyLength = uint32(len(existingHash))

	// 6. Hash the "new" password using the extracted parameters and salt
	newHash := argon2.IDKey(
		[]byte(rawPassword),
		salt,
		params.iterations,
		params.memory,
		params.parallelism,
		params.keyLength,
	)

	// 7. Compare them in constant time
	if subtle.ConstantTimeCompare(newHash, existingHash) == 1 {
		return true
	}

	return false
}

func PasswordFromDB(hash string) (Passworder, []string) {
	if hash == "" {
		return nil, []string{"Corrupted password comming from database"}
	}

	return &Password{
		hash: hash,
	}, nil
}

func NewPassword(p string) (Passworder, []string) {
	if p == "" {
		return nil, []string{"Password: is required"}
	}

	var errList []string

	if len(p) < 8 {
		errList = append(errList, "Password: must have a minimum of 8 characters.")
	}
	if len(p) > 72 {
		errList = append(errList, "Password: must have a maximum of 72 characters.")
	}

	var hasUpper, hasLower, hasNumber, hasSpecial bool
	for _, char := range p {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsDigit(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	if !hasUpper {
		errList = append(errList, "Password: must have at least one uppercase letter")
	}
	if !hasLower {
		errList = append(errList, "Password: must have at least one lowercase letter")
	}
	if !hasNumber {
		errList = append(errList, "Password: must have at least one number")
	}
	if !hasSpecial {
		errList = append(errList, "Password: must have at least one special character")
	}

	hash, err := hashPassword(p)
	if err != nil {
		errList = append(errList, "Password: "+err.Error())
	}

	if len(errList) > 0 {
		return nil, errList
	}

	return &Password{hash}, nil
}

func hashPassword(password string) (string, error) {
	p := &params{
		memory:      64 * 1024, // 64MB
		iterations:  3,
		parallelism: 2,
		saltLength:  16,
		keyLength:   32,
	}

	// 1. Generate a random salt
	salt := make([]byte, p.saltLength)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}

	// 2. Generate the hash
	hash := argon2.IDKey([]byte(password), salt, p.iterations, p.memory, p.parallelism, p.keyLength)

	// 3. Encode to a format you can store in the DB
	// Standard format: $argon2id$v=19$m=65536,t=3,p=2$salt$hash
	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	encodedHash := fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		argon2.Version, p.memory, p.iterations, p.parallelism, b64Salt, b64Hash)

	return encodedHash, nil
}
