package nutritionist

import (
	"net/http"

	nutritionistDtos "github.com/JooseMM/nutripia-backend-api/internal/nutritionist/dtos"
	"github.com/JooseMM/nutripia-backend-api/pkg/core/valueobject"
	"github.com/JooseMM/nutripia-backend-api/pkg/response"
)

type INutritionistHandler interface {
	GetNutritionistById(w http.ResponseWriter, r *http.Request)
	// UpdateNutritionistById(w http.ResponseWriter, r *http.Request)
	// DeleteNutritionistById(w http.ResponseWriter, r *http.Request)
}

type NutriotionistHandler struct {
	service INutritionistService
}

func NewNutritionistHandler(service INutritionistService) INutritionistHandler {
	return &NutriotionistHandler{service}
}

func (h *NutriotionistHandler) GetNutritionistById(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	rawId := r.PathValue("id")

	id, parseErr := valueobject.IdentifierFromString(rawId)
	if parseErr != nil {
		response.WriteJSON(w, parseErr.StatusCode, parseErr)
		return
	}

	user, failure := h.service.GetById(r.Context(), id)
	if failure != nil {
		response.WriteJSON(w, failure.StatusCode, failure)
		return
	}

	apiResponse := &response.ApiResponse[nutritionistDtos.Nutritionist]{
		Success: true,
		Data:    user.ToDTO(),
	}
	response.WriteJSON(w, http.StatusOK, apiResponse)
}

// func (h *NutriotionistHandler) UpdateNutritionistById(w http.ResponseWriter, r *http.Request) {
// 	defer r.Body.Close()
//
// 	rawId := r.PathValue("id")
// 	var dto nutritionistDtos.UpdateNutritionistIdentityRequest
//
// 	id, parseErr := uuid.Parse(rawId)
// 	if parseErr != nil {
// 		responseErr := core.ValidationError(
// 			[]string{"The identifier provided in the URL path is not a valid UUID format."},
// 		)
// 		response.WriteJSON(w, responseErr.StatusCode, responseErr)
// 		return
// 	}
//
// 	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
// 		http.Error(w, fmt.Sprintf("Bad request: %s", err.Error()), http.StatusBadRequest)
// 		return
// 	}
//
// 	if err := dto.Validate(); err != nil {
// 		response.WriteJSON(w, err.StatusCode, err)
// 		return
// 	}
//
// 	ctx := r.Context()
// 	failure := h.service.UpdateIdentityInformation(&id, &dto, &ctx)
// 	if failure != nil {
// 		response.WriteJSON(w, failure.StatusCode, failure)
// 		return
// 	}
//
// 	w.WriteHeader(http.StatusNoContent)
// }
//
// func (h *NutriotionistHandler) DeleteNutritionistById(w http.ResponseWriter, r *http.Request) {
// 	defer r.Body.Close()
// 	rawId := r.PathValue("id")
//
// 	id, parseErr := uuid.Parse(rawId)
// 	if parseErr != nil {
// 		responseErr := core.ValidationError(
// 			[]string{"The identifier provided in the URL path is not a valid UUID format."},
// 		)
// 		response.WriteJSON(w, responseErr.StatusCode, responseErr)
// 		return
// 	}
//
// 	ctx := r.Context()
// 	err := h.service.DeleteOne(&id, &ctx)
// 	if err != nil {
// 		response.WriteJSON(w, err.StatusCode, err)
// 		return
// 	}
//
// 	w.WriteHeader(http.StatusNoContent)
// }
