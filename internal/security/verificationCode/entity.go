package verification

import (
	"time"

	"github.com/JooseMM/nutripia-backend-api/pkg/core/valueobject"
)

type VerificationManager interface {
	ToDB() entityDB
	UserId() valueobject.Identifier
}

type verificationToken struct {
	id        valueobject.Identifier
	token     valueobject.Tokenizer 
	userId    valueobject.Identifier
	createdAt time.Time
	expiredAt time.Time
}

func (v *verificationToken) ToDB() entityDB {
	return entityDB{
		id:        v.id.Key(),
		token:     v.token.String(),
		userId:    v.userId.Key(),
		createdAt: v.createdAt,
		expiredAt: v.expiredAt,
	}
}

func NewVerificationToken(
	token valueobject.Tokenizer,
	userId valueobject.Identifier,
) VerificationManager {
	now := time.Now().UTC()
	return &verificationToken{
		id:        valueobject.NewIdentifier(),
		token:     token,
		userId:    userId,
		createdAt: now,
		expiredAt: now.Add(30 * time.Minute),
	}
}
