package session

import (
	"context"
	"errors"
	"time"

	"github.com/JooseMM/nutripia-backend-api/pkg/core"
	"github.com/JooseMM/nutripia-backend-api/pkg/core/valueobject"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type sessions struct {
	Id        uuid.UUID             `gorm:"type:uuid;primaryKey"`
	Token     string                `gorm:"column:token;not null"`
	UserId    uuid.UUID             `gorm:"type:uuid;index"`
	Role      valueobject.UserRoles `gorm:"type:int"`
	ExpiredAt time.Time             `gorm:"column:expired_at;index"`
	CreatedAt time.Time             `gorm:"column:created_at"`
}

func (sessions) TableName() string {
	return "sessions"
}

type repo struct {
	db *gorm.DB
}

type Repository interface {
	GetByHash(ctx context.Context, token valueobject.Tokenizer) (Sessioner, *core.BaseError)
	CreateSession(ctx context.Context, session Sessioner) *core.BaseError
	DeleteByUserId(ctx context.Context, userId valueobject.Identifier) *core.BaseError
	DeleteOne(ctx context.Context, sessionId valueobject.Identifier) *core.BaseError
}

func NewRepository(db *gorm.DB) (Repository, *core.BaseError) {
	if err := db.AutoMigrate(&sessions{}); err != nil {
		return nil, core.UnexpectedError(err.Error())
	}

	return &repo{db}, nil
}

func (s *repo) GetByHash(
	ctx context.Context,
	token valueobject.Tokenizer,
) (Sessioner, *core.BaseError) {
	var session session

	result := s.db.WithContext(ctx).First(&session, "token = ?", token.Hash())
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, SessionNotFound()
		}

		return nil, core.UnexpectedError(result.Error.Error())
	}

	return &session, nil
}

func (s *repo) CreateSession(
	ctx context.Context,
	session Sessioner,
) *core.BaseError {
	dto := session.ToDB()
	result := s.db.WithContext(ctx).Create(&dto)
	if result.Error != nil {
		return core.UnexpectedError(result.Error.Error())
	}
	return nil
}

func (s *repo) DeleteByUserId(
	ctx context.Context,
	userId valueobject.Identifier,
) *core.BaseError {
	result := s.db.WithContext(ctx).Where("user_id = ?", userId.Key()).Delete(&session{})
	if result.Error != nil {
		return core.UnexpectedError(result.Error.Error())
	}

	if result.RowsAffected == 0 {
		return SessionNotFound()
	}

	return nil
}

func (s *repo) DeleteOne(
	ctx context.Context,
	sessionId valueobject.Identifier,
) *core.BaseError {
	result := s.db.WithContext(ctx).Delete(&session{}, sessionId.Key())
	if result.Error != nil {
		return core.UnexpectedError(result.Error.Error())
	}

	if result.RowsAffected == 0 {
		return SessionNotFound()
	}

	return nil
}
