package nutritionist

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/JooseMM/nutripia-backend-api/internal/nutritionist/types/dtos"
	"github.com/JooseMM/nutripia-backend-api/pkg/core"
	"github.com/JooseMM/nutripia-backend-api/pkg/response"
	"github.com/google/uuid"
)

type INutritionistHandler interface {
	CreateNutritionist(w http.ResponseWriter, r *http.Request)
	GetClientById(w http.ResponseWriter, r *http.Request)
	UpdateClientById(w http.ResponseWriter, r *http.Request)
	DeleteClientById(w http.ResponseWriter, r *http.Request)
}

type NutriotionistHandler struct {
	service INutritionistService
}

func NewClientHandler(service INutritionistService) INutritionistHandler {
	return &NutriotionistHandler{service}
}

func (h *NutriotionistHandler) CreateNutritionist(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var dto nutritionistDtos.CreateNutritionistRequest

	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		http.Error(w, "Bad request: "+err.Error(), http.StatusBadRequest)
		return
	}

	if err := dto.Validate(); err != nil {
		response.WriteJSON(w, err.StatusCode, err)
		return
	}

	id, failure := h.service.CreateNutritionist(&dto, r.Context())
	if failure != nil {
		response.WriteJSON(w, failure.StatusCode, failure)
		return
	}

	idStr := id.String()
	responseJson := &response.ApiResponse[string]{
		Success: true,
		Data:    &idStr,
	}

	response.WriteJSON(w, http.StatusCreated, responseJson)
}

func (h *NutriotionistHandler) GetClientById(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	rawId := r.PathValue("id")

	id, parseErr := uuid.Parse(rawId)
	if parseErr != nil {
		responseErr := core.ValidationError(
			[]string{"The identifier provided in the URL path is not a valid UUID format."},
		)
		response.WriteJSON(w, responseErr.StatusCode, responseErr)
		return
	}

	user, failure := h.service.GetById(&id, r.Context())
	if failure != nil {
		response.WriteJSON(w, failure.StatusCode, failure)
		return
	}

	apiResponse := &response.ApiResponse[nutritionistDtos.UserDto]{
		Success: true,
		Data: &nutritionistDtos.UserDto{
			ID:           user.ID,
			Firstname:    user.Firstname,
			Lastname:     user.Lastname,
			EmailAddress: user.EmailAddress,
			BirthDate:    user.DateBirth,
		},
	}
	response.WriteJSON(w, http.StatusOK, apiResponse)
}

func (h *NutriotionistHandler) UpdateClientById(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	rawId := r.PathValue("id")
	var dto nutritionistDtos.UpdateNutritionistIdentityRequest

	id, parseErr := uuid.Parse(rawId)
	if parseErr != nil {
		responseErr := core.ValidationError(
			[]string{"The identifier provided in the URL path is not a valid UUID format."},
		)
		response.WriteJSON(w, responseErr.StatusCode, responseErr)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		http.Error(w, fmt.Sprintf("Bad request: %s", err.Error()), http.StatusBadRequest)
		return
	}

	if err := dto.Validate(); err != nil {
		response.WriteJSON(w, err.StatusCode, err)
		return
	}

	failure := h.service.UpdateIdentityInformation(&id, &dto, r.Context())
	if failure != nil {
		response.WriteJSON(w, failure.StatusCode, failure)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *NutriotionistHandler) DeleteClientById(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	rawId := r.PathValue("id")

	id, parseErr := uuid.Parse(rawId)
	if parseErr != nil {
		responseErr := core.ValidationError(
			[]string{"The identifier provided in the URL path is not a valid UUID format."},
		)
		response.WriteJSON(w, responseErr.StatusCode, responseErr)
		return
	}

	err := h.service.DeleteOne(&id, r.Context())
	if err != nil {
		response.WriteJSON(w, err.StatusCode, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
