package clients

import (
	"context"

	clientDtos "github.com/JooseMM/nutripia-backend-api/internal/clients/dtos"
	"github.com/JooseMM/nutripia-backend-api/pkg/core"
	"github.com/JooseMM/nutripia-backend-api/pkg/core/valueobject"
)

type ClientManager interface {
	CreateUser(
		ctx context.Context,
		dto clientDtos.CreateClientRequest,
		nutritionistOwnerId valueobject.Identifier,
	) (valueobject.Identifier, *core.BaseError)
	GetById(ctx context.Context, id valueobject.Identifier) (Client, *core.BaseError)
	GetByNutritionist(
		ctx context.Context,
		userId valueobject.Identifier,
	) ([]Client, *core.BaseError)
	Delete(ctx context.Context, id valueobject.Identifier) *core.BaseError
	Update(
		ctx context.Context,
		id valueobject.Identifier,
		dto clientDtos.UpdateClientRequest,
	) *core.BaseError
}

type service struct {
	Repo Repository
}

func NewService(repo Repository) ClientManager {
	return &service{repo}
}

func (u *service) CreateUser(
	ctx context.Context,
	dto clientDtos.CreateClientRequest,
	ownerId valueobject.Identifier,
) (valueobject.Identifier, *core.BaseError) {
	isFound, err := u.Repo.IsEmailTaken(ctx, dto.EmailAddress)
	if err != nil {
		return nil, err
	}
	if isFound {
		return nil, UserEmailAlreadyExisting(dto.EmailAddress.String())
	}

	client := NewClient(dto, ownerId)

	creationErr := u.Repo.Create(ctx, client)
	if creationErr != nil {
		return nil, creationErr
	}

	return client.Id(), nil
}

func (u *service) GetById(
	ctx context.Context,
	id valueobject.Identifier,
) (Client, *core.BaseError) {
	client, err := u.Repo.GetById(ctx, id)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func (u *service) GetByNutritionist(
	ctx context.Context,
	userId valueobject.Identifier,
) ([]Client, *core.BaseError) {
	clientList, unexpectedErr := u.Repo.GetAllByNutritionist(ctx, userId)
	if unexpectedErr != nil {
		return nil, unexpectedErr
	}
	if clientList == nil {
		return []Client{}, nil
	}

	return clientList, nil
}

func (u *service) Delete(
	ctx context.Context,
	id valueobject.Identifier,
) *core.BaseError {
	err := u.Repo.Delete(ctx, id)
	if err != nil {
		return err
	}

	return nil
}

func (u *service) Update(
	ctx context.Context,
	id valueobject.Identifier,
	dto clientDtos.UpdateClientRequest,
) *core.BaseError {
	client, err := u.Repo.GetById(ctx, id)
	if err != nil {
		return err
	}

	isFound, err := u.Repo.IsEmailTaken(ctx, dto.EmailAddress)
	if err != nil {
		return err
	}
	if isFound {
		return UserEmailAlreadyExisting(dto.EmailAddress.String())
	}

	client.Update(dto)
	return u.Repo.Update(ctx, client)
}
