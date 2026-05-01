package valueobject

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"math/big"

	"github.com/JooseMM/nutripia-backend-api/pkg/core"
)

const emailConfirmationTokenLength = 6
const passwordResetTokenLength = 16
const sessionTokenLength = 32

type TokenTypeEnum int

const (
	EmailConfirmation TokenTypeEnum = iota
	PasswordReset
	Session
)

type Tokenizer interface {
	String() string
	Hash() string
}

type Token struct {
	value     string
	tokenType TokenTypeEnum
}

func PasswordResetTokenFromString(raw string) (Tokenizer, *core.BaseError) {
	if len(raw) < passwordResetTokenLength {
		return nil, core.UnexpectedError("Token: Wrong password-reset token format")
	}

	return &Token{
		value:     raw,
		tokenType: PasswordReset,
	}, nil
}

func NewSessionToken() (Tokenizer, *core.BaseError) {
	b := make([]byte, sessionTokenLength)
	_, err := rand.Read(b)
	if err != nil {
		return nil, core.UnexpectedError(err.Error())
	}

	token := base64.RawURLEncoding.EncodeToString(b)
	return &Token{
		value:     token,
		tokenType: Session,
	}, nil
}

func SessionTokenFromString(raw string) (Tokenizer, *core.BaseError) {
	if len(raw) < sessionTokenLength {
		return nil, core.UnexpectedError("Token: wrong session format token")
	}

	return &Token{
		value:     raw,
		tokenType: Session,
	}, nil
}

func NewPasswordResetToken() (Tokenizer, *core.BaseError) {
	b := make([]byte, passwordResetTokenLength)
	_, err := rand.Read(b)
	if err != nil {
		return nil, core.UnexpectedError(err.Error())
	}

	token := base64.RawURLEncoding.EncodeToString(b)
	resp := &Token{
		value:     token,
		tokenType: PasswordReset,
	}
	return resp, nil
}

func NewEmailConfirmationToken() (Tokenizer, *core.BaseError) {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ123456789"

	token := make([]byte, emailConfirmationTokenLength)
	for i := range token {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return nil, core.UnexpectedError(err.Error())
		}
		token[i] = charset[num.Int64()]
	}

	resp := &Token{
		value:     string(token),
		tokenType: EmailConfirmation,
	}
	return resp, nil
}

func EmailConfirmationTokenFromString(raw string) (Tokenizer, *core.BaseError) {
	if len(raw) < emailConfirmationTokenLength {
		return nil, core.UnexpectedError("Token: wrong email-confirmation token format")
	}

	return &Token{
		value:     raw,
		tokenType: EmailConfirmation,
	}, nil
}

func (t *Token) String() string {
	return t.value
}

func (t *Token) Hash() string {
	hash := sha256.Sum256([]byte(t.value))
	return hex.EncodeToString(hash[:])
}
