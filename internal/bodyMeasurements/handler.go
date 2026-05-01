package bodyMeasurement

import (
	"encoding/json"
	"net/http"

	bodyMeasurementDtos "github.com/JooseMM/nutripia-backend-api/internal/bodyMeasurements/dtos"
	"github.com/JooseMM/nutripia-backend-api/internal/clients"
	"github.com/JooseMM/nutripia-backend-api/pkg/core"
	"github.com/JooseMM/nutripia-backend-api/pkg/core/valueobject"
	"github.com/JooseMM/nutripia-backend-api/pkg/response"
	"github.com/google/uuid"
)

type IBodyMeasurementHandler interface {
	Create(w http.ResponseWriter, r *http.Request)
	GetById(w http.ResponseWriter, r *http.Request)
	GetByRange(w http.ResponseWriter, r *http.Request)
	DeleteById(w http.ResponseWriter, r *http.Request)
}

type BodyMeasurementHandler struct {
	ClientService          clients.ClientManager
	bodyMeasurementService BodyMeasurementManager
}

func NewBodyMeasurementHandler(
	bodyMeasurementService BodyMeasurementManager,
	userService clients.ClientManager,
) IBodyMeasurementHandler {
	return &BodyMeasurementHandler{
		ClientService:          userService,
		bodyMeasurementService: bodyMeasurementService,
	}
}

func (h *BodyMeasurementHandler) GetByRange(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	queryParam := r.URL.Query()

	dto, err := bodyMeasurementDtos.NewRangeDate(queryParam.Get("start"), queryParam.Get("end"))
	if err != nil {
		response.WriteJSON(w, err.StatusCode, err)
		return
	}

	id, err := valueobject.IdentifierFromString(r.PathValue("userId"))
	if err != nil {
		response.WriteJSON(w, err.StatusCode, err)
		return
	}

	measurements, err := h.bodyMeasurementService.GetByRange(
		&userId,
		&start,
		&endDate,
		r.Context(),
	)
	if err != nil {
		response.WriteJSON(w, err.StatusCode, err)
		return
	}

	var dtos []bodyMeasurementDtos.BodyMeasurementDto
	for _, bm := range measurements {
		dtos = append(dtos, *bm.ToDTO())
	}
	apiResponse := &response.ApiResponse[[]bodyMeasurementDtos.BodyMeasurementDto]{
		Success: true,
		Data:    &dtos,
	}

	response.WriteJSON(w, http.StatusOK, apiResponse)
}

func (h *BodyMeasurementHandler) Create(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var dto bodyMeasurementDtos.CreateMeasurementDto

	err := json.NewDecoder(r.Body).Decode(&dto)
	if err != nil {
		http.Error(w, "Bad request: "+err.Error(), http.StatusBadRequest)
		return
	}
	if err := dto.Validate(); err != nil {
		response.WriteJSON(w, err.StatusCode, err)
		return
	}

	_, userErr := h.ClientService.GetById(&dto.ClientId, r.Context())
	if userErr != nil {
		response.WriteJSON(w, userErr.StatusCode, userErr)
	}

	createdId, failure := h.bodyMeasurementService.Create(&dto, r.Context())
	if failure != nil {
		response.WriteJSON(w, failure.StatusCode, failure)
		return
	}

	createdIdStr := createdId.String()
	responseJson := &response.ApiResponse[string]{
		Success: true,
		Data:    &createdIdStr,
	}
	response.WriteJSON(w, http.StatusCreated, responseJson)
}

func (h *BodyMeasurementHandler) GetById(w http.ResponseWriter, r *http.Request) {
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

	measurement, err := h.bodyMeasurementService.GetById(&id, r.Context())
	if err != nil {
		response.WriteJSON(w, err.StatusCode, err)
		return
	}

	dto := measurement.ToDTO()
	apiResponse := &response.ApiResponse[bodyMeasurementDtos.BodyMeasurementDto]{
		Success: true,
		Data:    dto,
	}

	response.WriteJSON(w, http.StatusOK, apiResponse)
}

func (h *BodyMeasurementHandler) DeleteById(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	rawId := r.PathValue("id")

	id, parseErr := uuid.Parse(rawId)
	if parseErr != nil {
		errResponse := core.ValidationError(
			[]string{"The identifier provided in the URL path is not a valid UUID format."},
		)
		response.WriteJSON(w, errResponse.StatusCode, errResponse)
		return
	}

	if err := h.bodyMeasurementService.DeleteOne(&id, r.Context()); err != nil {
		response.WriteJSON(w, err.StatusCode, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
