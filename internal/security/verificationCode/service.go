package verificationCode

import (
	"context"
	"time"

	verificationTypes "github.com/JooseMM/nutripia-backend-api/internal/security/verificationCode/types"
	"github.com/JooseMM/nutripia-backend-api/pkg/core"
	"github.com/google/uuid"
)

type VerificationCodeService struct {
	Repo IVerificationCodeRepository
}

type IVerificationCodeService interface {
	Create(userId uuid.UUID,
		token *string,
		expiredAt time.Time,
		ctx *context.Context) *core.BaseError
	Verify(token *string, ctx *context.Context) (*uuid.UUID, *core.BaseError)
}

func NewVerificationCodeService(repo IVerificationCodeRepository) IVerificationCodeService {
	return &VerificationCodeService{repo}
}

func (s *VerificationCodeService) Create(
	userId uuid.UUID,
	token *string,
	expiredAt time.Time,
	ctx *context.Context,
) *core.BaseError {
	hash := core.HashToken(token)
	verificationCode := &verificationTypes.VerificationToken{
		ID:        uuid.New(),
		Token:     hash,
		UserId:    userId,
		ExpiredAt: expiredAt,
	}

	if err := s.Repo.Create(ctx, verificationCode); err != nil {
		return err
	}

	return nil
}

func (s *VerificationCodeService) Verify(
	token *string,
	ctx *context.Context,
) (*uuid.UUID, *core.BaseError) {
	hash := core.HashToken(token)
	verificationCode, verificationErr := s.Repo.GetByToken(ctx, &hash)
	if verificationErr != nil {
		return nil, verificationErr
	}

	if verificationCode.ExpiredAt.Before(time.Now().UTC()) {
		return nil, VerificationCodeExpired()
	}

	if err := s.Repo.Delete(ctx, &verificationCode.ID); err != nil {
		return nil, err
	}

	return &verificationCode.UserId, nil
}
