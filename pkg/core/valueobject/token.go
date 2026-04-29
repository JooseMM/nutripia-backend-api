package valueobject

import (
	"crypto/rand"
	"encoding/base64"
	"math/big"
	"unicode/utf8"

	"github.com/JooseMM/nutripia-backend-api/pkg/core"
)

const EmailVerificationTokenLength = 6
const PasswordResetTokenLength = 16

type Tokenizer interface {
	String() string
}

type Token struct {
	value string
}

func PasswordResetTokenFromString(raw string) (Tokenizer, *core.BaseError) {
	if utf8.RuneCountInString(raw) < PasswordResetTokenLength {
		return nil, core.UnexpectedError("Token: wrong format")
	}

	return &Token{
		value: raw,
	}, nil
}

func NewPasswordResetTokenFromString(raw string) (Tokenizer, *core.BaseError) {
	b := make([]byte, PasswordResetTokenLength)
	_, err := rand.Read(b)
	if err != nil {
		return nil, core.UnexpectedError(err.Error())
	}

	token := base64.RawURLEncoding.EncodeToString(b)
	resp := &Token{
		value: token,
	}
	return resp, nil
}

func NewEmailVerificationToken() (Tokenizer, *core.BaseError) {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ123456789"

	token := make([]byte, EmailVerificationTokenLength)
	for i := range token {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return nil, core.UnexpectedError(err.Error())
		}
		token[i] = charset[num.Int64()]
	}

	resp := &Token{
		value: string(token),
	}
	return resp, nil
}

func EmailVerificationTokenFromString(raw string) (Tokenizer, *core.BaseError) {
	if utf8.RuneCountInString(raw) < EmailVerificationTokenLength {
		return nil, core.UnexpectedError("Token: wrong format")
	}

	return &Token{
		value: raw,
	}, nil
}

func (t Token) String() string {
	return t.value
}
