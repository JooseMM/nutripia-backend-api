package authentication

import (
	"context"
	"fmt"
	"time"

	"github.com/JooseMM/nutripia-backend-api/internal/notifications/email"
	"github.com/JooseMM/nutripia-backend-api/internal/notifications/templates"
	"github.com/JooseMM/nutripia-backend-api/internal/nutritionist"
	nutritionistTypes "github.com/JooseMM/nutripia-backend-api/internal/nutritionist/types"
	nutritionistDtos "github.com/JooseMM/nutripia-backend-api/internal/nutritionist/types/dtos"
	authenticationTypes "github.com/JooseMM/nutripia-backend-api/internal/security/authentication/types"
	authenticationDtos "github.com/JooseMM/nutripia-backend-api/internal/security/authentication/types/dtos"
	"github.com/JooseMM/nutripia-backend-api/internal/security/session"
	"github.com/JooseMM/nutripia-backend-api/internal/security/verificationCode"
	"github.com/JooseMM/nutripia-backend-api/pkg/core"
	"github.com/google/uuid"
)

type IAuthenticationService interface {
	RegisterNutritionist(
		dto *nutritionistDtos.CreateNutritionistRequest,
		ctx *context.Context,
	) *core.BaseError
	Login(
		dto *authenticationDtos.LoginRequest,
		role *authenticationTypes.UserRoles,
		ctx *context.Context,
	) (*string, *core.BaseError)
	ConfirmedEmail(
		token *string,
		ctx *context.Context,
	) (*string, *core.BaseError)
	SendResetPasswordToken(
		emailAddress *string,
		ctx *context.Context,
	) *core.BaseError
}

type AuthenticationService struct {
	NutritionistRepo    nutritionist.INutritionistRepository
	SessionService      session.ISessionService
	VerificationService verificationCode.IVerificationCodeService
}

func NewAuthenticationService(
	nutritionistRepo nutritionist.INutritionistRepository,
	sessionService session.ISessionService,
	verificationService verificationCode.IVerificationCodeService,
) IAuthenticationService {
	return &AuthenticationService{
		nutritionistRepo,
		sessionService,
		verificationService,
	}
}

func (s *AuthenticationService) RegisterNutritionist(
	userDto *nutritionistDtos.CreateNutritionistRequest,
	ctx *context.Context,
) *core.BaseError {
	foundEmailOwner, queryEmailErr := s.NutritionistRepo.GetByEmailAddress(
		ctx,
		userDto.EmailAddress,
	)
	if queryEmailErr != nil &&
		queryEmailErr.ErrorCode != string(nutritionist.NUTRITIONIST_NOT_FOUND) {
		return queryEmailErr
	}
	if foundEmailOwner != nil {
		return nutritionist.UserEmailAlreadyExisting(userDto.EmailAddress)
	}

	hash, hashErr := HashPassword(userDto.Password)
	if hashErr != nil {
		return core.UnexpectedError(hashErr.Error())
	}

	now := time.Now().UTC()
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
		return creationErr
	}

	token, tokenErr := core.GenerateAZToken(6)
	if tokenErr != nil {
		return core.UnexpectedError(tokenErr.Error())
	}

	if err := s.VerificationService.Create(
		user.ID,
		&token,
		now.Add(15*time.Minute),
		ctx,
	); err != nil {
		return err
	}

	return s.sendRegistrationNotification(&user.Firstname, &token)
}

func (s *AuthenticationService) Login(
	dto *authenticationDtos.LoginRequest,
	role *authenticationTypes.UserRoles,
	ctx *context.Context,
) (*string, *core.BaseError) {
	user, userErr := s.NutritionistRepo.GetByEmailAddress(ctx, dto.EmailAddress)
	if userErr != nil {
		return nil, WrongCredentials()
	}

	if !user.IsEmailConfirmed {
		return nil, EmailNotConfirmed()
	}

	if !VerifyPassword(dto.Password, user.PasswordHash) {
		return nil, WrongCredentials()
	}

	token, tokenErr := s.SessionService.CreateSession(&user.ID, role, ctx)
	if tokenErr != nil {
		return nil, tokenErr
	}

	return token, nil
}
func (s *AuthenticationService) ConfirmedEmail(
	token *string,
	ctx *context.Context,
) (*string, *core.BaseError) {
	userId, verificationErr := s.VerificationService.Verify(token, ctx)
	if verificationErr != nil {
		return nil, verificationErr
	}

	user, userErr := s.NutritionistRepo.GetById(ctx, userId)
	if userErr != nil {
		return nil, userErr
	}

	user.IsEmailConfirmed = true
	if err := s.NutritionistRepo.Update(ctx, user); err != nil {
		return nil, err
	}

	role := authenticationTypes.NUTRITIONIST
	token, tokenErr := s.SessionService.CreateSession(&user.ID, &role, ctx)
	if tokenErr != nil {
		return nil, tokenErr
	}

	return token, nil
}

func (s *AuthenticationService) sendRegistrationNotification(
	nutritionistName *string,
	token *string,
) *core.BaseError {
	frontURL := "https://www.google.com"

	replacements := map[string]string{
		"{{firstname}}": *nutritionistName,
		"{{code}}":      *token,
		"{{url}}":       frontURL,
	}

	sender, senderErr := email.NewMailSender(
		[]string{"josexmoreno1998@gmail.com"},
		"Finaliza tu registro",
		templates.REGISTRATION,
		replacements,
	)
	if senderErr != nil {
		return senderErr
	}

	return sender.Send()
}

func (s *AuthenticationService) SendResetPasswordToken(
	emailAddress *string,
	ctx *context.Context,
) *core.BaseError {
	user, userErr := s.NutritionistRepo.GetByEmailAddress(ctx, *emailAddress)
	if userErr != nil {
		return userErr
	}

	if !user.IsEmailConfirmed {
		return EmailNotConfirmed()
	}

	token, tokenErr := core.GenerateToken(16)
	if tokenErr != nil {
		return core.UnexpectedError(tokenErr.Error())
	}

	now := time.Now().UTC()
	if err := s.VerificationService.Create(
		user.ID,
		token,
		now.Add(15*time.Minute),
		ctx,
	); err != nil {
		return err
	}

	return s.sendResetPasswordNotification(token)
}

func (s *AuthenticationService) sendResetPasswordNotification(
	token *string,
) *core.BaseError {
	frontURL := fmt.Sprintf("https://www.google.com/%s", *token)

	replacements := map[string]string{
		"{{url}}": frontURL,
	}

	sender, senderErr := email.NewMailSender(
		[]string{"josexmoreno1998@gmail.com"},
		"Finaliza tu registro",
		templates.RESET_PASSWORD_TOKEN,
		replacements,
	)
	if senderErr != nil {
		return senderErr
	}

	return sender.Send()
}
