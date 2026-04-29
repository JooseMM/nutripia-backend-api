package valueobject

import (
	"crypto/rand"
	"encoding/base64"
	"math/big"
	"unicode/utf8"

	"github.com/JooseMM/nutripia-backend-api/pkg/core"
)

const emailConfirmationTokenLength = 6
const passwordResetTokenLength = 16

type TokenTypeEnum int

const (
	EmailConfirmation TokenTypeEnum = iota
	PasswordReset
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
	if utf8.RuneCountInString(raw) < passwordResetTokenLength {
		return nil, core.UnexpectedError("Token: wrong format")
	}

	return &Token{
		value:     raw,
		tokenType: PasswordReset,
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
	if utf8.RuneCountInString(raw) < emailConfirmationTokenLength {
		return nil, core.UnexpectedError("Token: wrong format")
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
	return core.HashToken(&t.value)
}
