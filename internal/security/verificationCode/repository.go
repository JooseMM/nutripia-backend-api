package verification

import (
	"context"
	"errors"
	"time"

	"github.com/JooseMM/nutripia-backend-api/pkg/core"
	"github.com/JooseMM/nutripia-backend-api/pkg/core/valueobject"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// TODO: add metadata
type entityDB struct {
	id        uuid.UUID
	token     string
	userId    uuid.UUID
	tokenType valueobject.TokenTypeEnum
	createdAt time.Time
	expiredAt time.Time
}

func (e *entityDB) ToEntity() (*verificationToken, *core.BaseError) {
	var token valueobject.Tokenizer

	switch e.tokenType {
	case valueobject.EmailConfirmation:
		t, err := valueobject.PasswordResetTokenFromString(e.token)
		if err != nil {
			return nil, err
		}
		token = t
	case valueobject.PasswordReset:
		t, err := valueobject.EmailConfirmationTokenFromString(e.token)
		if err != nil {
			return nil, err
		}
		token = t
	}

	resp := verificationToken{
		id:        valueobject.IdentifierFromValue(e.id),
		token:     token,
		tokenType: e.tokenType,
		userId:    valueobject.IdentifierFromValue(e.id),
		createdAt: e.createdAt,
		expiredAt: e.expiredAt,
	}

	return &resp, nil
}

func (entityDB) TableName() string {
	return "verification_tokens"
}

type entityRepository struct {
	db *gorm.DB
}

type Repository interface {
	GetByToken(
		ctx context.Context,
		token valueobject.Tokenizer,
	) (*verificationToken, *core.BaseError)
	Create(
		ctx context.Context,
		verificationCode VerificationManager,
	) *core.BaseError
	Delete(ctx context.Context, userId valueobject.Identifier) *core.BaseError
}

func NewRepository(db *gorm.DB) (Repository, *core.BaseError) {
	if err := db.AutoMigrate(&entityDB{}); err != nil {
		return nil, core.UnexpectedError(err.Error())
	}
	return &entityRepository{db}, nil
}

func (s *entityRepository) GetByToken(
	ctx context.Context,
	token valueobject.Tokenizer,
) (*verificationToken, *core.BaseError) {
	var dto entityDB

	result := s.db.WithContext(ctx).Find(&dto, "token = ?", token.String())
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, VerificationCodeNotFound()
		}

		return nil, core.UnexpectedError(result.Error.Error())
	}

	entity, err := dto.ToEntity()
	if err != nil {
		return nil, err
	}

	return entity, nil
}

func (s *entityRepository) Create(
	ctx context.Context,
	verificationToken VerificationManager,
) *core.BaseError {
	dto := verificationToken.ToDB()
	result := s.db.WithContext(ctx).Create(dto)
	if result.Error != nil {
		return core.UnexpectedError(result.Error.Error())
	}
	return nil
}

func (s *entityRepository) Delete(
	ctx context.Context,
	id valueobject.Identifier,
) *core.BaseError {
	result := s.db.WithContext(ctx).
		Where("id = ?", id).
		Delete(&entityDB{})
	if result.Error != nil {
		return core.UnexpectedError(result.Error.Error())
	}

	if result.RowsAffected == 0 {
		return VerificationCodeNotFound()
	}

	return nil
}
