package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/JooseMM/nutripia-backend-api/internal/users/models"
	"github.com/JooseMM/nutripia-backend-api/internal/users/models/value_objects"
	"github.com/JooseMM/nutripia-backend-api/internal/users/repository/interfaces"
	"github.com/google/uuid"
)

type UserHandler struct {
	Repo interfaces.UserRepository
}

func NewUserHandler(repo interfaces.UserRepository) *UserHandler {
	return &UserHandler{repo}
}

func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var dto models.CreateUserDto

	err := json.NewDecoder(r.Body).Decode(&dto)
	if err != nil {
		http.Error(w, "Bad request: "+err.Error(), http.StatusBadRequest)
		return
	}

	errors := dto.Validate()
	if errors != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
	}

	now := time.Now()
	user := &models.User{
		UserIdentity: value_objects.UserIdentity{
			ID:           uuid.New(),
			Firstname:    dto.Firstname,
			Lastname:     dto.Lastname,
			EmailAddress: dto.EmailAddress,
			DateBirth:    dto.BirthDate,
		},
		AuthenticationInformation: value_objects.AuthenticationInformation{},
		TrackingInformation: value_objects.TrackingInformation{
			CreatedAt: now,
			UpdatedAt: now,
		},
	}

	h.Repo.Create(r.Context(), user)

	w.WriteHeader(http.StatusCreated)
}
