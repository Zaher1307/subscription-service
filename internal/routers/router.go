package router

import (
	"database/sql"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"

	"github.com/zaher1307/subscription-service/internal/handlers"
	"github.com/zaher1307/subscription-service/internal/repositories"
	"github.com/zaher1307/subscription-service/internal/services"
)

func SetupRouter(db *sql.DB, redis *redis.Client) *gin.Engine {
	r := gin.Default()

	userRepo := repositories.NewUserRepository(db)
	productRepo := repositories.NewProductRepository(db)
	subscriptionRepo := repositories.NewSubscriptionRepository(db)
	billRepo := repositories.NewBillRepository(db)

	subscriptionService := services.NewSubscriptionService(subscriptionRepo, productRepo, billRepo, userRepo)
	billingService := services.NewBillingService(subscriptionRepo, productRepo, billRepo, userRepo)
	userService := services.NewUserService(userRepo)
	productService := services.NewProductService(productRepo)

	userHandler := handlers.NewUserHandler(userService)
	productHandler := handlers.NewProductHandler(productService)
	subscriptionHandler := handlers.NewSubscriptionHandler(subscriptionService)
	billHandler := handlers.NewBillHandler(billingService)
	healthHandler := handlers.NewHealthHandler(db, redis)

	r.GET("/health", healthHandler.Check)

	api := r.Group("/api")
	{
		users := api.Group("/users")
		{
			users.POST("", userHandler.Create)
			users.GET("/:id", userHandler.GetByID)
		}

		products := api.Group("/products")
		{
			products.GET("", productHandler.GetAll)
			products.GET("/:id", productHandler.GetByID)
		}

		subscriptions := api.Group("/subscriptions")
		{
			subscriptions.POST("", subscriptionHandler.Create)
			subscriptions.GET("/:id", subscriptionHandler.GetByID)
		}

		bills := api.Group("/bills")
		{
			bills.GET("/user/:user_id", billHandler.GetUserBills)
			bills.GET("/:id", billHandler.GetBill)
			bills.POST("/:id/pay", billHandler.PayBill)
		}
	}

	return r
}
