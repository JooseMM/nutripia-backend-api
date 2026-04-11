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
	Create(userId uuid.UUID, ctx *context.Context) (*string, *core.BaseError)
	Verify(token *string, ctx *context.Context) (*uuid.UUID, *core.BaseError)
}

func NewVerificationCodeService(repo IVerificationCodeRepository) IVerificationCodeService {
	return &VerificationCodeService{repo}
}

func (s *VerificationCodeService) Create(
	userId uuid.UUID,
	ctx *context.Context,
) (*string, *core.BaseError) {
	token, tokenErr := core.GenerateToken(6)
	if tokenErr != nil {
		return nil, core.UnexpectedError(tokenErr.Error())
	}

	now := time.Now().UTC()
	verificationCode := &verificationTypes.VerificationToken{
		ID:        uuid.New(),
		Token:     *token,
		UserId:    userId,
		ExpiredAt: now.Add(24 * time.Hour),
	}

	if err := s.Repo.Create(ctx, verificationCode); err != nil {
		return nil, err
	}

	return token, nil
}

func (s *VerificationCodeService) Verify(
	token *string,
	ctx *context.Context,
) (*uuid.UUID, *core.BaseError) {
	verificationCode, verificationErr := s.Repo.GetByToken(ctx, token)
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
