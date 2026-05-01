package main

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/JooseMM/nutripia-backend-api/cmd/app/app"
	"github.com/JooseMM/nutripia-backend-api/internal/storage"
	"github.com/JooseMM/nutripia-backend-api/pkg/core"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		fmt.Print(err.Error())
		return
	}

	var db, err = storage.InitDB("host=localhost user=postgres password=Password123! dbname=nutripia_db port=5432 sslmode=disable")
	if err != nil {
		slog.Error("Database initialization failed", "error", err)
		return
	}

	mux := http.NewServeMux()
	protectedMux := core.RecoveryMiddleware(mux)
	app := app.NewApp(db)

	/* Authentication */
	mux.HandleFunc(
		"POST /nutritionist/register",
		app.AuthenticationHandler.RegisterNutritionist,
	)

	mux.HandleFunc(
		"POST /nutritionist/login",
		app.AuthenticationHandler.LoginNutritionist,
	)
	mux.HandleFunc(
		"GET /nutritionist/{id}",
		app.NutritionistHandler.GetNutritionistById,
	)
	mux.HandleFunc(
		"POST /nutritionist/send-reset-password",
		app.AuthenticationHandler.SendResetPasswordToken,
	)
	mux.HandleFunc(
		"POST /nutritionist/verify-reset-password",
		app.AuthenticationHandler.VerifyResetPasswordToken,
	)
	mux.HandleFunc(
		"POST /nutritionist/complete-reset-password",
		app.AuthenticationHandler.CompleteResetPassword,
	)

	mux.HandleFunc(
		"POST /nutritionist/confirm-email",
		app.AuthenticationHandler.ConfirmNutritionistEmail,
	)

	/* Nutritionists */
	mux.HandleFunc("PUT /nutritionist/{id}",
		app.NutritionistHandler.UpdateNutritionistById,
	)

	/* Clients */
	mux.HandleFunc(
		"POST /nutritionist/client/{nutritionistId}",
		app.ClientHandler.CreateClient,
	)

	mux.HandleFunc("GET /clients/{nutritionistId}",
		app.ClientHandler.GetClientByNutritionist,
	)

	mux.HandleFunc("GET /client/{clientId}",
		app.ClientHandler.GetClientById,
	)

	mux.HandleFunc("DELETE /nutritionist/client/{clientId}",
		app.ClientHandler.DeleteClientById,
	)

	mux.HandleFunc("PUT /nutritionist/client/{clientId}",
		app.ClientHandler.UpdateClientById,
	)

	/* Measurements */
	mux.HandleFunc("POST /body-measurements/{clientId}",
		app.MeasurementHandler.Create,
	)
	mux.HandleFunc("GET /body-measurements/{id}",
		app.MeasurementHandler.GetById,
	)

	mux.HandleFunc("DELETE /body-measurements/{id}",
		app.MeasurementHandler.DeleteById,
	)

	fmt.Println("Server starting on port:3000")
	if err := http.ListenAndServe(":3000", protectedMux); err != nil {
		fmt.Printf("Error starting server: %s\n", err)
	}
}
