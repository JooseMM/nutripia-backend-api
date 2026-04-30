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

type entityDB struct {
	id               uuid.UUID `gorm:"type:uuid;primaryKey"`
	firstname        string    `gorm:"column:firstname;not null"`
	lastname         string    `gorm:"column:lastname;not null"`
	isEmailConfirmed bool      `gorm:"column:is_email_confirmed;not null"`
	email            string    `gorm:"column:email;uniqueIndex;not null"`
	birthDate        time.Time `gorm:"column:birth_date"`
	password         string    `gorm:"column:password;not null"`
	rut              string    `gorm:"column:rut;uniqueIndex;not null"`
	createdAt        time.Time `gorm:"column:created_at"`
	updatedAt        time.Time `gorm:"column:updated_at"`
}

func fromEntity(entity Nutritionist) entityDB {
	return entityDB{
		id:        entity.Id().Key(),
		firstname: entity.Name().Firstname(),
		lastname:  entity.Name().Lastname(),
		email:     entity.EmailAddress().String(),
		birthDate: entity.BirthDate().Date(),
		password:  entity.Password().String(),
	}
}

func (e *entityDB) toEntity() (Nutritionist, *core.BaseError) {
	var errList []string

	name, err := valueobject.NewName(e.firstname, e.lastname)
	errList = append(errList, err...)

	emailAddress, err := valueobject.NewEmailAddress(e.email)
	errList = append(errList, err...)

	birthDate, err := valueobject.NewBirthDate(e.birthDate)
	errList = append(errList, err...)

	password, err := valueobject.PasswordFromDB(e.password)
	errList = append(errList, err...)

	rut, err := valueobject.NewRUT(e.rut)
	errList = append(errList, err...)

	if len(errList) > 0 {
		return nil, core.ValidationError(errList)
	}

	return NewEntity(name, emailAddress, e.isEmailConfirmed, birthDate, password, rut), nil
}

func (entityDB) TableName() string {
	return "nutritionists"
}

type Repository struct {
	db *gorm.DB
}

type RepositoryManager interface {
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

func NewNutritionistRepository(db *gorm.DB) RepositoryManager {
	return &Repository{db}
}

func (r *Repository) IsEmailTaken(
	ctx context.Context,
	emalAddress valueobject.Emailer,
) (bool, *core.BaseError) {
	var exists bool

	err := r.db.WithContext(ctx).
		Model(&entityDB{}).
		Select("count(*) > 0").
		Where("email = ?", emalAddress.String()).
		Find(&exists).
		Error

	if err != nil {
		return false, core.UnexpectedError(err.Error())
	}
	return exists, nil
}

func (r *Repository) GetAll(
	ctx context.Context,
) ([]Nutritionist, *core.BaseError) {
	var rawList []*entityDB
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

func (r *Repository) GetById(
	ctx context.Context,
	id valueobject.Identifier,
) (Nutritionist, *core.BaseError) {
	var user entityDB

	result := r.db.WithContext(ctx).First(&user, "id = ?", id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, UserNotFound(id.String())
		}
		return nil, core.UnexpectedError(result.Error.Error())
	}

	entity, err := user.toEntity()
	if err != nil {
		return nil, err
	}

	return entity, nil
}

func (r *Repository) GetByEmailAddress(
	ctx context.Context,
	emailAddress valueobject.Emailer,
) (Nutritionist, *core.BaseError) {
	var user entityDB

	result := r.db.WithContext(ctx).First(&user, "email_address = ?", emailAddress.String())
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

func (r *Repository) Create(
	ctx context.Context,
	entity Nutritionist,
) *core.BaseError {
	entityDB := fromEntity(entity)
	result := r.db.WithContext(ctx).Create(entityDB)
	if result.Error != nil {
		return core.UnexpectedError(result.Error.Error())
	}
	return nil
}

func (r *Repository) Update(
	ctx context.Context,
	entity Nutritionist,
) *core.BaseError {
	entityDB := fromEntity(entity)
	result := r.db.WithContext(ctx).Save(entityDB)
	if result.Error != nil {
		return core.UnexpectedError(result.Error.Error())
	}
	return nil
}

func (r *Repository) Delete(ctx context.Context, id valueobject.Identifier) *core.BaseError {
	result := r.db.WithContext(ctx).Delete(&Entity{}, id.Key())
	if result.Error != nil {
		return core.UnexpectedError(result.Error.Error())
	}
	if result.RowsAffected == 0 {
		return UserNotFound(id.String())
	}

	return nil
}
