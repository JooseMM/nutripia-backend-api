package users

import (
	"encoding/json"
	"net/http"

	userModels "github.com/JooseMM/nutripia-backend-api/internal/users/models"
	"github.com/JooseMM/nutripia-backend-api/pkg/core"
)

type IUserHandler interface {
	Create(w http.ResponseWriter, r *http.Request)
}

type UserHandler struct {
	service IUserService
}

func NewUserHandler(service IUserService) IUserHandler {
	return &UserHandler{service}
}

func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var dto userModels.CreateUserDto

	err := json.NewDecoder(r.Body).Decode(&dto)
	if err != nil {
		http.Error(w, "Bad request: "+err.Error(), http.StatusBadRequest)
		return
	}

	validationErr := dto.Validate()
	if validationErr != nil {
		e := core.ValidationError(*validationErr)

		json, jsonErr := json.Marshal(e)
		if jsonErr != nil {
			http.Error(
				w,
				"Error tryinh to serialize an error response",
				http.StatusInternalServerError,
			)
			return
		}

		w.WriteHeader(http.StatusBadRequest)
		w.Write(json)
		return
	}

	id, e := h.service.CreateUser(dto, r.Context())
	if e != nil {
		errJson, err := json.Marshal(id)
		if err != nil {
			http.Error(
				w,
				"Error tryinh to serialize an error response",
				http.StatusInternalServerError,
			)
			return
		}

		w.WriteHeader(http.StatusConflict)
		w.Write(errJson)
		return
	}

	responseJson, responseErr := json.Marshal(id)
	if responseErr != nil {
		http.Error(
			w,
			"Error tryinh to serialize an error response",
			http.StatusInternalServerError,
		)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write(responseJson)
}
