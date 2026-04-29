package session

import (
	"context"
	"errors"
	"time"

	authenticationTypes "github.com/JooseMM/nutripia-backend-api/internal/security/authentication/types"
	"github.com/JooseMM/nutripia-backend-api/pkg/core"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SessionDB struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey"`
	SessionHash string    `gorm:"column:session_hash;not null"`
	UserId      uuid.UUID `gorm:"type:uuid;index"`
	Role        authenticationTypes.UserRoles
	ExpiredAt   time.Time `gorm:"column:expired_at"`
}

type SessionRepository struct {
	db *gorm.DB
}

type ISessionRepository interface {
	GetByHash(ctx context.Context, hash *string) (*Session, *core.BaseError)
	CreateSession(ctx context.Context, session Session) *core.BaseError
	DeleteByUserId(ctx context.Context, userId *uuid.UUID) *core.BaseError
	DeleteOne(ctx context.Context, sessionId *uuid.UUID) *core.BaseError
}

func NewSessionRepository(db *gorm.DB) ISessionRepository {
	return &SessionRepository{db}
}

func (s *SessionRepository) GetByHash(
	ctx context.Context,
	sessionHash *string,
) (*Session, *core.BaseError) {
	var session Session

	result := s.db.WithContext(ctx).Find(&session, "session_hash = ?", sessionHash)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, SessionNotFound()
		}

		return nil, core.UnexpectedError(result.Error.Error())
	}

	return &session, nil
}

func (s *SessionRepository) CreateSession(
	ctx context.Context,
	session Session,
) *core.BaseError {
	result := s.db.WithContext(ctx).Create(session)
	if result.Error != nil {
		return core.UnexpectedError(result.Error.Error())
	}
	return nil
}

func (s *SessionRepository) DeleteByUserId(
	ctx context.Context,
	userId *uuid.UUID,
) *core.BaseError {
	result := s.db.WithContext(ctx).Where("user_id = ?", userId).Delete(&Session{})
	if result.Error != nil {
		return core.UnexpectedError(result.Error.Error())
	}

	if result.RowsAffected == 0 {
		return SessionNotFound()
	}

	return nil
}

func (s *SessionRepository) DeleteOne(
	ctx context.Context,
	sessionId *uuid.UUID,
) *core.BaseError {
	result := s.db.WithContext(ctx).Delete(&Session{}, sessionId)
	if result.Error != nil {
		return core.UnexpectedError(result.Error.Error())
	}

	if result.RowsAffected == 0 {
		return SessionNotFound()
	}

	return nil
}
