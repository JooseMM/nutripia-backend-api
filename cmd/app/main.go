package main

import (
	"fmt"
	"net/http"

	"github.com/JooseMM/nutripia-backend-api/cmd/app/app"
	"github.com/JooseMM/nutripia-backend-api/internal/storage"
	userModels "github.com/JooseMM/nutripia-backend-api/internal/users/models"
)

func main() {
	var db, err = storage.InitDB("host=localhost user=postgres password=Password123! dbname=nutripia_db port=5432 sslmode=disable")
	if err != nil {
		fmt.Print("Error")
		return
	}

	err = db.AutoMigrate(&userModels.User{})
	if err != nil {
		fmt.Print("Migration failed: ", err)
		return
	}

	mux := http.NewServeMux()
	protectedMux := app.RecoveryMiddleware(mux)
	app := app.NewApp(db)

	mux.HandleFunc("POST /users", app.UserHandler.Create)
	mux.HandleFunc("GET /users/{id}", app.UserHandler.GetById)
	mux.HandleFunc("DELETE /users/{id}", app.UserHandler.DeleteById)

	fmt.Println("Server starting on :3000...")
	serveErr := http.ListenAndServe(":3000", protectedMux)
	if serveErr != nil {
		fmt.Printf("Error starting server: %s\n", serveErr)
	}
}
