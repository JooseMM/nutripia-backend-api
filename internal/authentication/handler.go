package authentication

import (
	"fmt"
	"net/http"

	authenticationDtos "github.com/JooseMM/nutripia-backend-api/internal/authentication/dtos"
	nutritionistDtos "github.com/JooseMM/nutripia-backend-api/internal/nutritionist/dtos"
	"github.com/JooseMM/nutripia-backend-api/pkg/core/valueobject"
	"github.com/JooseMM/nutripia-backend-api/pkg/request"
	"github.com/JooseMM/nutripia-backend-api/pkg/response"
)

type Handler interface {
	RegisterNutritionist(w http.ResponseWriter, r *http.Request)
	LoginNutritionist(w http.ResponseWriter, r *http.Request)
	ConfirmNutritionistEmail(w http.ResponseWriter, r *http.Request)
	SendResetPasswordToken(w http.ResponseWriter, r *http.Request)
	VerifyResetPasswordToken(w http.ResponseWriter, r *http.Request)
	CompleteResetPassword(w http.ResponseWriter, r *http.Request)
}

type handler struct {
	service IAuthenticationService
}

func NewAuthenticationHandler(service IAuthenticationService) Handler {
	return &handler{service}
}

func (h *handler) RegisterNutritionist(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	dto, err := request.DecodeJSON[nutritionistDtos.RegisterNutritionistInput](w, r)
	if err != nil {
		response.WriteJSON(w, err.StatusCode, err)
		return
	}
	fmt.Printf("firstname: %v", dto.Fullname.Firstname)
	fmt.Printf("lastname: %v", dto.Fullname.Lastname)
	fmt.Printf("%+v", dto.BirthDate)
	fmt.Printf("%+v", dto.EmailAddress)
	fmt.Printf("%+v", dto.Password)
	fmt.Printf("%+v", dto.Rut)

	payload, err := dto.ToEntity()
	if err != nil {
		response.WriteJSON(w, err.StatusCode, err)
		return
	}

	fmt.Printf("%+v", dto)

	failure := h.service.RegisterNutritionist(r.Context(), *payload)
	if failure != nil {
		response.WriteJSON(w, failure.StatusCode, failure)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *handler) LoginNutritionist(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	dto, err := request.DecodeJSON[authenticationDtos.RawLoginRequest](w, r)
	if err != nil {
		response.WriteJSON(w, err.StatusCode, err)
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

func (h *handler) ConfirmNutritionistEmail(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	dto, err := request.DecodeJSON[authenticationDtos.Token](w, r)
	if err != nil {
		response.WriteJSON(w, err.StatusCode, err)
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

func (h *handler) SendResetPasswordToken(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	dto, err := request.DecodeJSON[authenticationDtos.SendResetPasswordTokenRequest](w, r)
	if err != nil {
		response.WriteJSON(w, err.StatusCode, err)
		return
	}

	email, err := valueobject.NewEmailAddress(dto.EmailAddress)
	if err != nil {
		response.WriteJSON(w, err.StatusCode, err)
		return
	}

	if err := h.service.SendResetPasswordToken(r.Context(), email); err != nil {
		fmt.Println("Err: " + err.Description)
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *handler) VerifyResetPasswordToken(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	dto, err := request.DecodeJSON[authenticationDtos.Token](w, r)
	if err != nil {
		response.WriteJSON(w, err.StatusCode, err)
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

func (h *handler) CompleteResetPassword(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	dto, err := request.DecodeJSON[authenticationDtos.ChangePasswordRequest](w, r)
	if err != nil {
		response.WriteJSON(w, err.StatusCode, err)
		return
	}

	id := valueobject.IdentifierFromValue(dto.UserId)

	password, err := valueobject.NewPassword(dto.Password)
	if err != nil {
		response.WriteJSON(w, err.StatusCode, err)
		return
	}

	if err := h.service.CompleteResetPassword(r.Context(), id, password); err != nil {
		response.WriteJSON(w, err.StatusCode, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
