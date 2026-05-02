package nutritionist

import (
	"context"
	"errors"
	"time"

	"github.com/JooseMM/nutripia-backend-api/pkg/core"
	"github.com/JooseMM/nutripia-backend-api/pkg/core/valueobject"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type nutritionists struct {
	Id               uuid.UUID `gorm:"column:id;type:uuid;primaryKey"`
	Firstname        string    `gorm:"column:firstname;not null"`
	Lastname         string    `gorm:"column:lastname;not null"`
	IsEmailConfirmed bool      `gorm:"column:is_email_confirmed;not null"`
	Email            string    `gorm:"column:email;uniqueIndex;not null"`
	BirthDate        time.Time `gorm:"column:birth_date"`
	Password         string    `gorm:"column:password;not null"`
	Rut              string    `gorm:"column:rut;uniqueIndex;not null"`
	CreatedAt        time.Time `gorm:"column:created_at"`
	UpdatedAt        time.Time `gorm:"column:updated_at"`
}

func fromEntity(entity Nutritionist) nutritionists {
	return nutritionists{
		Id:               entity.Id().Key(),
		Firstname:        entity.Name().Firstname(),
		Lastname:         entity.Name().Lastname(),
		Email:            entity.EmailAddress().String(),
		IsEmailConfirmed: entity.IsEmailConfirmed(),
		BirthDate:        entity.BirthDate().ToTime(),
		Password:         entity.Password().String(),
		Rut:              entity.RUT().ToString(),
		UpdatedAt:        entity.UpdatedAt().ToTime(),
	}
}

func (e *nutritionists) toEntity() (Nutritionist, *core.BaseError) {
	var errList []string

	name, err := valueobject.NewFullName(e.Firstname, e.Lastname)
	if err != nil {
		errList = append(errList, err.Details...)
	}

	emailAddress, err := valueobject.NewEmailAddress(e.Email)
	if err != nil {
		errList = append(errList, err.Details...)
	}

	if err != nil {
		errList = append(errList, err.Details...)
	}

	birthDate, err := valueobject.NewBirthDate(e.BirthDate)
	if err != nil {
		errList = append(errList, err.Details...)
	}
	password, err := valueobject.PasswordFromDB(e.Password)
	if err != nil {
		errList = append(errList, err.Details...)
	}

	rut, err := valueobject.NewRUT(e.Rut)
	if err != nil {
		errList = append(errList, err.Details...)
	}

	createdAt, err := valueobject.NewTrackerFromTime(&e.CreatedAt)
	if err != nil {
		errList = append(errList, err.Details...)
	}

	updatedAt, err := valueobject.NewTrackerFromTime(&e.UpdatedAt)
	if err != nil {
		errList = append(errList, err.Details...)
	}

	if len(errList) > 0 {
		return nil, core.ValidationError(errList)
	}

	return FromDB(
		valueobject.IdentifierFromValue(e.Id),
		name,
		emailAddress,
		e.IsEmailConfirmed,
		birthDate,
		password,
		rut,
		createdAt,
		updatedAt,
	), nil
}

func (nutritionists) TableName() string {
	return "nutritionists"
}

type repository struct {
	db *gorm.DB
}

type Repository interface {
	GetAll(ctx context.Context) ([]Nutritionist, *core.BaseError)
	GetById(ctx context.Context, id valueobject.Identifier) (Nutritionist, *core.BaseError)
	IsEmailTaken(
		ctx context.Context,
		emalAddress valueobject.Emailer,
	) (bool, *core.BaseError)
	GetByEmailAddress(
		ctx context.Context,
		emalAddress valueobject.Emailer,
	) (Nutritionist, *core.BaseError)
	Create(ctx context.Context, u Nutritionist) *core.BaseError
	Update(ctx context.Context, u Nutritionist) *core.BaseError
	Delete(ctx context.Context, id valueobject.Identifier) *core.BaseError
}

func NewRepository(db *gorm.DB) (Repository, *core.BaseError) {
	if err := db.AutoMigrate(&nutritionists{}); err != nil {
		return nil, core.UnexpectedError(err.Error())
	}

	return &repository{db}, nil
}

func (r *repository) IsEmailTaken(
	ctx context.Context,
	emalAddress valueobject.Emailer,
) (bool, *core.BaseError) {
	var exists bool

	err := r.db.WithContext(ctx).
		Model(&nutritionists{}).
		Select("count(*) > 0").
		Where("email = ?", emalAddress.String()).
		Find(&exists).
		Error

	if err != nil {
		return false, core.UnexpectedError(err.Error())
	}
	return exists, nil
}

func (r *repository) GetAll(
	ctx context.Context,
) ([]Nutritionist, *core.BaseError) {
	var rawList []*nutritionists
	var list []Nutritionist

	result := r.db.WithContext(ctx).Find(&rawList)
	if result.Error != nil {
		return nil, core.UnexpectedError(result.Error.Error())
	}

	for _, raw := range rawList {
		entity, err := raw.toEntity()
		if err != nil {
			return nil, err
		}
		list = append(list, entity)
	}

	return list, nil
}

func (r *repository) GetById(
	ctx context.Context,
	id valueobject.Identifier,
) (Nutritionist, *core.BaseError) {
	var raw nutritionists

	result := r.db.WithContext(ctx).First(&raw, "id = ?", id.Key())
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, UserNotFound(id.String())
		}
		return nil, core.UnexpectedError(result.Error.Error())
	}

	entity, err := raw.toEntity()
	if err != nil {
		return nil, err
	}

	return entity, nil
}

func (r *repository) GetByEmailAddress(
	ctx context.Context,
	emailAddress valueobject.Emailer,
) (Nutritionist, *core.BaseError) {
	var user nutritionists

	result := r.db.WithContext(ctx).First(&user, "email = ?", emailAddress.String())
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, UserNotFound(emailAddress.String())
		}
		return nil, core.UnexpectedError(result.Error.Error())
	}

	entity, err := user.toEntity()
	if err != nil {
		return nil, err
	}

	return entity, nil
}

func (r *repository) Create(
	ctx context.Context,
	entity Nutritionist,
) *core.BaseError {
	entityDB := fromEntity(entity)
	result := r.db.WithContext(ctx).Create(&entityDB)
	if result.Error != nil {
		return core.UnexpectedError(result.Error.Error())
	}
	return nil
}

func (r *repository) Update(
	ctx context.Context,
	nutritionist Nutritionist,
) *core.BaseError {
	data := fromEntity(nutritionist)
	result := r.db.WithContext(ctx).Save(&data)
	if result.Error != nil {
		return core.UnexpectedError(result.Error.Error())
	}
	return nil
}

func (r *repository) Delete(ctx context.Context, id valueobject.Identifier) *core.BaseError {
	result := r.db.WithContext(ctx).Delete(&Entity{}, id.Key())
	if result.Error != nil {
		return core.UnexpectedError(result.Error.Error())
	}
	if result.RowsAffected == 0 {
		return UserNotFound(id.String())
	}

	return nil
}
