package session

import (
	"context"

	sessionTypes "github.com/JooseMM/nutripia-backend-api/internal/security/session/types"
	"github.com/JooseMM/nutripia-backend-api/pkg/core"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type postgresRepository struct {
	db *gorm.DB
}

type ISessionRepository interface {
	GetByHash(ctx context.Context, hash *string) (*sessionTypes.Session, *core.BaseError)
	CreateSession(ctx context.Context, userId *uuid.UUID) (*string, *core.BaseError)
	DeleteByUserId(ctx context.Context, userId *uuid.UUID) (*string, *core.BaseError)
	DeleteOne(ctx context.Context, sessionId *uuid.UUID) (*string, *core.BaseError)
}

func NewSessionRepository(db *gorm.DB) ISessionRepository {
	return &postgresRepository{db}
}
