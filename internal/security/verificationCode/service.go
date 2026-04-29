package verification

import (
	"context"
	"time"

	"github.com/JooseMM/nutripia-backend-api/pkg/core"
	"github.com/JooseMM/nutripia-backend-api/pkg/core/valueobject"
	"github.com/google/uuid"
)

type VerificationCodeService struct {
	Repo Repository
}

type IVerificationCodeService interface {
	Create(ctx context.Context, verificationToken verificationToken) *core.BaseError
	Verify(ctx context.Context, token valueobject.Tokenizer) (*uuid.UUID, *core.BaseError)
}

func NewVerificationCodeService(repo Repository) IVerificationCodeService {
	return &VerificationCodeService{repo}
}

func (s *VerificationCodeService) Create(
	ctx context.Context,
	verificationToken VerificationManager,
) *core.BaseError {
	hash := core.HashToken(token)
	verificationCode := &verificationTypes.VerificationToken{
		ID:        uuid.New(),
		Token:     hash,
		UserId:    verificationToken.userId,
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
