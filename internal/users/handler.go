package users

import (
	"encoding/json"
	"net/http"

	"github.com/JooseMM/nutripia-backend-api/internal/users/types"
	"github.com/JooseMM/nutripia-backend-api/internal/users/types/dtos"
	"github.com/JooseMM/nutripia-backend-api/pkg/core"
	"github.com/JooseMM/nutripia-backend-api/pkg/response"
	"github.com/google/uuid"
)

type IUserHandler interface {
	Create(w http.ResponseWriter, r *http.Request)
	GetById(w http.ResponseWriter, r *http.Request)
	DeleteById(w http.ResponseWriter, r *http.Request)
}

type UserHandler struct {
	service IUserService
}

func NewMeasurementHandler(service IUserService) IUserHandler {
	return &UserHandler{service}
}

func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var dto userDtos.UserIdentityDto

	err := json.NewDecoder(r.Body).Decode(&dto)
	if err != nil {
		http.Error(w, "Bad request: "+err.Error(), http.StatusBadRequest)
		return
	}

	validationErr := dto.Validate()
	if validationErr != nil {
		json, jsonErr := json.Marshal(validationErr)
		if jsonErr != nil {
			http.Error(
				w,
				"Error trying to serialize a response",
				http.StatusInternalServerError,
			)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write(json)
		return
	}

	id, failure := h.service.CreateUser(&dto, r.Context())
	if failure != nil {
		errJson, err := json.Marshal(failure)
		if err != nil {
			http.Error(
				w,
				"Error trying to serialize a response",
				http.StatusInternalServerError,
			)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(int(failure.StatusCode))
		w.Write(errJson)
		return
	}

	idStr := id.String()
	response := &response.ApiResponse[string]{
		Success: true,
		Data:    &idStr,
	}
	responseJson, responseErr := json.Marshal(response)
	if responseErr != nil {
		http.Error(
			w,
			"Error trying to serialize a response",
			http.StatusInternalServerError,
		)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(responseJson)
}

func (h *UserHandler) GetById(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	rawId := r.PathValue("id")

	id, parseErr := uuid.Parse(rawId)
	if parseErr != nil {
		responseErr := core.ValidationError(
			[]string{"The identifier provided in the URL path is not a valid UUID format."},
		)
		jsonResponse, jsonErr := json.Marshal(responseErr)
		if jsonErr != nil {
			http.Error(
				w,
				"Error trying to serialize a response",
				http.StatusInternalServerError,
			)
			return
		}

		w.WriteHeader(int(responseErr.StatusCode))
		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonResponse)
		return
	}

	user, err := h.service.GetById(&id, r.Context())
	if err != nil {
		errResponse, e := json.Marshal(err)
		if e != nil {
			http.Error(
				w,
				"Error trying to serialize a response",
				http.StatusInternalServerError,
			)
			return
		}
		w.WriteHeader(int(err.StatusCode))
		w.Header().Set("Content-Type", "application/json")
		w.Write(errResponse)
		return
	}

	apiResponse := &response.ApiResponse[userTypes.User]{
		Success: true,
		Data:    user,
	}
	userResponse, responseErr := json.Marshal(apiResponse)
	if responseErr != nil {
		http.Error(
			w,
			"Error trying to serialize a response",
			http.StatusInternalServerError,
		)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(userResponse)
}

func (h *UserHandler) DeleteById(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	rawId := r.PathValue("id")

	id, parseErr := uuid.Parse(rawId)
	if parseErr != nil {
		responseErr := core.ValidationError(
			[]string{"The identifier provided in the URL path is not a valid UUID format."},
		)
		jsonResponse, jsonErr := json.Marshal(responseErr)
		if jsonErr != nil {
			http.Error(
				w,
				"Error trying to serialize a response",
				http.StatusInternalServerError,
			)
			return
		}

		w.WriteHeader(int(responseErr.StatusCode))
		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonResponse)
		return
	}

	err := h.service.DeleteOne(&id, r.Context())
	if err != nil {
		errResponse, e := json.Marshal(err)
		if e != nil {
			http.Error(
				w,
				"Error trying to serialize a response",
				http.StatusInternalServerError,
			)
			return
		}
		w.WriteHeader(int(err.StatusCode))
		w.Header().Set("Content-Type", "application/json")
		w.Write(errResponse)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
