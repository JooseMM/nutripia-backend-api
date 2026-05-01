package verification

import (
	"github.com/JooseMM/nutripia-backend-api/pkg/core/valueobject"
)

type VerificationManager interface {
	ToDB() verificationCodes
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

func (v *verificationToken) ToDB() verificationCodes {
	return verificationCodes{
		Id:        v.id.Key(),
		Token:     v.token.String(),
		UserId:    v.userId.Key(),
		CreatedAt: v.createdAt.ToTime(),
		ExpiredAt: v.expiredAt.ToTime(),
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
