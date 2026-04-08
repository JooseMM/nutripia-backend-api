package main

import (
	"fmt"
	"net/http"

	"github.com/JooseMM/nutripia-backend-api/cmd/app/app"
	"github.com/JooseMM/nutripia-backend-api/internal/storage"
	userTypes "github.com/JooseMM/nutripia-backend-api/internal/users/types"
	"github.com/JooseMM/nutripia-backend-api/pkg/core"
)

func main() {
	var db, err = storage.InitDB("host=localhost user=postgres password=Password123! dbname=nutripia_db port=5432 sslmode=disable")
	if err != nil {
		fmt.Print("Error")
		return
	}

	err = db.AutoMigrate(&userTypes.User{})
	if err != nil {
		fmt.Print("Migration failed: ", err)
		return
	}

	mux := http.NewServeMux()
	protectedMux := core.RecoveryMiddleware(mux)
	app := app.NewApp(db)

	mux.HandleFunc("POST /body-measurements", app.MeasurementHandler.Create)
	mux.HandleFunc("GET /body-measurements/{id}", app.MeasurementHandler.GetById)
	mux.HandleFunc("DELETE /users/{id}", app.MeasurementHandler.DeleteById)


	mux.HandleFunc("POST /users", app.UserHandler.Create)
	mux.HandleFunc("GET /users/{id}", app.UserHandler.GetById)
	mux.HandleFunc("DELETE /users/{id}", app.UserHandler.DeleteById)

	fmt.Println("Server starting on :3000...")
	serveErr := http.ListenAndServe(":3000", protectedMux)
	if serveErr != nil {
		fmt.Printf("Error starting server: %s\n", serveErr)
	}
}
