package session

import (
	"context"

	"github.com/JooseMM/nutripia-backend-api/pkg/core"
	"github.com/JooseMM/nutripia-backend-api/pkg/core/valueobject"
)

type SessionService struct {
	Repo Repository
}

type SessionManger interface {
	CreateSession(
		ctx context.Context,
		userId valueobject.Identifier,
		role valueobject.UserRoles,
	) (valueobject.Tokenizer, *core.BaseError)
	VerifySession(
		ctx context.Context,
		token valueobject.Tokenizer,
	) (Sessioner, *core.BaseError)
	DeleteSession(ctx context.Context, sessionId valueobject.Identifier) *core.BaseError
	DeleteAllSessionByUserId(ctx context.Context, userId valueobject.Identifier) *core.BaseError
}

func NewSessionService(repo Repository) SessionManger {
	return &SessionService{repo}
}

func (s *SessionService) CreateSession(
	ctx context.Context,
	userId valueobject.Identifier,
	role valueobject.UserRoles,
) (valueobject.Tokenizer, *core.BaseError) {
	session, err := NewSession(userId, role)
	if err != nil {
		return nil, err
	}

	if err := s.Repo.CreateSession(ctx, session); err != nil {
		return nil, err
	}

	return session.Token(), nil
}

func (s *SessionService) VerifySession(
	ctx context.Context,
	token valueobject.Tokenizer,
) (Sessioner, *core.BaseError) {
	session, sessionErr := s.Repo.GetByHash(ctx, token)
	if sessionErr != nil {
		return nil, sessionErr
	}

	if session.IsExpired() {
		return nil, SessionExpired()
	}

	return session, nil
}

func (s *SessionService) DeleteSession(
	ctx context.Context,
	sessionId valueobject.Identifier,
) *core.BaseError {
	return s.Repo.DeleteOne(ctx, sessionId)
}

func (s *SessionService) DeleteAllSessionByUserId(
	ctx context.Context,
	userId valueobject.Identifier,
) *core.BaseError {
	return s.Repo.DeleteByUserId(ctx, userId)
}
