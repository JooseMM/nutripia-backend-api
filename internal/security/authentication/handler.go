package authentication

import (
	"fmt"
	"net/http"

	nutritionistDtos "github.com/JooseMM/nutripia-backend-api/internal/nutritionist/dtos"
	authenticationDtos "github.com/JooseMM/nutripia-backend-api/internal/security/authentication/dtos"
	"github.com/JooseMM/nutripia-backend-api/pkg/core"
	"github.com/JooseMM/nutripia-backend-api/pkg/core/valueobject"
	"github.com/JooseMM/nutripia-backend-api/pkg/response"
)

type IAuthenticationHandler interface {
	RegisterNutritionist(w http.ResponseWriter, r *http.Request)
	LoginNutritionist(w http.ResponseWriter, r *http.Request)
	ConfirmNutritionistEmail(w http.ResponseWriter, r *http.Request)
	SendResetPasswordToken(w http.ResponseWriter, r *http.Request)
	VerifyResetPasswordToken(w http.ResponseWriter, r *http.Request)
	CompleteResetPassword(w http.ResponseWriter, r *http.Request)
}

type AuthenticationHandler struct {
	service IAuthenticationService
}

func NewAuthenticationHandler(service IAuthenticationService) IAuthenticationHandler {
	return &AuthenticationHandler{service}
}

func (h *AuthenticationHandler) RegisterNutritionist(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	dto, errList := core.DecodeJSON[nutritionistDtos.RawRegisterNutritionist](w, r)
	if errList != nil {
		err := core.ValidationError(errList)
		response.WriteJSON(w, err.StatusCode, err)
		return
	}

	payload, err := dto.ToValueObject()
	if err != nil {
		err := core.ValidationError(errList)
		response.WriteJSON(w, err.StatusCode, err)
		return
	}

	failure := h.service.RegisterNutritionist(r.Context(), *payload)
	if failure != nil {
		response.WriteJSON(w, failure.StatusCode, failure)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *AuthenticationHandler) LoginNutritionist(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	dto, errList := core.DecodeJSON[authenticationDtos.RawLoginRequest](w, r)
	if errList != nil {
		e := core.ValidationError(errList)
		response.WriteJSON(w, e.StatusCode, e)
		return
	}

	valueObject, err := dto.ToValueObject()
	if err != nil {
		response.WriteJSON(w, err.StatusCode, err)
		return
	}

	loginResponse, err := h.service.Login(
		r.Context(),
		*valueObject,
		valueobject.NUTRITIONIST,
	)
	if err != nil {
		response.WriteJSON(w, err.StatusCode, err)
		return
	}

	apiResponse := &response.ApiResponse[authenticationDtos.LoginResponse]{
		Success: true,
		Data:    *loginResponse,
	}
	response.WriteJSON(w, http.StatusOK, apiResponse)
}

func (h *AuthenticationHandler) ConfirmNutritionistEmail(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	dto, parseErr := core.DecodeJSON[authenticationDtos.Token](w, r)
	if parseErr != nil {
		e := core.ValidationError(parseErr)
		response.WriteJSON(w, e.StatusCode, e)
		return
	}

	token, err := valueobject.EmailConfirmationTokenFromString(dto.Token)
	if err != nil {
		response.WriteJSON(w, err.StatusCode, err)
		return
	}

	sessionToken, err := h.service.ConfirmEmail(r.Context(), token)
	if err != nil {
		response.WriteJSON(w, err.StatusCode, err)
		return
	}

	apiResponse := &response.ApiResponse[string]{
		Success: true,
		Data:    sessionToken.String(),
	}

	response.WriteJSON(w, http.StatusOK, apiResponse)
}

func (h *AuthenticationHandler) SendResetPasswordToken(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	dto, parseErr := core.DecodeJSON[authenticationDtos.SendResetPasswordTokenRequest](w, r)
	if parseErr != nil {
		e := core.ValidationError(parseErr)
		response.WriteJSON(w, e.StatusCode, e)
		return
	}

	email, parseErr := valueobject.NewEmailAddress(dto.EmailAddress)
	if parseErr != nil {
		e := core.ValidationError(parseErr)
		response.WriteJSON(w, e.StatusCode, e)
		return
	}

	if err := h.service.SendResetPasswordToken(r.Context(), email); err != nil {
		fmt.Println("Err: " + err.Description)
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *AuthenticationHandler) VerifyResetPasswordToken(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	dto, parseErr := core.DecodeJSON[authenticationDtos.Token](w, r)
	if parseErr != nil {
		e := core.ValidationError(parseErr)
		response.WriteJSON(w, e.StatusCode, e)
		return
	}

	token, err := valueobject.PasswordResetTokenFromString(dto.Token)
	if err != nil {
		response.WriteJSON(w, err.StatusCode, err)
		return
	}

	userId, verificationErr := h.service.VerifyResetPasswordToken(r.Context(), token)
	if verificationErr != nil {
		response.WriteJSON(w, verificationErr.StatusCode, verificationErr)
		return
	}

	apiResponse := &response.ApiResponse[string]{
		Success: true,
		Data:    userId.String(),
	}
	response.WriteJSON(w, http.StatusOK, apiResponse)
}

func (h *AuthenticationHandler) CompleteResetPassword(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	dto, parseErr := core.DecodeJSON[authenticationDtos.ChangePasswordRequest](w, r)
	if parseErr != nil {
		e := core.ValidationError(parseErr)
		response.WriteJSON(w, e.StatusCode, e)
		return
	}

	id := valueobject.IdentifierFromValue(dto.UserId)

	password, err := valueobject.NewPassword(dto.Password)
	if err != nil {
		e := core.ValidationError(err)
		response.WriteJSON(w, e.StatusCode, e)
		return
	}

	if err := h.service.CompleteResetPassword(r.Context(), id, password); err != nil {
		response.WriteJSON(w, err.StatusCode, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
