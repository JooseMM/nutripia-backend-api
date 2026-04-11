package authentication

import (
	"encoding/json"
	"net/http"

	nutritionistDtos "github.com/JooseMM/nutripia-backend-api/internal/nutritionist/types/dtos"
	authenticationTypes "github.com/JooseMM/nutripia-backend-api/internal/security/authentication/types"
	authenticationDtos "github.com/JooseMM/nutripia-backend-api/internal/security/authentication/types/dtos"
	"github.com/JooseMM/nutripia-backend-api/pkg/core"
	"github.com/JooseMM/nutripia-backend-api/pkg/response"
)

type IAuthenticationHandler interface {
	RegisterNutritionist(w http.ResponseWriter, r *http.Request)
	LoginNutritionist(w http.ResponseWriter, r *http.Request)
}

type AuthenticationHandler struct {
	service IAuthenticationService
}

func NewAuthenticationHandler(service IAuthenticationService) IAuthenticationHandler {
	return &AuthenticationHandler{service}
}

func (h *AuthenticationHandler) RegisterNutritionist(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var dto nutritionistDtos.CreateNutritionistRequest

	err := json.NewDecoder(r.Body).Decode(&dto)
	if err != nil {
		http.Error(w, "Bad request: "+err.Error(), http.StatusBadRequest)
		return
	}

	if err := dto.Validate(); err != nil {
		response.WriteJSON(w, err.StatusCode, err)
		return
	}

	id, failure := h.service.RegisterNutritionist(&dto, r.Context())
	if failure != nil {
		response.WriteJSON(w, failure.StatusCode, failure)
		return
	}

	idStr := id.String()
	apiResponse := &response.ApiResponse[string]{
		Success: true,
		Data:    &idStr,
	}
	response.WriteJSON(w, http.StatusCreated, apiResponse)
}

func (h *AuthenticationHandler) LoginNutritionist(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var dto authenticationDtos.LoginRequest

	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		responseErr := core.ValidationError([]string{err.Error()})
		response.WriteJSON(w, responseErr.StatusCode, responseErr)
		return
	}

	role := authenticationTypes.NUTRITIONIST
	token, err := h.service.Login(&dto, &role, r.Context())
	if err != nil {
		response.WriteJSON(w, err.StatusCode, err)
		return
	}

	apiResponse := &response.ApiResponse[string]{
		Success: true,
		Data:    token,
	}

	response.WriteJSON(w, http.StatusOK, apiResponse)
}
