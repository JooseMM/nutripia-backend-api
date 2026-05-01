package clients

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/JooseMM/nutripia-backend-api/pkg/core"
	"github.com/JooseMM/nutripia-backend-api/pkg/core/valueobject"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type clients struct {
	Id        uuid.UUID `gorm:"type:uuid;primaryKey"`
	OwnerId   uuid.UUID `gorm:"column:owner_id;type:uuid;index"`
	Firstname string    `gorm:"column:firstname;not null"`
	Lastname  string    `gorm:"column:lastname;not null"`
	Email     string    `gorm:"column:email;uniqueIndex;not null"`
	BirthDate time.Time `gorm:"column:birth_date"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

func (clients) TableName() string {
	return "clients"
}

func (d *clients) toEntity() (Client, *core.BaseError) {
	var errList []string

	id := valueobject.IdentifierFromValue(d.Id)

	ownerId := valueobject.IdentifierFromValue(d.OwnerId)

	name, err := valueobject.NewName(d.Firstname, d.Lastname)
	if err != nil {
		errList = append(errList, err.Details...)
	}

	email, err := valueobject.NewEmailAddress(d.Email)
	if err != nil {
		errList = append(errList, err.Details...)
	}

	birthDate, err := valueobject.NewBirthDate(d.BirthDate)
	if err != nil {
		errList = append(errList, err.Details...)
	}

	created, err := valueobject.NewTrackerFromTime(&d.CreatedAt)
	if err != nil {
		errList = append(errList, err.Details...)
	}

	updated, err := valueobject.NewTrackerFromTime(&d.UpdatedAt)
	if err != nil {
		errList = append(errList, err.Details...)
	}

	if len(errList) > 0 {
		return nil, core.CorrupetedDatabase(errList)
	}

	return &entity{
		id:        id,
		name:      name,
		email:     email,
		birthDate: birthDate,
		ownerId:   ownerId,

		createdAt: created,
		updateAt:  updated,
	}, nil
}

type repo struct {
	db *gorm.DB
}

type Repository interface {
	GetAllByNutritionist(
		ctx context.Context,
		userId valueobject.Identifier,
	) ([]Client, *core.BaseError)
	GetById(ctx context.Context, id valueobject.Identifier) (Client, *core.BaseError)
	IsEmailTaken(
		ctx context.Context,
		emalAddress valueobject.Emailer,
	) (bool, *core.BaseError)
	GetByEmailAddress(
		ctx context.Context,
		emalAddress valueobject.Emailer,
	) (Client, *core.BaseError)
	Create(ctx context.Context, u Client) *core.BaseError
	Update(ctx context.Context, u Client) *core.BaseError
	Delete(ctx context.Context, id valueobject.Identifier) *core.BaseError
}

func NewRepository(db *gorm.DB) (Repository, *core.BaseError) {
	if err := db.AutoMigrate(&clients{}); err != nil {
		return nil, core.UnexpectedError(err.Error())
	}
	return &repo{db}, nil
}

func (r *repo) GetAllByNutritionist(
	ctx context.Context,
	userId valueobject.Identifier,
) ([]Client, *core.BaseError) {
	var rawList []clients

	result := r.db.WithContext(ctx).Where("owner_id = ?", userId.Key()).Find(&rawList)
	if result.Error != nil {
		fmt.Printf("error: %s", result.Error)
		return nil, core.UnexpectedError(result.Error.Error())
	}

	var list []Client

	for _, raw := range rawList {
		c, err := raw.toEntity()
		if err != nil {
			return nil, err
		}

		list = append(list, c)
	}

	return list, nil
}

func (r *repo) GetById(
	ctx context.Context,
	id valueobject.Identifier,
) (Client, *core.BaseError) {
	var raw clients

	result := r.db.WithContext(ctx).First(&raw, "id = ?", id.Key())
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, UserNotFound(id.String())
		}
		return nil, core.UnexpectedError(result.Error.Error())
	}

	client, err := raw.toEntity()
	if err != nil {
		return nil, err
	}

	return client, nil
}

func (r *repo) IsEmailTaken(
	ctx context.Context,
	emalAddress valueobject.Emailer,
) (bool, *core.BaseError) {
	var exists bool

	err := r.db.WithContext(ctx).
		Model(&clients{}).
		Select("count(*) > 0").
		Where("email = ?", emalAddress.String()).
		Find(&exists).
		Error

	if err != nil {
		return false, core.UnexpectedError(err.Error())
	}
	return exists, nil
}

func (r *repo) GetByEmailAddress(
	ctx context.Context,
	emailAddress valueobject.Emailer,
) (Client, *core.BaseError) {
	var raw clients

	result := r.db.WithContext(ctx).First(&raw, "email = ?", emailAddress.String())
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, UserNotFound(emailAddress.String())
		}
		return nil, core.UnexpectedError(result.Error.Error())
	}

	client, err := raw.toEntity()
	if err != nil {
		return nil, err
	}

	return client, nil
}

func (r *repo) Create(ctx context.Context, client Client) *core.BaseError {
	data := client.ToDB()
	result := r.db.WithContext(ctx).Create(&data)
	if result.Error != nil {
		return core.UnexpectedError(result.Error.Error())
	}
	return nil
}

func (r *repo) Update(ctx context.Context, client Client) *core.BaseError {
	data := client.ToDB()
	result := r.db.WithContext(ctx).Save(&data)
	if result.Error != nil {
		return core.UnexpectedError(result.Error.Error())
	}
	return nil
}

func (r *repo) Delete(ctx context.Context, id valueobject.Identifier) *core.BaseError {
	result := r.db.WithContext(ctx).Delete(&clients{}, id.Key())
	if result.Error != nil {
		return core.UnexpectedError(result.Error.Error())
	}
	if result.RowsAffected == 0 {
		return UserNotFound(id.String())
	}

	return nil
}
