package authentication

import (
	"context"
	"time"

	"github.com/JooseMM/nutripia-backend-api/internal/notifications/email"
	"github.com/JooseMM/nutripia-backend-api/internal/notifications/templates"
	"github.com/JooseMM/nutripia-backend-api/internal/nutritionist"
	nutritionistDtos "github.com/JooseMM/nutripia-backend-api/internal/nutritionist/dtos"
	authenticationTypes "github.com/JooseMM/nutripia-backend-api/internal/security/authentication/types"
	authenticationDtos "github.com/JooseMM/nutripia-backend-api/internal/security/authentication/types/dtos"
	"github.com/JooseMM/nutripia-backend-api/internal/security/session"
	"github.com/JooseMM/nutripia-backend-api/internal/security/verificationCode"
	"github.com/JooseMM/nutripia-backend-api/pkg/core"
	"github.com/JooseMM/nutripia-backend-api/pkg/core/valueobject"
	"github.com/google/uuid"
)

// Make this an enviroment variable
const FRONT_URL = "https://www.example.com"

type IAuthenticationService interface {
	RegisterNutritionist(
		ctx context.Context,
		dto nutritionistDtos.RegisterNutritionist,
	) *core.BaseError
	// Login(
	// 	ctx context.Context,
	// 	dto authenticationDtos.LoginRequest,
	// 	role authenticationTypes.UserRoles,
	// ) (*authenticationDtos.LoginResponse, *core.BaseError)
	// ConfirmedEmail(
	// 	ctx context.Context,
	// 	token string,
	// ) (*string, *core.BaseError)
	// SendResetPasswordToken(
	// 	ctx context.Context,
	// 	emailAddress string,
	// ) *core.BaseError
	// VerifyResetPasswordToken(
	// 	token *string,
	// 	ctx *context.Context,
	// ) (*uuid.UUID, *core.BaseError)
	// CompleteResetPassword(
	// 	ctx context.Context,
	// 	userId uuid.UUID,
	// 	password string,
	// ) *core.BaseError
}

type AuthenticationService struct {
	NutritionistRepo    nutritionist.RepositoryManager
	SessionService      session.ISessionService
	VerificationService verification.VerificationCodeService
}

func NewAuthenticationService(
	nutritionistRepo nutritionist.RepositoryManager,
	sessionService session.ISessionService,
	verificationService verification.VerificationCodeService,
) IAuthenticationService {
	return &AuthenticationService{
		nutritionistRepo,
		sessionService,
		verificationService,
	}
}

func (s *AuthenticationService) RegisterNutritionist(
	ctx context.Context,
	data nutritionistDtos.RegisterNutritionist,
) *core.BaseError {
	foundEmailOwner, queryEmailErr := s.NutritionistRepo.GetByEmailAddress(
		ctx,
		data.Email(),
	)
	if queryEmailErr != nil &&
		queryEmailErr.ErrorCode != string(nutritionist.NUTRITIONIST_NOT_FOUND) {
		return queryEmailErr
	}
	if foundEmailOwner != nil {
		return nutritionist.UserEmailAlreadyExisting(data.Email().String())
	}

	nutritionist := nutritionist.FromRegisterDTO(data)
	if err := s.NutritionistRepo.Create(ctx, nutritionist); err != nil {
		return err
	}

	token, err := valueobject.NewEmailConfirmationToken()
	if err != nil {
		return err
	}

	verificationToken := verification.NewVerificationToken(token, nutritionist.Id())
	if err := s.VerificationService.Create(verificationToken); err != nil {
		return err
	}
	return s.sendEmailConfirmationToken(nutritionist.EmailAddress(), nutritionist.Name(), token)
}

