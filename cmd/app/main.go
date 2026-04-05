package main

import (
	"fmt"
	"net/http"

	"github.com/JooseMM/nutripia-backend-api/internal/storage"
	"github.com/JooseMM/nutripia-backend-api/internal/users/handlers"
	"github.com/JooseMM/nutripia-backend-api/internal/users/models"
	"github.com/JooseMM/nutripia-backend-api/internal/users/repository"
)

func main() {
	var db, err = storage.InitDB("host=localhost user=postgres password=Password123! dbname=nutripia_db port=5432 sslmode=disable")
	if err != nil {
		fmt.Print("Error")
	}

	err = db.AutoMigrate(&models.User{})
	if err != nil {
		fmt.Print("Migration failed: ", err)
	}

	mux := http.NewServeMux()

	userRepo := repository.NewUserRepository(db)
	userHandler := handlers.NewUserHandler(userRepo)
	mux.HandleFunc("POST /users", userHandler.Create)

	fmt.Println("Server starting on :3000...")
	serveErr := http.ListenAndServe(":3000", mux)
	if serveErr != nil {
		fmt.Printf("Error starting server: %s\n", serveErr)
	}
}
