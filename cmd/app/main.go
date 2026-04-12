package main

import (
	"fmt"
	"net/http"

	"github.com/JooseMM/nutripia-backend-api/cmd/app/app"
	bodyMeasurementTypes "github.com/JooseMM/nutripia-backend-api/internal/bodyMeasurements/types"
	"github.com/JooseMM/nutripia-backend-api/internal/clients/types"
	nutritionistTypes "github.com/JooseMM/nutripia-backend-api/internal/nutritionist/types"
	sessionTypes "github.com/JooseMM/nutripia-backend-api/internal/security/session/types"
	verificationTypes "github.com/JooseMM/nutripia-backend-api/internal/security/verificationCode/types"
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
		fmt.Print("Error")
		return
	}

	if err = db.AutoMigrate(&clientTypes.Client{}); err != nil {
		fmt.Print("Migration failed: ", err)
		return
	}

	if err = db.AutoMigrate(&nutritionistTypes.Nutritionist{}); err != nil {
		fmt.Print("Migration failed: ", err)
		return
	}

	if err = db.AutoMigrate(&bodyMeasurementTypes.BodyMeasurement{}); err != nil {
		fmt.Print("Migration failed: ", err)
		return
	}

	if err = db.AutoMigrate(&sessionTypes.Session{}); err != nil {
		fmt.Print("Migration failed: ", err)
		return
	}

	if err = db.AutoMigrate(&verificationTypes.VerificationToken{}); err != nil {
		fmt.Print("Migration failed: ", err)
		return
	}

	mux := http.NewServeMux()
	protectedMux := core.RecoveryMiddleware(mux)
	app := app.NewApp(db)

	/* TODO: admin routes
	 	mux.HandleFunc("DELETE /nutritionist/{id}", app.NutritionistHandler.DeleteNutritionistById)
	*/

	mux.HandleFunc("GET /nutritionist/{id}", app.NutritionistHandler.GetNutritionistById)
	mux.HandleFunc("POST /nutritionist/register", app.AuthenticationHandler.RegisterNutritionist)
	mux.HandleFunc("POST /nutritionist/login", app.AuthenticationHandler.LoginNutritionist)
	mux.HandleFunc("POST /nutritionist/send-reset-password", app.AuthenticationHandler.SendResetPasswordToken)

	mux.HandleFunc(
		"POST /nutritionist/confirm-email",
		app.AuthenticationHandler.ConfirmedNutritionistEmail,
	)

	mux.Handle("PUT /nutritionist/{id}",
		app.AuthorizationMiddlewares.NutritionistOnly(
			app.NutritionistHandler.UpdateNutritionistById,
		),
	)

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
	serveErr := http.ListenAndServe(":3000", protectedMux)
	if serveErr != nil {
		fmt.Printf("Error starting server: %s\n", serveErr)
	}
}