// func (s *AuthenticationService) Login(
// 	dto *authenticationDtos.LoginRequest,
// 	role *authenticationTypes.UserRoles,
// 	ctx *context.Context,
// ) (*authenticationDtos.LoginResponse, *core.BaseError) {
// 	user, userErr := s.NutritionistRepo.GetByEmailAddress(ctx, dto.EmailAddress)
// 	if userErr != nil {
// 		return nil, WrongCredentials()
// 	}
//
// 	if !user.IsEmailConfirmed {
// 		return nil, EmailNotConfirmed()
// 	}
//
// 	if !VerifyPassword(dto.Password, user.PasswordHash) {
// 		return nil, WrongCredentials()
// 	}
//
// 	token, tokenErr := s.SessionService.CreateSession(&user.ID, role, ctx)
// 	if tokenErr != nil {
// 		return nil, tokenErr
// 	}
//
// 	response := &authenticationDtos.LoginResponse{
// 		UserId:    user.ID,
// 		Firstname: user.Firstname,
// 		Role:      authenticationTypes.NUTRITIONIST,
// 		Token:     *token,
// 	}
// 	return response, nil
// }
// func (s *AuthenticationService) ConfirmedEmail(
// 	token *string,
// 	ctx *context.Context,
// ) (*string, *core.BaseError) {
// 	userId, verificationErr := s.VerificationService.Verify(token, ctx)
// 	if verificationErr != nil {
// 		return nil, verificationErr
// 	}
//
// 	user, userErr := s.NutritionistRepo.GetById(ctx, userId)
// 	if userErr != nil {
// 		return nil, userErr
// 	}
//
// 	user.IsEmailConfirmed = true
// 	if err := s.NutritionistRepo.Update(ctx, user); err != nil {
// 		return nil, err
// 	}
//
// 	role := authenticationTypes.NUTRITIONIST
// 	token, tokenErr := s.SessionService.CreateSession(&user.ID, &role, ctx)
// 	if tokenErr != nil {
// 		return nil, tokenErr
// 	}
//
// 	return token, nil
// }
//
// func (s *AuthenticationService) VerifyResetPasswordToken(
// 	token *string,
// 	ctx *context.Context,
// ) (*uuid.UUID, *core.BaseError) {
// 	userId, verificationErr := s.VerificationService.Verify(token, ctx)
// 	if verificationErr != nil {
// 		return nil, verificationErr
// 	}
//
// 	return userId, nil
// }
//
// func (s *AuthenticationService) SendResetPasswordToken(
// 	emailAddress *string,
// 	ctx *context.Context,
// ) *core.BaseError {
// 	user, userErr := s.NutritionistRepo.GetByEmailAddress(ctx, *emailAddress)
// 	if userErr != nil {
// 		return userErr
// 	}
//
// 	if !user.IsEmailConfirmed {
// 		return EmailNotConfirmed()
// 	}
//
// 	token, tokenErr := core.GenerateToken(16)
// 	if tokenErr != nil {
// 		return core.UnexpectedError(tokenErr.Error())
// 	}
//
// 	now := time.Now().UTC()
// 	if err := s.VerificationService.Create(
// 		user.ID,
// 		token,
// 		now.Add(15*time.Minute),
// 		ctx,
// 	); err != nil {
// 		return err
// 	}
//
// 	return s.sendResetPasswordToken(&user.EmailAddress, token)
// }
//
// func (s *AuthenticationService) CompleteResetPassword(
// 	ctx context.Context,
// 	userId valueobject.Identifier,
// 	password valueobject.Passworder,
// ) *core.BaseError {
// 	user, userErr := s.NutritionistRepo.GetById(ctx, userId)
// 	if userErr != nil {
// 		return userErr
// 	}
//
// 	if !user.IsEmailConfirmed() {
// 		return EmailNotConfirmed()
// 	}
//
// 	user.PasswordHash = core.HashToken(password)
// 	user.UpdatedAt = time.Now().UTC()
// 	if err := s.NutritionistRepo.Update(ctx, user); err != nil {
// 		return err
// 	}
//
// 	if err := s.SessionService.DeleteAllSessionByUserId(userId, ctx); err != nil {
// 		return err
// 	}
//
// 	return nil
// }
//

func (s *AuthenticationService) sendEmailConfirmationToken(
	targetEmail valueobject.Emailer,
	name valueobject.Namer,
	token valueobject.Tokenizer,
) *core.BaseError {
	replacements := map[string]string{
		"{{firstname}}": name.Firstname(),
		"{{code}}":      token.String(),
		"{{url}}":       FRONT_URL + "/authentication/verify-email",
	}

	sender, senderErr := email.NewMailSender(
		[]string{targetEmail.String()},
		"Finaliza tu registro",
		templates.REGISTRATION,
		replacements,
	)
	if senderErr != nil {
		return senderErr
	}

	return sender.Send()
}

//
// func (s *AuthenticationService) sendResetPasswordToken(
// 	targetEmail valueobject.Emailer,
// 	token valueobject.TokenManager,
// ) *core.BaseError {
//
// 	replacements := map[string]string{
// 		"{{url}}": FRONT_URL + "/authentication/change-password" + token.String(),
// 	}
//
// 	sender, senderErr := email.NewMailSender(
// 		[]string{targetEmail.String()},
// 		"Cambio de Contraseña",
// 		templates.RESET_PASSWORD_TOKEN,
// 		replacements,
// 	)
// 	if senderErr != nil {
// 		return senderErr
// 	}
//
// 	return sender.Send()
// }
