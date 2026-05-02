package clients

import (
	"net/http"

	clientDtos "github.com/JooseMM/nutripia-backend-api/internal/client/dtos"
	"github.com/JooseMM/nutripia-backend-api/pkg/core/valueobject"
	"github.com/JooseMM/nutripia-backend-api/pkg/request"
	"github.com/JooseMM/nutripia-backend-api/pkg/response"
)

const UserIdKey = "userId"
const RoleIdKey = "roleId"

type Handler interface {
	CreateClient(w http.ResponseWriter, r *http.Request)
	GetClientById(w http.ResponseWriter, r *http.Request)
	GetClientByNutritionist(w http.ResponseWriter, r *http.Request)
	DeleteClientById(w http.ResponseWriter, r *http.Request)
	UpdateClientById(w http.ResponseWriter, r *http.Request)
}

type handler struct {
	service ClientManager
}

func NewClientHandler(service ClientManager) Handler {
	return &handler{service}
}

func (h *handler) CreateClient(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	dto, err := request.DecodeJSON[clientDtos.RawCreateClientRequest](w, r)
	if err != nil {
		response.WriteJSON(w, err.StatusCode, err)
		return
	}

	payload, err := dto.ToValueObject()
	if err != nil {
		response.WriteJSON(w, err.StatusCode, err)
		return
	}

	id, err := valueobject.IdentifierFromString(r.PathValue("nutritionistId"))
	if err != nil {
		response.WriteJSON(w, err.StatusCode, err)
		return
	}

	createdId, failure := h.service.CreateUser(r.Context(), *payload, id)
	if failure != nil {
		response.WriteJSON(w, failure.StatusCode, failure)
		return
	}

	apiResponse := &response.ApiResponse[string]{
		Success: true,
		Data:    createdId.String(),
	}
	response.WriteJSON(w, http.StatusCreated, apiResponse)
}

func (h *handler) GetClientByNutritionist(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	id, err := valueobject.IdentifierFromString(r.PathValue("nutritionistId"))
	if err != nil {
		response.WriteJSON(w, err.StatusCode, err)
		return
	}

	clientList, err := h.service.GetByNutritionist(r.Context(), id)
	if err != nil {
		response.WriteJSON(w, err.StatusCode, err)
		return
	}

	dtoList := []clientDtos.ClientDto{}
	for _, c := range clientList {
		dtoList = append(dtoList, c.ToDTO())
	}

	apiResponse := &response.ApiResponse[[]clientDtos.ClientDto]{
		Success: true,
		Data:    dtoList,
	}
	response.WriteJSON(w, http.StatusOK, apiResponse)
}

func (h *handler) GetClientById(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	id, err := valueobject.IdentifierFromString(r.PathValue("clientId"))
	if err != nil {
		response.WriteJSON(w, err.StatusCode, err)
		return
	}

	client, err := h.service.GetById(r.Context(), id)
	if err != nil {
		response.WriteJSON(w, err.StatusCode, err)
		return
	}

	apiResponse := &response.ApiResponse[clientDtos.ClientDto]{
		Success: true,
		Data:    client.ToDTO(),
	}
	response.WriteJSON(w, http.StatusOK, apiResponse)
}

func (h *handler) DeleteClientById(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	id, parseErr := valueobject.IdentifierFromString(r.PathValue("clientId"))
	if parseErr != nil {
		response.WriteJSON(w, parseErr.StatusCode, parseErr)
		return
	}

	err := h.service.Delete(r.Context(), id)
	if err != nil {
		response.WriteJSON(w, err.StatusCode, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *handler) UpdateClientById(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	id, parseErr := valueobject.IdentifierFromString(r.PathValue("clientId"))
	if parseErr != nil {
		response.WriteJSON(w, parseErr.StatusCode, parseErr)
		return
	}

	rawPayload, err := request.DecodeJSON[clientDtos.RawUpdateClientRequest](w, r)
	if err != nil {
		response.WriteJSON(w, err.StatusCode, err)
		return
	}

	payload, err := rawPayload.ToValueObject()
	if err != nil {
		response.WriteJSON(w, err.StatusCode, err)
		return
	}

	err = h.service.Update(r.Context(), id, *payload)
	if err != nil {
		response.WriteJSON(w, err.StatusCode, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
