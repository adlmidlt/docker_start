package main

import (
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"start_golang_with_docker/internal/database"
	"start_golang_with_docker/internal/handlers"
	"start_golang_with_docker/internal/redis"
	"strings"
)

func main() {
	if err := godotenv.Load("/.env"); err != nil {
		log.Print("No .env file found")
	}

	database.Init()

	redis.Init()

	http.HandleFunc("/items/", func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/items/") && r.Method == http.MethodGet {
			handlers.GetItems(w, r)
		} else if strings.HasPrefix(r.URL.Path, "/items/") && r.Method == http.MethodPost {
			handlers.CreateItem(w, r)
		} else {
			handlers.GetItem(w, r)
		}
	})

	http.HandleFunc("/health", handlers.HealthCheck)

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Listening on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
