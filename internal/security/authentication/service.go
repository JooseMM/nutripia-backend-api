package authentication

import (
	"context"

	"github.com/JooseMM/nutripia-backend-api/internal/notifications/email"
	"github.com/JooseMM/nutripia-backend-api/internal/notifications/templates"
	"github.com/JooseMM/nutripia-backend-api/internal/nutritionist"
	nutritionistDtos "github.com/JooseMM/nutripia-backend-api/internal/nutritionist/dtos"
	authenticationDtos "github.com/JooseMM/nutripia-backend-api/internal/security/authentication/dtos"
	"github.com/JooseMM/nutripia-backend-api/internal/security/session"
	"github.com/JooseMM/nutripia-backend-api/internal/security/verificationCode"
	"github.com/JooseMM/nutripia-backend-api/pkg/core"
	"github.com/JooseMM/nutripia-backend-api/pkg/core/valueobject"
)

const FRONT_URL = "https://www.example.com"

type IAuthenticationService interface {
	RegisterNutritionist(
		ctx context.Context,
		dto nutritionistDtos.RegisterNutritionist,
	) *core.BaseError
	Login(
		ctx context.Context,
		dto authenticationDtos.LoginRequest,
		role valueobject.UserRoles,
	) (*authenticationDtos.LoginResponse, *core.BaseError)
	ConfirmEmail(
		ctx context.Context,
		token valueobject.Tokenizer,
	) (valueobject.Tokenizer, *core.BaseError)
	SendResetPasswordToken(
		ctx context.Context,
		emailAddress valueobject.Emailer,
	) *core.BaseError
	VerifyResetPasswordToken(
		ctx context.Context,
		token valueobject.Tokenizer,
	) (valueobject.Identifier, *core.BaseError)
	CompleteResetPassword(
		ctx context.Context,
		userId valueobject.Identifier,
		password valueobject.Passworder,
	) *core.BaseError
}

type AuthenticationService struct {
	NutritionistRepo    nutritionist.RepositoryManager
	SessionService      session.SessionManger
	VerificationService verification.VerificationCodeService
}

func NewAuthenticationService(
	nutritionistRepo nutritionist.RepositoryManager,
	sessionService session.SessionManger,
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
	isFound, err := s.NutritionistRepo.IsEmailTaken(ctx, data.Email())
	if err != nil {
		return err
	}
	if isFound {
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
	if err := s.VerificationService.Create(ctx, verificationToken); err != nil {
		return err
	}
	return s.sendEmailConfirmationToken(nutritionist.EmailAddress(), nutritionist.Name(), token)
}

func (s *AuthenticationService) Login(
	ctx context.Context,
	dto authenticationDtos.LoginRequest,
	role valueobject.UserRoles,
) (*authenticationDtos.LoginResponse, *core.BaseError) {
	user, userErr := s.NutritionistRepo.GetByEmailAddress(ctx, dto.EmailAddress())
	if userErr != nil {
		return nil, WrongCredentials()
	}

	if !user.IsEmailConfirmed() {
		return nil, EmailNotConfirmed()
	}

	if user.Password().IsEqual(dto.Password()) {
		return nil, WrongCredentials()
	}

	token, tokenErr := s.SessionService.CreateSession(ctx, user.Id(), role)
	if tokenErr != nil {
		return nil, tokenErr
	}

	response := authenticationDtos.NewLoginResponse(
		user.Id(),
		user.Name(),
		valueobject.NUTRITIONIST,
		token,
	)
	return &response, nil
}

func (s *AuthenticationService) ConfirmEmail(
	ctx context.Context,
	token valueobject.Tokenizer,
) (valueobject.Tokenizer, *core.BaseError) {
	userId, verificationErr := s.VerificationService.Verify(ctx, token)
	if verificationErr != nil {
		return nil, verificationErr
	}

	user, userErr := s.NutritionistRepo.GetById(ctx, *userId)
	if userErr != nil {
		return nil, userErr
	}

	user.MarkAsEmailAddressConfirmed()
	if err := s.NutritionistRepo.Update(ctx, user); err != nil {
		return nil, err
	}

	token, tokenErr := s.SessionService.CreateSession(ctx, *userId, valueobject.NUTRITIONIST)
	if tokenErr != nil {
		return nil, tokenErr
	}

	return token, nil
}

func (s *AuthenticationService) VerifyResetPasswordToken(
	ctx context.Context,
	token valueobject.Tokenizer,
) (valueobject.Identifier, *core.BaseError) {
	userId, verificationErr := s.VerificationService.Verify(ctx, token)
	if verificationErr != nil {
		return nil, verificationErr
	}

	return *userId, nil
}

func (s *AuthenticationService) SendResetPasswordToken(
	ctx context.Context,
	emailAddress valueobject.Emailer,
) *core.BaseError {
	user, userErr := s.NutritionistRepo.GetByEmailAddress(ctx, emailAddress)
	if userErr != nil {
		return userErr
	}

	if !user.IsEmailConfirmed() {
		return EmailNotConfirmed()
	}

	token, err := valueobject.NewPasswordResetToken()
	if err != nil {
		return err
	}

	verificationToken := verification.NewVerificationToken(token, user.Id())
	if err := s.VerificationService.Create(
		ctx,
		verificationToken,
	); err != nil {
		return err
	}

	return s.sendResetPasswordToken(user.EmailAddress(), token)
}

func (s *AuthenticationService) CompleteResetPassword(
	ctx context.Context,
	userId valueobject.Identifier,
	password valueobject.Passworder,
) *core.BaseError {
	user, userErr := s.NutritionistRepo.GetById(ctx, userId)
	if userErr != nil {
		return userErr
	}

	if !user.IsEmailConfirmed() {
		return EmailNotConfirmed()
	}

	user.UpdatePassword(password)
	if err := s.NutritionistRepo.Update(ctx, user); err != nil {
		return err
	}

	if err := s.SessionService.DeleteAllSessionByUserId(ctx, user.Id()); err != nil {
		return err
	}

	return nil
}

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

func (s *AuthenticationService) sendResetPasswordToken(
	targetEmail valueobject.Emailer,
	token valueobject.Tokenizer,
) *core.BaseError {

	replacements := map[string]string{
		"{{url}}": FRONT_URL + "/authentication/change-password" + token.String(),
	}

	sender, senderErr := email.NewMailSender(
		[]string{targetEmail.String()},
		"Cambio de Contraseña",
		templates.RESET_PASSWORD_TOKEN,
		replacements,
	)
	if senderErr != nil {
		return senderErr
	}

	return sender.Send()
}
