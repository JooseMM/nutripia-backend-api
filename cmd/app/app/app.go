package app

import (
	bodyMeasurement "github.com/JooseMM/nutripia-backend-api/internal/bodyMeasurements"
	"github.com/JooseMM/nutripia-backend-api/internal/clients"
	"github.com/JooseMM/nutripia-backend-api/internal/nutritionist"
	"github.com/JooseMM/nutripia-backend-api/internal/security/authentication"
	"github.com/JooseMM/nutripia-backend-api/internal/security/authorization"
	"github.com/JooseMM/nutripia-backend-api/internal/security/session"
	"gorm.io/gorm"
)

type App struct {
	NutritionistHandler      nutritionist.INutritionistHandler
	ClientHandler            clients.IClientHandler
	MeasurementHandler       bodyMeasurement.IBodyMeasurementHandler
	AuthenticationHandler    authentication.IAuthenticationHandler
	AuthorizationMiddlewares authorization.IAuthorizationMiddlewares
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

	sessionRepo := session.NewSessionRepository(db)
	sessionService := session.NewSessionService(sessionRepo)
	authenticationService := authentication.NewAuthenticationService(
		nutritionistRepo,
		sessionService,
	)
	authenticationHandler := authentication.NewAuthenticationHandler(authenticationService)

	authorizationMiddleware := authorization.NewAuthorizationMiddlewares(
		nutritionistRepo,
		clientRepo,
		sessionService,
	)
	return &App{
		ClientHandler:            clients.NewClientHandler(clientService),
		NutritionistHandler:      nutritionist.NewNutritionistHandler(nutritionistService),
		MeasurementHandler:       measurementHandler,
		AuthenticationHandler:    authenticationHandler,
		AuthorizationMiddlewares: authorizationMiddleware,
	}
}
