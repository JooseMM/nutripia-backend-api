package authentication

import (
	"context"
	"time"

	"github.com/JooseMM/nutripia-backend-api/internal/nutritionist"
	"github.com/JooseMM/nutripia-backend-api/internal/nutritionist/types"
	"github.com/JooseMM/nutripia-backend-api/internal/nutritionist/types/dtos"
	authenticationDtos "github.com/JooseMM/nutripia-backend-api/internal/security/authentication/types/dtos"
	"github.com/JooseMM/nutripia-backend-api/pkg/core"
	"github.com/google/uuid"
)

type IAuthenticationService interface {
	RegisterNutritionist(
		dto *nutritionistDtos.CreateNutritionistRequest,
		ctx context.Context,
	) (*uuid.UUID, *core.BaseError)
	LoginNutritionist(
		dto *authenticationDtos.LoginRequest,
		ctx context.Context,
	) (*string, *core.BaseError)
}

type AuthenticationService struct {
	NutritionistRepo nutritionist.INutritionistRepository
}

func NewAuthenticationService(
	nutritionistRepo nutritionist.INutritionistRepository,
) IAuthenticationService {
	return &AuthenticationService{nutritionistRepo}
}

func (s *AuthenticationService) RegisterNutritionist(
	userDto *nutritionistDtos.CreateNutritionistRequest,
	ctx context.Context,
) (*uuid.UUID, *core.BaseError) {
	foundEmailOwner, queryEmailErr := s.NutritionistRepo.GetByEmailAddress(
		ctx,
		userDto.EmailAddress,
	)
	if queryEmailErr != nil &&
		queryEmailErr.ErrorCode != string(nutritionist.NUTRITIONIST_NOT_FOUND) {
		return nil, queryEmailErr
	}
	if foundEmailOwner != nil {
		return nil, nutritionist.UserEmailAlreadyExisting(userDto.EmailAddress)
	}

	hash, hashErr := HashPassword(userDto.Password)
	if hashErr != nil {
		return nil, core.UnexpectedError(hashErr.Error())
	}

	now := time.Now()
	user := &nutritionistTypes.Nutritionist{
		ID: uuid.New(),
		NutritionistIdentity: nutritionistTypes.NutritionistIdentity{
			Firstname:    userDto.Firstname,
			Lastname:     userDto.Lastname,
			EmailAddress: userDto.EmailAddress,
			DateBirth:    userDto.BirthDate,
		},
		AuthenticationInformation: nutritionistTypes.AuthenticationInformation{
			PasswordHash:     hash,
			IsEmailConfirmed: false,
		},
		TrackingInformation: nutritionistTypes.TrackingInformation{
			CreatedAt: now,
			UpdatedAt: now,
		},
	}

	creationErr := s.NutritionistRepo.Create(ctx, user)
	if creationErr != nil {
		return nil, creationErr
	}

	return &user.ID, nil
}

func (s *AuthenticationService) LoginNutritionist(
	dto *authenticationDtos.LoginRequest,
	ctx context.Context,
) (*string, *core.BaseError) {

	var sessionToken = ""
	return &sessionToken, nil
}
