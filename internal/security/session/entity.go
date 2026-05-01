package session

import (
	"github.com/JooseMM/nutripia-backend-api/pkg/core"
	"github.com/JooseMM/nutripia-backend-api/pkg/core/valueobject"
)

type Sessioner interface {
	ToDB() sessions
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

func (s *session) ToDB() sessions {
	return sessions{
		Id:        s.id.Key(),
		Token:     s.token.Hash(),
		UserId:    s.userId.Key(),
		Role:      s.role,
		ExpiredAt: s.expiredAt.ToTime(),
		CreatedAt: s.createdAt.ToTime(),
	}
}

func (s *session) Token() valueobject.Tokenizer {
	return s.token
}

func (s *session) IsExpired() bool {
	return s.expiredAt.IsPastOrNow()
}
