package main

import (
	"fmt"
	"log"
	"net/http"

	"leaderboard-service/db"
	"leaderboard-service/db/migrations"
	_ "leaderboard-service/docs" // Import generated Swagger docs
	"leaderboard-service/models"
	"leaderboard-service/routes"

	"github.com/joho/godotenv"
)

// @title Leaderboard Service API
// @version 1.0
// @description API for managing leaderboards, entries, participants, and metrics
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.url http://www.example.com/support
// @contact.email support@example.com
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:8080
// @BasePath /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and the JWT token.
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
	err = db.DB.AutoMigrate(
		&models.Leaderboard{},
		&models.LeaderboardMetric{},
		&models.LeaderboardEntry{},
		&models.Participant{},
		&models.Metric{},
		&models.MetricValue{},
	)
	if err != nil {
		log.Fatal("Error migrating database: ", err)
	}

	r := router.Router()

	fmt.Println("Server is running on port 8080")
	fmt.Println("Swagger UI is available at http://localhost:8080/swagger/index.html")
	log.Fatal(http.ListenAndServe("localhost:8080", r))
}
