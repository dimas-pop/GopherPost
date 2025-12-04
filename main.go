package main

import (
	"gopher-post/db"
	"gopher-post/handlers"
	"gopher-post/routes"
	"log/slog"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

// @title 			GopherPost API
// @version			1.0.0
// @description 	Ini adalah server API untuk GopherPost
// @termsOfService 	http://swagger.io/terms/

// @contact.name 	Dimas Saputra
// @contact.email 	dsaputra5403@gmail.com

// @host 			localhost:8080
// @BasePath 		/api

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	err := godotenv.Load()
	if err != nil {
		slog.Info("File .env not found, using system environment variables")
	}

	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		slog.Error("DB_URL not found in .env file")
		os.Exit(1)
	}

	dbpool := db.InitDB(dbURL)
	defer dbpool.Close()

	srv := &handlers.Server{
		DB: dbpool,
	}

	r := routes.SetupRoutes(srv)

	slog.Info("Server starting", "port", 8080)
	if err := http.ListenAndServe(":8080", r); err != nil {
		slog.Error("Server failed to start", "error", err)
		os.Exit(1)
	}
}
