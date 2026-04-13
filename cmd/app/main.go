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

	/* TODO: admin routes
	mux.HandleFunc("DELETE /nutritionist/{id}", app.NutritionistHandler.DeleteNutritionistById)
	*/

	/* Authentication */
	mux.HandleFunc(
		"GET /nutritionist/{id}",
		app.NutritionistHandler.GetNutritionistById,
	)
	mux.HandleFunc(
		"POST /nutritionist/register",
		app.AuthenticationHandler.RegisterNutritionist,
	)
	mux.HandleFunc(
		"POST /nutritionist/login",
		app.AuthenticationHandler.LoginNutritionist,
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
		app.AuthenticationHandler.ConfirmedNutritionistEmail,
	)

	/* Nutritionists */
	mux.Handle("PUT /nutritionist/{id}",
		app.AuthorizationMiddlewares.NutritionistOnly(
			app.NutritionistHandler.UpdateNutritionistById,
		),
	)

	/* Clients */
	mux.Handle(
		"POST /nutritionist/client",
		app.AuthorizationMiddlewares.NutritionistOnly(
			app.ClientHandler.CreateClient,
		),
	)

	mux.Handle("GET /clients/",
		app.AuthorizationMiddlewares.NutritionistOnly(
			app.ClientHandler.GetClientByNutritionist,
		),
	)

	mux.Handle("GET /client/{id}",
		app.AuthorizationMiddlewares.NutritionistOnly(
			app.ClientHandler.GetClientById,
		),
	)

	mux.Handle("DELETE /nutritionist/client/{id}",
		app.AuthorizationMiddlewares.NutritionistOnly(
			app.ClientHandler.DeleteClientById,
		),
	)

	mux.Handle("PUT /nutritionist/client",
		app.AuthorizationMiddlewares.NutritionistOnly(
			app.ClientHandler.UpdateClientById,
		),
	)

	/* Measurements */
	mux.Handle("POST /body-measurements",
		app.AuthorizationMiddlewares.NutritionistOnly(
			app.MeasurementHandler.Create,
		),
	)
	mux.Handle("GET /body-measurements/{id}",
		app.AuthorizationMiddlewares.NutritionistOnly(
			app.MeasurementHandler.GetById,
		),
	)

	mux.Handle("DELETE /body-measurements/{id}",
		app.AuthorizationMiddlewares.NutritionistOnly(
			app.MeasurementHandler.DeleteById,
		),
	)

	fmt.Println("Server starting on port:3000")
	if err := http.ListenAndServe(":3000", protectedMux); err != nil {
		fmt.Printf("Error starting server: %s\n", err)
	}
}
