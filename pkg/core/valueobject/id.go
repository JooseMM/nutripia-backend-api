package valueobject

import (
	"github.com/JooseMM/nutripia-backend-api/pkg/core"
	"github.com/google/uuid"
)

type Identifier interface {
	String() string
	Key() uuid.UUID
	IsEqual(target Identifier) bool
}

type Id struct {
	key uuid.UUID
}

func (i *Id) String() string {
	return i.key.String()
}

func (i *Id) Key() uuid.UUID {
	return i.key
}

func (i *Id) IsEqual(target Identifier) bool {
	return target.Key() == i.Key()
}

func IdentifierFromString(raw string) (Identifier, *core.BaseError) {
	uuid, err := uuid.Parse(raw)
	if err != nil {
		return nil, core.ValidationError([]string{"ID: The identifier provided is not a valid UUID format."})
	}
	return &Id{key: uuid}, nil
}

func NewIdentifier() Identifier {
	return &Id{uuid.New()}
}
