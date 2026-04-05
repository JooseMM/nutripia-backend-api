package value_objects

type AuthenticationInformation struct {
	PasswordHash     string `gorm:"not null"`
	Role             string
	IsEmailConfirmed bool
}
