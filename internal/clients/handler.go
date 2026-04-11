package clients

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/JooseMM/nutripia-backend-api/internal/clients/types/dtos"
	"github.com/JooseMM/nutripia-backend-api/pkg/core"
	"github.com/JooseMM/nutripia-backend-api/pkg/response"
	"github.com/google/uuid"
)

type IClientHandler interface {
	CreateClient(w http.ResponseWriter, r *http.Request)
	GetClientById(w http.ResponseWriter, r *http.Request)
	GetClientByNutritionist(w http.ResponseWriter, r *http.Request)
	DeleteClientById(w http.ResponseWriter, r *http.Request)
	UpdateClientById(w http.ResponseWriter, r *http.Request)
}

type ClientHandler struct {
	service IClientService
}

func NewClientHandler(service IClientService) IClientHandler {
	return &ClientHandler{service}
}

func (h *ClientHandler) CreateClient(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var dto clientDtos.CreateClientRequest

	err := json.NewDecoder(r.Body).Decode(&dto)
	if err != nil {
		responseErr := core.ValidationError([]string{err.Error()})
		response.WriteJSON(w, responseErr.StatusCode, responseErr)
		return
	}

	if validationErr := dto.Validate(); validationErr != nil {
		response.WriteJSON(w, validationErr.StatusCode, validationErr)
		return
	}

	id, failure := h.service.CreateUser(&dto, r.Context())
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

func (h *ClientHandler) GetClientByNutritionist(w http.ResponseWriter, r *http.Request) {
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

	user, err := h.service.GetById(&id, r.Context())
	if err != nil {
		response.WriteJSON(w, err.StatusCode, err)
		return
	}

	apiResponse := &response.ApiResponse[clientDtos.ClientDto]{
		Success: true,
		Data: &clientDtos.ClientDto{
			ID:           user.ID,
			Firstname:    user.Firstname,
			Lastname:     user.Lastname,
			EmailAddress: user.EmailAddress,
			BirthDate:    user.DateBirth,
		},
	}
	response.WriteJSON(w, http.StatusOK, apiResponse)
}

func (h *ClientHandler) GetClientById(w http.ResponseWriter, r *http.Request) {
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

	apiResponse := &response.ApiResponse[clientDtos.ClientDto]{
		Success: true,
		Data: &clientDtos.ClientDto{
			ID:           user.ID,
			Firstname:    user.Firstname,
			Lastname:     user.Lastname,
			EmailAddress: user.EmailAddress,
			BirthDate:    user.DateBirth,
		},
	}
	response.WriteJSON(w, http.StatusOK, apiResponse)
}

func (h *ClientHandler) DeleteClientById(w http.ResponseWriter, r *http.Request) {
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

func (h *ClientHandler) UpdateClientById(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	rawId := r.PathValue("id")
	var dto clientDtos.UpdateClientRequest

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
		w.Write(jsonResponse)
		return
	}

	err := json.NewDecoder(r.Body).Decode(&dto)
	if err != nil {
		http.Error(w, fmt.Sprintf("Bad request: %s", err.Error()), http.StatusBadRequest)
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

	failure := h.service.UpdateIdentityInformation(&id, &dto, r.Context())
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

	w.WriteHeader(http.StatusNoContent)
}
