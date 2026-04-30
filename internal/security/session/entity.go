package session

import (
	"time"

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
	now := time.Now().UTC()
	return &session{
		id:        valueobject.NewIdentifier(),
		token:     token,
		role:      role,
		userId:    userId,
		expiredAt: now.Add(72 * time.Hour),
		createdAt: now,
	}, nil
}

type session struct {
	id        valueobject.Identifier
	token     valueobject.Tokenizer
	userId    valueobject.Identifier
	role      valueobject.UserRoles
	expiredAt time.Time
	createdAt time.Time
}

func (s *session) ToDB() entityDB {
	return entityDB{
		id:        s.id.Key(),
		token:     s.token.Hash(),
		userId:    s.userId.Key(),
		role:      s.role,
		expiredAt: s.expiredAt,
		createdAt: s.createdAt,
	}
}

func (s *session) Token() valueobject.Tokenizer {
	return s.token
}

func (s *session) IsExpired() bool {
	now := time.Now().UTC()
	return s.expiredAt.Equal(now) || s.expiredAt.Before(now)
}
