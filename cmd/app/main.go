package main

import (
	"fmt"
	"net/http"

	"github.com/JooseMM/nutripia-backend-api/cmd/app/app"
	bodyMeasurementTypes "github.com/JooseMM/nutripia-backend-api/internal/bodyMeasurements/types"
	"github.com/JooseMM/nutripia-backend-api/internal/clients/types"
	nutritionistTypes "github.com/JooseMM/nutripia-backend-api/internal/nutritionist/types"
	sessionTypes "github.com/JooseMM/nutripia-backend-api/internal/security/session/types"
	"github.com/JooseMM/nutripia-backend-api/internal/storage"
	"github.com/JooseMM/nutripia-backend-api/pkg/core"
)

func main() {
	var db, err = storage.InitDB("host=localhost user=postgres password=Password123! dbname=nutripia_db port=5432 sslmode=disable")
	if err != nil {
		fmt.Print("Error")
		return
	}

	err = db.AutoMigrate(&clientTypes.Client{})
	err = db.AutoMigrate(&nutritionistTypes.Nutritionist{})
	err = db.AutoMigrate(&bodyMeasurementTypes.BodyMeasurement{})
	err = db.AutoMigrate(&sessionTypes.Session{})

	if err != nil {
		fmt.Print("Migration failed: ", err)
		return
	}

	mux := http.NewServeMux()
	protectedMux := core.RecoveryMiddleware(mux)
	app := app.NewApp(db)

	mux.HandleFunc("POST /nutritionist/register", app.AuthenticationHandler.RegisterNutritionist)
	mux.HandleFunc("POST /nutritionist/login", app.AuthenticationHandler.LoginNutritionist)

	mux.HandleFunc("GET /nutritionist/{id}/clients", app.NutritionistHandler.GetNutritionistById)
	mux.HandleFunc("GET /nutritionist/{id}", app.NutritionistHandler.GetNutritionistById)
	mux.HandleFunc("PUT /nutritionist/{id}", app.NutritionistHandler.UpdateNutritionistById)
	mux.HandleFunc("DELETE /nutritionist/{id}", app.NutritionistHandler.DeleteNutritionistById)

	mux.HandleFunc("POST /nutritionist/client", app.ClientHandler.CreateClient)
	mux.HandleFunc("PUT /nutritionist/client", app.ClientHandler.UpdateClientById)
	mux.HandleFunc("GET /client/{id}", app.ClientHandler.GetClientById)
	mux.HandleFunc("DELETE /nutritionist/client/{id}", app.ClientHandler.DeleteClientById)

	mux.HandleFunc("POST /body-measurements", app.MeasurementHandler.Create)
	mux.HandleFunc("GET /body-measurements/{id}", app.MeasurementHandler.GetById)
	mux.HandleFunc("DELETE /body-measurements/{id}", app.MeasurementHandler.DeleteById)

	fmt.Println("Server starting on :3000...")
	serveErr := http.ListenAndServe(":3000", protectedMux)
	if serveErr != nil {
		fmt.Printf("Error starting server: %s\n", serveErr)
	}
}
