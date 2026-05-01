package app

import (
	bodyMeasurement "github.com/JooseMM/nutripia-backend-api/internal/bodyMeasurements"
	"github.com/JooseMM/nutripia-backend-api/internal/clients"
	"github.com/JooseMM/nutripia-backend-api/internal/nutritionist"
	"github.com/JooseMM/nutripia-backend-api/internal/security/authentication"
	"github.com/JooseMM/nutripia-backend-api/internal/security/session"
	verification "github.com/JooseMM/nutripia-backend-api/internal/security/verificationCode"
	"gorm.io/gorm"
)

type App struct {
	NutritionistHandler   nutritionist.Handler
	ClientHandler         clients.Handler
	MeasurementHandler    bodyMeasurement.Handler
	AuthenticationHandler authentication.Handler
}

func NewApp(db *gorm.DB) *App {
	clientRepo, err := clients.NewRepository(db)
	if err != nil {
		panic(err)
	}
	clientService := clients.NewService(clientRepo)

	nutritionistRepo, err := nutritionist.NewRepository(db)
	if err != nil {
		panic(err)
	}

	nutritionistService := nutritionist.NewService(nutritionistRepo)

	measurementRepo, err := bodyMeasurement.NewRepository(db)
	if err != nil {
		panic(err)
	}

	measurementService := bodyMeasurement.NewService(measurementRepo)
	measurementHandler := bodyMeasurement.NewHandler(
		measurementService,
		clientService,
	)

	sessionRepo, err := session.NewRepository(db)
	if err != nil {
		panic(err)
	}
	sessionService := session.NewService(sessionRepo)

	verificationRepo, err := verification.NewRepository(db)
	if err != nil {
		panic(err)
	}
	verificationService := verification.NewService(verificationRepo)

	authenticationService := authentication.NewService(
		nutritionistRepo,
		sessionService,
		verificationService,
	)
	authenticationHandler := authentication.NewAuthenticationHandler(authenticationService)
	return &App{
		ClientHandler:         clients.NewClientHandler(clientService),
		NutritionistHandler:   nutritionist.NewNutritionistHandler(nutritionistService),
		MeasurementHandler:    measurementHandler,
		AuthenticationHandler: authenticationHandler,
	}
}
