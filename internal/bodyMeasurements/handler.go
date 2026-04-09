package bodyMeasurement

import (
	"encoding/json"
	"net/http"
	"time"

	bodyMeasurementDtos "github.com/JooseMM/nutripia-backend-api/internal/bodyMeasurements/types/dtos"
	"github.com/JooseMM/nutripia-backend-api/internal/users"
	userTypes "github.com/JooseMM/nutripia-backend-api/internal/users/types"
	"github.com/JooseMM/nutripia-backend-api/pkg/core"
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
	userService            users.IUserService
	bodyMeasurementService IBodyMeasurementService
}

func NewBodyMeasurementHandler(
	bodyMeasurementService IBodyMeasurementService,
	userService users.IUserService,
) IBodyMeasurementHandler {
	return &BodyMeasurementHandler{
		userService:            userService,
		bodyMeasurementService: bodyMeasurementService,
	}
}

func (h *BodyMeasurementHandler) GetByRange(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	layout := "2000-01-01"
	queryParam := r.URL.Query()

	startDate, startErr := time.Parse(layout, queryParam.Get("start"))
	if startErr != nil {
		startDate = time.Now().AddDate(-1, 0, 0)
	}

	endDate, endErr := time.Parse(layout, queryParam.Get("end"))
	if endErr != nil {
		endDate = time.Now()
	}

	rawUserId := r.PathValue("userId")
	userId, parseErr := uuid.Parse(rawUserId)
	if parseErr != nil {
		responseErr := core.ValidationError(
			[]string{"The identifier provided in the URL path is not a valid UUID format."},
		)
		response.WriteJSON(w, responseErr.StatusCode, responseErr)
		return
	}

	measurements, err := h.bodyMeasurementService.GetByRange(
		&userId,
		&startDate,
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

	user, userErr := h.userService.GetById(&dto.UserId, r.Context())
	if userErr != nil {
		response.WriteJSON(w, userErr.StatusCode, userErr)
	}

	if user.Role != userTypes.CLIENT {
		e := BodyMeasurementNotAClient(user.ID)
		response.WriteJSON(w, e.StatusCode, e)
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
