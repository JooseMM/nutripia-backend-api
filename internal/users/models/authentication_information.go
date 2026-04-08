package userModels

type AuthenticationInformation struct {
	PasswordHash     string `gorm:"not null"`
	Role             Role
	IsEmailConfirmed bool
}
