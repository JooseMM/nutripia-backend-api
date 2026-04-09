package authentication

import "gorm.io/gorm"

type postgresRepository struct {
	db *gorm.DB
}
