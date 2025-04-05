package main

import (
	"fmt"
	"log"
	"net/http"

	"leaderboard-service/db"
	"leaderboard-service/migrations"
	"leaderboard-service/models"
	"leaderboard-service/router"

	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	db.InitDB()

	// Run custom migrations first
	err = migrations.RegisterMigrations(db.DB)
	if err != nil {
		log.Fatal("Error running custom migrations: ", err)
	}

	// Then run auto-migration
	err = db.DB.AutoMigrate(&models.Leaderboard{}, &models.LeaderboardMetric{}, &models.LeaderboardEntry{}, &models.Participant{})
	if err != nil {
		log.Fatal("Error migrating database: ", err)
	}

	// http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	// 	fmt.Fprintf(w, "Hello World!")
	// })

	r := router.Router()

	fmt.Println("Server is running on port 8080")
	log.Fatal(http.ListenAndServe("localhost:8080", r))
}
