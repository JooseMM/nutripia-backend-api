package nutritionist

import (
	"net/http"

	nutritionistDtos "github.com/JooseMM/nutripia-backend-api/internal/nutritionist/dtos"
	"github.com/JooseMM/nutripia-backend-api/pkg/core"
	"github.com/JooseMM/nutripia-backend-api/pkg/core/valueobject"
	"github.com/JooseMM/nutripia-backend-api/pkg/response"
)

type Handler interface {
	GetNutritionistById(w http.ResponseWriter, r *http.Request)
	UpdateNutritionistById(w http.ResponseWriter, r *http.Request)
	DeleteNutritionistById(w http.ResponseWriter, r *http.Request)
}

type handler struct {
	service NutritionistManager
}

func NewNutritionistHandler(service NutritionistManager) Handler {
	return &handler{service}
}

func (h *handler) GetNutritionistById(w http.ResponseWriter, r *http.Request) {
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

func (h *handler) UpdateNutritionistById(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var payload nutritionistDtos.UpdateNutritionist

	id, idErr := valueobject.IdentifierFromString(r.PathValue("id"))
	if idErr != nil {
		response.WriteJSON(w, idErr.StatusCode, idErr)
		return
	}

	payload, errMessage := core.DecodeJSON[nutritionistDtos.UpdateNutritionist](w, r)
	if errMessage != nil {
		res := core.ValidationError(errMessage)
		response.WriteJSON(w, http.StatusBadRequest, res)
		return
	}

	if err := h.service.Update(r.Context(), id, payload); err != nil {
		response.WriteJSON(w, err.StatusCode, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *handler) DeleteNutritionistById(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	id, idErr := valueobject.IdentifierFromString(r.PathValue("id"))
	if idErr != nil {
		response.WriteJSON(w, idErr.StatusCode, idErr)
		return
	}

	if err := h.service.Delete(r.Context(), id); err != nil {
		response.WriteJSON(w, err.StatusCode, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
