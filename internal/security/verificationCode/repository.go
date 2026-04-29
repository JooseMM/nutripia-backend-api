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

type entityDB struct {
	id        uuid.UUID
	token     string
	userId    uuid.UUID
	createdAt time.Time
	expiredAt time.Time
}

func (e *entityDB) ToEntity() verificationToken {
	return verificationToken{
		id:        valueobject.IdentifierFromDB(e.id),
		token:     e.token,
		userId:    valueobject.IdentifierFromDB(e.id),
		createdAt: e.createdAt,
		expiredAt: e.expiredAt,
	}
}

type entityRepository struct {
	db *gorm.DB
}

type Repository interface {
	GetByToken(
		ctx context.Context,
		token string,
	) (*verificationToken, *core.BaseError)
	Create(
		ctx context.Context,
		verificationCode *verificationToken,
	) *core.BaseError
	Delete(ctx context.Context, userId valueobject.Identifier) *core.BaseError
}

func NewVerificationCodeRepository(db *gorm.DB) Repository {
	return &entityRepository{db}
}

func (s *entityRepository) GetByToken(
	ctx context.Context,
	token string,
) (*verificationToken, *core.BaseError) {
	var dto entityDB

	result := s.db.WithContext(ctx).Find(&dto, "token = ?", token)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, VerificationCodeNotFound()
		}

		return nil, core.UnexpectedError(result.Error.Error())
	}

	entity := dto.ToEntity()
	return &entity, nil
}

func (s *entityRepository) Create(
	ctx context.Context,
	verificationToken *verificationToken,
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
