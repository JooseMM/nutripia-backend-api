package nutritionistTypes

type AuthenticationInformation struct {
	PasswordHash     string `gorm:"not null"`
	IsEmailConfirmed bool
}
