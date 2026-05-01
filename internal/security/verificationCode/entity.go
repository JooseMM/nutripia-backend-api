package verification

import (
	"github.com/JooseMM/nutripia-backend-api/pkg/core/valueobject"
)

type VerificationManager interface {
	ToDB() entityDB
	IsExpired() bool
}

type verificationToken struct {
	id        valueobject.Identifier
	token     valueobject.Tokenizer
	tokenType valueobject.TokenTypeEnum
	userId    valueobject.Identifier
	createdAt valueobject.Dater
	expiredAt valueobject.Dater
}

func (v *verificationToken) ToDB() entityDB {
	return entityDB{
		id:        v.id.Key(),
		token:     v.token.String(),
		userId:    v.userId.Key(),
		createdAt: v.createdAt.ToTime(),
		expiredAt: v.expiredAt.ToTime(),
	}
}

func (v *verificationToken) IsExpired() bool {
	return v.expiredAt.IsPastOrNow()
}

func NewVerificationToken(
	token valueobject.Tokenizer,
	userId valueobject.Identifier,
) VerificationManager {
	now := valueobject.NewTracker()
	return &verificationToken{
		id:        valueobject.NewIdentifier(),
		token:     token,
		userId:    userId,
		createdAt: now,
		expiredAt: valueobject.NewExpiredDate(30),
	}
}
