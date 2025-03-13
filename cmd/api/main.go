package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"github.com/zaher1307/subscription-service/internal/database"
	"github.com/zaher1307/subscription-service/internal/jobs"
	"github.com/zaher1307/subscription-service/internal/repositories"
	router "github.com/zaher1307/subscription-service/internal/routers"
	"github.com/zaher1307/subscription-service/internal/services"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	ginMode := os.Getenv("GIN_MODE")
	if ginMode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}
	redis, err := database.InitRedis()
	if err != nil {
		log.Fatalf("Redis connection failed: %v", err)
	}

	db, err := database.Connect()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	r := router.SetupRouter(db, redis)

	billingService := services.NewBillingService(
		repositories.NewSubscriptionRepository(db),
		repositories.NewProductRepository(db),
		repositories.NewBillRepository(db),
		repositories.NewUserRepository(db),
	)

	jobs.StartBillingJob(billingService, redis)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server running on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
