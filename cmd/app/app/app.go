package app

import (
	bodyMeasurement "github.com/JooseMM/nutripia-backend-api/internal/bodyMeasurements"
	"github.com/JooseMM/nutripia-backend-api/internal/users"
	"gorm.io/gorm"
)

type App struct {
	UserHandler        users.IUserHandler
	MeasurementHandler bodyMeasurement.IBodyMeasurementHandler
}

func NewApp(db *gorm.DB) *App {
	userRepo := users.NewUserRepository(db)
	userService := users.NewUserService(userRepo)
	userHandler := users.NewMeasurementHandler(userService)

	measurementRepo := bodyMeasurement.NewBodyMeasurementRepository(db)
	measurementService := bodyMeasurement.NewBodyMeasurementService(measurementRepo)
	measurementHandler := bodyMeasurement.NewBodyMeasurementHandler(measurementService, userService)

	return &App{
		UserHandler:        userHandler,
		MeasurementHandler: measurementHandler,
	}
}
