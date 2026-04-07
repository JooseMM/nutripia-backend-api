package app

import (
	"github.com/JooseMM/nutripia-backend-api/internal/users"
	"gorm.io/gorm"
)

type App struct {
	UserHandler users.IUserHandler
}

func NewApp(db *gorm.DB) *App {
	userRepo := users.NewUserRepository(db)
	userService := users.NewUserService(userRepo)
	userHandler := users.NewUserHandler(userService)

	return &App{
		UserHandler: userHandler,
	}
}
