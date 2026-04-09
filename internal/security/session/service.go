package session

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"time"

	sessionTypes "github.com/JooseMM/nutripia-backend-api/internal/security/session/types"
	"github.com/JooseMM/nutripia-backend-api/pkg/core"
	"github.com/google/uuid"
)

type SessionService struct {
	Repo ISessionRepository
}

type ISessionService interface {
	CreateSession(userId *uuid.UUID, ctx context.Context) (*string, *core.BaseError)
	VerifySession(sessionHash *string, ctx context.Context) (*sessionTypes.Session, *core.BaseError)
	DeleteSession(sessionId *uuid.UUID, ctx context.Context) *core.BaseError
	DeleteAllSessionByUserId(userId *uuid.UUID, ctx context.Context) *core.BaseError
}

func NewSessionService(repo ISessionRepository) ISessionService {
	return &SessionService{repo}
}

func (s *SessionService) CreateSession(
	userId *uuid.UUID,
	ctx context.Context,
) (*string, *core.BaseError) {
	token, tokenErr := s.generateToken()
	if tokenErr != nil {
		return nil, core.UnexpectedError(tokenErr.Error())
	}

	now := time.Now().UTC()
	session := &sessionTypes.Session{
		ID:          uuid.New(),
		SessionHash: s.hashToken(token),
		UserId:      *userId,
		ExpiredAt:   now.Add(72 * time.Hour),
	}

	if err := s.Repo.CreateSession(ctx, *session); err != nil {
		return nil, err
	}

	return token, nil
}

func (s *SessionService) VerifySession(
	sessionToken *string,
	ctx context.Context,
) (*sessionTypes.Session, *core.BaseError) {
	hash := s.hashToken(sessionToken)

	session, sessionErr := s.Repo.GetByHash(ctx, &hash)
	if sessionErr != nil {
		return nil, sessionErr
	}

	if session.ExpiredAt.Before(time.Now().UTC()) {
		return nil, SessionExpired()
	}

	return session, nil
}

func (s *SessionService) DeleteSession(sessionId *uuid.UUID, ctx context.Context) *core.BaseError {
	return s.Repo.DeleteOne(ctx, sessionId)
}

func (s *SessionService) DeleteAllSessionByUserId(
	userId *uuid.UUID,
	ctx context.Context,
) *core.BaseError {
	return s.Repo.DeleteByUserId(ctx, userId)
}

func (s *SessionService) generateToken() (*string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}

	token := base64.RawURLEncoding.EncodeToString(b)
	return &token, nil
}

func (s *SessionService) hashToken(token *string) string {
	hash := sha256.Sum256([]byte(*token))
	return hex.EncodeToString(hash[:])
}
