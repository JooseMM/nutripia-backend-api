package nutritionistTypes

import (
	"github.com/google/uuid"
)

type Nutritionist struct {
	ID uuid.UUID `gorm:"type:uuid;primaryKey"`

	NutritionistIdentity
	AuthenticationInformation

	TrackingInformation
}
