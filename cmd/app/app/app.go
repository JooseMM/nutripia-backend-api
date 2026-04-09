package app

import (
	bodyMeasurement "github.com/JooseMM/nutripia-backend-api/internal/bodyMeasurements"
	"github.com/JooseMM/nutripia-backend-api/internal/clients"
	"github.com/JooseMM/nutripia-backend-api/internal/nutritionist"
	"gorm.io/gorm"
)

type App struct {
	NutritionistHandler nutritionist.INutritionistHandler
	ClientHandler       clients.IClientHandler
	MeasurementHandler  bodyMeasurement.IBodyMeasurementHandler
}

func NewApp(db *gorm.DB) *App {
	clientRepo := clients.NewClientRepository(db)
	clientService := clients.NewClientService(clientRepo)

	nutritionistRepo := nutritionist.NewNutritionistRepository(db)
	nutritionistService := nutritionist.NewUserService(nutritionistRepo)

	measurementRepo := bodyMeasurement.NewBodyMeasurementRepository(db)
	measurementService := bodyMeasurement.NewBodyMeasurementService(measurementRepo)
	measurementHandler := bodyMeasurement.NewBodyMeasurementHandler(
		measurementService,
		clientService,
	)

	return &App{
		ClientHandler:       clients.NewClientHandler(clientService),
		NutritionistHandler: nutritionist.NewClientHandler(nutritionistService),
		MeasurementHandler:  measurementHandler,
	}
}
