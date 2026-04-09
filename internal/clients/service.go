package clients

import (
	"context"
	"time"

	clientTypes "github.com/JooseMM/nutripia-backend-api/internal/clients/types"
	"github.com/JooseMM/nutripia-backend-api/internal/clients/types/dtos"
	"github.com/JooseMM/nutripia-backend-api/pkg/core"
	"github.com/google/uuid"
)

type IClientService interface {
	CreateUser(
		dto *clientDtos.CreateClientRequest,
		ctx context.Context,
	) (*uuid.UUID, *core.BaseError)
	GetById(id *uuid.UUID, ctx context.Context) (*clientTypes.Client, *core.BaseError)
	DeleteOne(id *uuid.UUID, ctx context.Context) *core.BaseError
	UpdateIdentityInformation(
		id *uuid.UUID,
		userDto *clientDtos.UpdateClientRequest,
		ctx context.Context,
	) *core.BaseError
}

type UserService struct {
	Repo IClientRepository
}

func NewClientService(repo IClientRepository) IClientService {
	return &UserService{repo}
}

func (u *UserService) CreateUser(
	userDto *clientDtos.CreateClientRequest,
	ctx context.Context,
) (*uuid.UUID, *core.BaseError) {

	foundEmailOwner, queryEmailErr := u.Repo.GetByEmailAddress(ctx, userDto.EmailAddress)
	if queryEmailErr != nil && queryEmailErr.ErrorCode != string(CLIENT_NOT_FOUND) {
		return nil, queryEmailErr
	}
	if foundEmailOwner != nil {
		return nil, UserEmailAlreadyExisting(userDto.EmailAddress)
	}

	now := time.Now()
	user := &clientTypes.Client{
		ID: uuid.New(),
		ClientIdentity: clientTypes.ClientIdentity{
			Firstname:    userDto.Firstname,
			Lastname:     userDto.Lastname,
			EmailAddress: userDto.EmailAddress,
			DateBirth:    userDto.BirthDate,
		},
		TrackingInformation: clientTypes.TrackingInformation{
			CreatedAt: now,
			UpdatedAt: now,
		},
	}

	creationErr := u.Repo.Create(ctx, user)
	if creationErr != nil {
		return nil, creationErr
	}

	return &user.ID, nil
}

func (u *UserService) GetById(
	id *uuid.UUID,
	ctx context.Context,
) (*clientTypes.Client, *core.BaseError) {

	foundUser, unexpectedErr := u.Repo.GetById(ctx, id)
	if unexpectedErr != nil {
		return nil, unexpectedErr
	}

	return foundUser, nil
}

func (u *UserService) DeleteOne(
	id *uuid.UUID,
	ctx context.Context,
) *core.BaseError {
	err := u.Repo.Delete(ctx, id)
	if err != nil {
		return err
	}

	return nil
}

func (u *UserService) UpdateIdentityInformation(
	id *uuid.UUID,
	userDto *clientDtos.UpdateClientRequest,
	ctx context.Context,
) *core.BaseError {
	foundUser, unexpectedErr := u.Repo.GetById(ctx, id)
	if unexpectedErr != nil {
		return unexpectedErr
	}

	foundEmailOwner, queryEmailErr := u.Repo.GetByEmailAddress(ctx, userDto.EmailAddress)
	if queryEmailErr != nil && queryEmailErr.ErrorCode != string(CLIENT_NOT_FOUND) {
		return queryEmailErr
	}
	if foundEmailOwner != nil && foundEmailOwner.ID != *id {
		return UserEmailAlreadyExisting(userDto.EmailAddress)
	}

	foundUser.Firstname = userDto.Firstname
	foundUser.Lastname = userDto.Lastname
	foundUser.EmailAddress = userDto.EmailAddress
	foundUser.DateBirth = userDto.BirthDate

	foundUser.UpdatedAt = time.Now()

	return u.Repo.Update(ctx, foundUser)
}
