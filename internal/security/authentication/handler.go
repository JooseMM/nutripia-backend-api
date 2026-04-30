package authentication

import (
	"net/http"

	nutritionistDtos "github.com/JooseMM/nutripia-backend-api/internal/nutritionist/dtos"
	"github.com/JooseMM/nutripia-backend-api/pkg/core"
	"github.com/JooseMM/nutripia-backend-api/pkg/response"
)

type IAuthenticationHandler interface {
	RegisterNutritionist(w http.ResponseWriter, r *http.Request)
	// LoginNutritionist(w http.ResponseWriter, r *http.Request)
	// ConfirmedNutritionistEmail(w http.ResponseWriter, r *http.Request)
	// SendResetPasswordToken(w http.ResponseWriter, r *http.Request)
	// VerifyResetPasswordToken(w http.ResponseWriter, r *http.Request)
	// CompleteResetPassword(w http.ResponseWriter, r *http.Request)
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

// func (h *AuthenticationHandler) LoginNutritionist(w http.ResponseWriter, r *http.Request) {
// 	defer r.Body.Close()
// 	var dto authenticationDtos.LoginRequest
//
// 	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
// 		responseErr := core.ValidationError([]string{err.Error()})
// 		response.WriteJSON(w, responseErr.StatusCode, responseErr)
// 		return
// 	}
//
// 	if err := dto.Validate(); err != nil {
// 		response.WriteJSON(w, err.StatusCode, err)
// 		return
// 	}
//
// 	role := authenticationTypes.NUTRITIONIST
// 	ctx := r.Context()
// 	loginResponse, err := h.service.Login(&dto, &role, &ctx)
// 	if err != nil {
// 		response.WriteJSON(w, err.StatusCode, err)
// 		return
// 	}
//
// 	apiResponse := &response.ApiResponse[authenticationDtos.LoginResponse]{
// 		Success: true,
// 		Data:    loginResponse,
// 	}
//
// 	response.WriteJSON(w, http.StatusOK, apiResponse)
// }
//
// func (h *AuthenticationHandler) ConfirmedNutritionistEmail(w http.ResponseWriter, r *http.Request) {
// 	defer r.Body.Close()
// 	var dto authenticationDtos.Token
//
// 	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
// 		responseErr := core.ValidationError([]string{err.Error()})
// 		response.WriteJSON(w, responseErr.StatusCode, responseErr)
// 		return
// 	}
//
// 	if err := dto.Validate(); err != nil {
// 		response.WriteJSON(w, err.StatusCode, err)
// 		return
// 	}
//
// 	ctx := r.Context()
// 	sessionToken, err := h.service.ConfirmedEmail(&dto.Token, &ctx)
// 	if err != nil {
// 		response.WriteJSON(w, err.StatusCode, err)
// 		return
// 	}
//
// 	apiResponse := &response.ApiResponse[string]{
// 		Success: true,
// 		Data:    sessionToken,
// 	}
//
// 	response.WriteJSON(w, http.StatusOK, apiResponse)
// }
//
// func (h *AuthenticationHandler) SendResetPasswordToken(w http.ResponseWriter, r *http.Request) {
// 	defer r.Body.Close()
// 	var dto authenticationDtos.SendResetPasswordTokenRequest
//
// 	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
// 		responseErr := core.ValidationError([]string{err.Error()})
// 		response.WriteJSON(w, responseErr.StatusCode, responseErr)
// 		return
// 	}
//
// 	if err := dto.Validate(); err != nil {
// 		response.WriteJSON(w, err.StatusCode, err)
// 		return
// 	}
//
// 	ctx := r.Context()
// 	if err := h.service.SendResetPasswordToken(&dto.EmailAddress, &ctx); err != nil {
// 		fmt.Println("Err: " + err.Description)
// 	}
//
// 	w.WriteHeader(http.StatusNoContent)
// }
//
// func (h *AuthenticationHandler) VerifyResetPasswordToken(w http.ResponseWriter, r *http.Request) {
// 	defer r.Body.Close()
// 	var dto authenticationDtos.Token
//
// 	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
// 		responseErr := core.ValidationError([]string{err.Error()})
// 		response.WriteJSON(w, responseErr.StatusCode, responseErr)
// 		return
// 	}
//
// 	if err := dto.Validate(); err != nil {
// 		response.WriteJSON(w, err.StatusCode, err)
// 		return
// 	}
//
// 	ctx := r.Context()
// 	userId, verificationErr := h.service.VerifyResetPasswordToken(&dto.Token, &ctx)
// 	if verificationErr != nil {
// 		response.WriteJSON(w, verificationErr.StatusCode, verificationErr)
// 		return
// 	}
//
// 	userIdStr := userId.String()
// 	apiResponse := &response.ApiResponse[string]{
// 		Success: true,
// 		Data:    &userIdStr,
// 	}
// 	response.WriteJSON(w, http.StatusOK, apiResponse)
// }
//
// func (h *AuthenticationHandler) CompleteResetPassword(w http.ResponseWriter, r *http.Request) {
// 	defer r.Body.Close()
// 	var dto authenticationDtos.ChangePasswordRequest
//
// 	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
// 		responseErr := core.ValidationError([]string{err.Error()})
// 		response.WriteJSON(w, responseErr.StatusCode, responseErr)
// 		return
// 	}
//
// 	if err := dto.Validate(); err != nil {
// 		response.WriteJSON(w, err.StatusCode, err)
// 		return
// 	}
//
// 	ctx := r.Context()
// 	if err := h.service.CompleteResetPassword(&dto.UserId, &dto.Password, &ctx); err != nil {
// 		response.WriteJSON(w, err.StatusCode, err)
// 		return
// 	}
//
// 	w.WriteHeader(http.StatusNoContent)
// }
