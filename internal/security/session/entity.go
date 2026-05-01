package session

import (
	"github.com/JooseMM/nutripia-backend-api/pkg/core"
	"github.com/JooseMM/nutripia-backend-api/pkg/core/valueobject"
)

type Sessioner interface {
	ToDB() entityDB
	IsExpired() bool
	Token() valueobject.Tokenizer
}

func NewSession(
	userId valueobject.Identifier,
	role valueobject.UserRoles,
) (Sessioner, *core.BaseError) {
	token, err := valueobject.NewSessionToken()
	if err != nil {
		return nil, err
	}
	return &session{
		id:        valueobject.NewIdentifier(),
		token:     token,
		role:      role,
		userId:    userId,
		expiredAt: valueobject.NewExpiredDate(60 * 24),
		createdAt: valueobject.NewTracker(),
	}, nil
}

type session struct {
	id        valueobject.Identifier
	token     valueobject.Tokenizer
	userId    valueobject.Identifier
	role      valueobject.UserRoles
	expiredAt valueobject.Dater
	createdAt valueobject.Dater
}

func (s *session) ToDB() entityDB {
	return entityDB{
		id:        s.id.Key(),
		token:     s.token.Hash(),
		userId:    s.userId.Key(),
		role:      s.role,
		expiredAt: s.expiredAt.ToTime(),
		createdAt: s.createdAt.ToTime(),
	}
}

func (s *session) Token() valueobject.Tokenizer {
	return s.token
}

func (s *session) IsExpired() bool {
	return s.expiredAt.IsPastOrNow()
}
