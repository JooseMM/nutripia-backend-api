package verificationCode

import (
	"context"
	"errors"

	verificationTypes "github.com/JooseMM/nutripia-backend-api/internal/security/verificationCode/types"
	"github.com/JooseMM/nutripia-backend-api/pkg/core"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type VerificationCodeRepository struct {
	db *gorm.DB
}

type IVerificationCodeRepository interface {
	GetByToken(
		ctx *context.Context,
		token *string,
	) (*verificationTypes.VerificationToken, *core.BaseError)
	Create(
		ctx *context.Context,
		verificationCode *verificationTypes.VerificationToken,
	) *core.BaseError
	Delete(ctx *context.Context, userId *uuid.UUID) *core.BaseError
}

func NewVerificationCodeRepository(db *gorm.DB) IVerificationCodeRepository {
	return &VerificationCodeRepository{db}
}

func (s *VerificationCodeRepository) GetByToken(
	ctx *context.Context,
	token *string,
) (*verificationTypes.VerificationToken, *core.BaseError) {
	var verificationToken verificationTypes.VerificationToken

	result := s.db.WithContext(*ctx).Find(&verificationToken, "token = ?", token)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, VerificationCodeNotFound()
		}

		return nil, core.UnexpectedError(result.Error.Error())
	}

	return &verificationToken, nil
}

func (s *VerificationCodeRepository) Create(
	ctx *context.Context,
	verificationToken *verificationTypes.VerificationToken,
) *core.BaseError {
	result := s.db.WithContext(*ctx).Create(verificationToken)
	if result.Error != nil {
		return core.UnexpectedError(result.Error.Error())
	}
	return nil
}

func (s *VerificationCodeRepository) Delete(
	ctx *context.Context,
	id *uuid.UUID,
) *core.BaseError {
	result := s.db.WithContext(*ctx).
		Where("id = ?", id).
		Delete(&verificationTypes.VerificationToken{})
	if result.Error != nil {
		return core.UnexpectedError(result.Error.Error())
	}

	if result.RowsAffected == 0 {
		return VerificationCodeNotFound()
	}

	return nil
}
