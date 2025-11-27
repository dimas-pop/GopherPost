package main

import (
	"gopher-post/db"
	"gopher-post/handlers"
	"gopher-post/routes"
	"log"
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
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file, Error:", err)
	}

	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL not found in .env file, Error:")
	}

	dbpool := db.InitDB(dbURL)
	defer dbpool.Close()

	srv := &handlers.Server{
		DB: dbpool,
	}

	r := routes.SetupRoutes(srv)

	log.Println("Menjalankan server di port 8080...")
	http.ListenAndServe(":8080", r)
}
