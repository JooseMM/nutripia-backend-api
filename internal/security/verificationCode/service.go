package verification

import (
	"context"

	"github.com/JooseMM/nutripia-backend-api/pkg/core"
	"github.com/JooseMM/nutripia-backend-api/pkg/core/valueobject"
)

type VerificationCodeService struct {
	Repo Repository
}

type IVerificationCodeService interface {
	Create(ctx context.Context, verificationToken VerificationManager) *core.BaseError
	Verify(
		ctx context.Context,
		token valueobject.Tokenizer,
	) (*valueobject.Identifier, *core.BaseError)
}

func NewVerificationCodeService(repo Repository) IVerificationCodeService {
	return &VerificationCodeService{repo}
}

func (s *VerificationCodeService) Create(
	ctx context.Context,
	verificationToken VerificationManager,
) *core.BaseError {
	if err := s.Repo.Create(ctx, verificationToken); err != nil {
		return err
	}

	return nil
}

func (s *VerificationCodeService) Verify(
	ctx context.Context,
	token valueobject.Tokenizer,
) (*valueobject.Identifier, *core.BaseError) {
	verificationCode, verificationErr := s.Repo.GetByToken(ctx, token)
	if verificationErr != nil {
		return nil, verificationErr
	}

	if verificationCode.IsExpired() {
		return nil, VerificationCodeExpired()
	}

	if err := s.Repo.Delete(ctx, verificationCode.id); err != nil {
		return nil, err
	}

	return &verificationCode.userId, nil
}
