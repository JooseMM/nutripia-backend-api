package bodyMeasurement

import (
	"net/http"

	bodyMeasurementDtos "github.com/JooseMM/nutripia-backend-api/internal/bodyMeasurement/dtos"
	"github.com/JooseMM/nutripia-backend-api/internal/client"
	"github.com/JooseMM/nutripia-backend-api/pkg/core/valueobject"
	"github.com/JooseMM/nutripia-backend-api/pkg/request"
	"github.com/JooseMM/nutripia-backend-api/pkg/response"
)

type Handler interface {
	Create(w http.ResponseWriter, r *http.Request)
	GetById(w http.ResponseWriter, r *http.Request)
	GetByRange(w http.ResponseWriter, r *http.Request)
	DeleteById(w http.ResponseWriter, r *http.Request)
}

type handler struct {
	ClientService          clients.ClientManager
	bodyMeasurementService BodyMeasurementManager
}

func NewHandler(
	bodyMeasurementService BodyMeasurementManager,
	userService clients.ClientManager,
) Handler {
	return &handler{
		ClientService:          userService,
		bodyMeasurementService: bodyMeasurementService,
	}
}

func (h *handler) GetByRange(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	queryParam := r.URL.Query()

	dto, err := bodyMeasurementDtos.NewRangeDate(queryParam.Get("start"), queryParam.Get("end"))
	if err != nil {
		response.WriteJSON(w, err.StatusCode, err)
		return
	}

	clientId, err := valueobject.IdentifierFromString(r.PathValue("userId"))
	if err != nil {
		response.WriteJSON(w, err.StatusCode, err)
		return
	}

	measurements, err := h.bodyMeasurementService.GetByRange(
		r.Context(),
		clientId,
		*dto,
	)
	if err != nil {
		response.WriteJSON(w, err.StatusCode, err)
		return
	}

	var dtos []bodyMeasurementDtos.BodyMeasurementDto
	for _, bm := range measurements {
		dtos = append(dtos, bm.ToDTO())
	}
	apiResponse := &response.ApiResponse[[]bodyMeasurementDtos.BodyMeasurementDto]{
		Success: true,
		Data:    dtos,
	}

	response.WriteJSON(w, http.StatusOK, apiResponse)
}

func (h *handler) Create(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	dto, err := request.DecodeJSON[bodyMeasurementDtos.RawCreateMeasurement](w, r)
	if err != nil {
		response.WriteJSON(w, err.StatusCode, err)
		return
	}

	payload, err := dto.ToValueObject()
	if err != nil {
		response.WriteJSON(w, err.StatusCode, err)
	}

	clientId, err := valueobject.IdentifierFromString(r.PathValue("clientId"))
	if err != nil {
		response.WriteJSON(w, err.StatusCode, err)
		return
	}

	_, clientErr := h.ClientService.GetById(r.Context(), clientId)
	if clientErr != nil {
		response.WriteJSON(w, clientErr.StatusCode, clientErr)
		return
	}

	createdId, failure := h.bodyMeasurementService.Create(r.Context(), *payload, clientId)
	if failure != nil {
		response.WriteJSON(w, failure.StatusCode, failure)
		return
	}

	createdIdStr := createdId.String()
	responseJson := &response.ApiResponse[string]{
		Success: true,
		Data:    createdIdStr,
	}
	response.WriteJSON(w, http.StatusCreated, responseJson)
}

func (h *handler) GetById(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	id, err := valueobject.IdentifierFromString(r.PathValue("id"))
	if err != nil {
		response.WriteJSON(w, err.StatusCode, err)
		return
	}

	measurement, err := h.bodyMeasurementService.GetById(r.Context(), id)
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

func (h *handler) DeleteById(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	id, err := valueobject.IdentifierFromString(r.PathValue("id"))
	if err != nil {
		response.WriteJSON(w, err.StatusCode, err)
		return
	}

	if err := h.bodyMeasurementService.DeleteOne(r.Context(), id); err != nil {
		response.WriteJSON(w, err.StatusCode, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
