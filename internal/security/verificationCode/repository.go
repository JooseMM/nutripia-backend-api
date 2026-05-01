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

type verificationCodes struct {
	Id        uuid.UUID                 `gorm:"type:uuid;primaryKey"`
	Token     string                    `gorm:"type:varchar(255);not null;index"`
	UserId    uuid.UUID                 `gorm:"type:uuid;not null;index"`
	TokenType valueobject.TokenTypeEnum `gorm:"type:int;not null"`
	CreatedAt time.Time                 `gorm:"not null"`
	ExpiredAt time.Time                 `gorm:"not null;index"`
}

func (e *verificationCodes) ToEntity() (*verificationToken, *core.BaseError) {
	var token valueobject.Tokenizer

	switch e.TokenType {
	case valueobject.EmailConfirmation:
		t, err := valueobject.EmailConfirmationTokenFromString(e.Token)
		if err != nil {
			return nil, core.CorrupetedDatabase(err.Details)
		}
		token = t
	case valueobject.PasswordReset:
		t, err := valueobject.PasswordResetTokenFromString(e.Token)
		if err != nil {
			return nil, core.CorrupetedDatabase(err.Details)
		}
		token = t
	}

	created, err := valueobject.NewTrackerFromTime(&e.CreatedAt)
	if err != nil {
		e := core.CorrupetedDatabase(err.Details)
		return nil, e
	}

	expired, err := valueobject.ExpiredDateFromTime(e.ExpiredAt)
	if err != nil {
		e := core.CorrupetedDatabase(err.Details)
		return nil, e
	}

	resp := verificationToken{
		id:        valueobject.IdentifierFromValue(e.Id),
		token:     token,
		tokenType: e.TokenType,
		userId:    valueobject.IdentifierFromValue(e.UserId),
		createdAt: created,
		expiredAt: expired,
	}

	return &resp, nil
}

func (verificationCodes) TableName() string {
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
	if err := db.AutoMigrate(&verificationCodes{}); err != nil {
		return nil, core.UnexpectedError(err.Error())
	}
	return &entityRepository{db}, nil
}

func (s *entityRepository) GetByToken(
	ctx context.Context,
	token valueobject.Tokenizer,
) (*verificationToken, *core.BaseError) {
	var dto verificationCodes

	result := s.db.WithContext(ctx).First(&dto, "token = ?", token.String())
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
	result := s.db.WithContext(ctx).Create(&dto)
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
		Delete(&verificationCodes{})
	if result.Error != nil {
		return core.UnexpectedError(result.Error.Error())
	}

	if result.RowsAffected == 0 {
		return VerificationCodeNotFound()
	}

	return nil
}
