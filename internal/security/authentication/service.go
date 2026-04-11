package authentication

import (
	"context"
	"time"

	"github.com/JooseMM/nutripia-backend-api/internal/nutritionist"
	"github.com/JooseMM/nutripia-backend-api/internal/nutritionist/types"
	"github.com/JooseMM/nutripia-backend-api/internal/nutritionist/types/dtos"
	authenticationTypes "github.com/JooseMM/nutripia-backend-api/internal/security/authentication/types"
	authenticationDtos "github.com/JooseMM/nutripia-backend-api/internal/security/authentication/types/dtos"
	"github.com/JooseMM/nutripia-backend-api/internal/security/session"
	"github.com/JooseMM/nutripia-backend-api/pkg/core"
	"github.com/google/uuid"
)

type IAuthenticationService interface {
	RegisterNutritionist(
		dto *nutritionistDtos.CreateNutritionistRequest,
		ctx context.Context,
	) (*uuid.UUID, *core.BaseError)
	Login(
		dto *authenticationDtos.LoginRequest,
		role *authenticationTypes.UserRoles,
		ctx context.Context,
	) (*string, *core.BaseError)
}

type AuthenticationService struct {
	NutritionistRepo nutritionist.INutritionistRepository
	SessionService   session.ISessionService
}

func NewAuthenticationService(
	nutritionistRepo nutritionist.INutritionistRepository,
	sessionService session.ISessionService,
) IAuthenticationService {
	return &AuthenticationService{
		nutritionistRepo,
		sessionService,
	}
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

func (s *AuthenticationService) Login(
	dto *authenticationDtos.LoginRequest,
	role *authenticationTypes.UserRoles,
	ctx context.Context,
) (*string, *core.BaseError) {
	user, userErr := s.NutritionistRepo.GetByEmailAddress(ctx, dto.EmailAddress)
	if userErr != nil {
		return nil, WrongCredentials()
	}

	// if !user.IsEmailConfirmed {
	// 	return nil, EmailNotConfirmed()
	// }

	if !VerifyPassword(dto.Password, user.PasswordHash) {
		return nil, WrongCredentials()
	}

	token, tokenErr := s.SessionService.CreateSession(&user.ID, role, ctx)
	if tokenErr != nil {
		return nil, tokenErr
	}

	return token, nil
}
