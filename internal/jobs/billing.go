package jobs

import (
	"context"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/robfig/cron/v3"
	"github.com/zaher1307/subscription-service/internal/services"
)

func StartBillingJob(billingService *services.BillingService, redis *redis.Client) {
	c := cron.New(cron.WithLocation(time.UTC))

    _, err := c.AddFunc("0 0 * * *", func() {
		log.Println("Attempting to run billing job...")
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		lockKey := "billing_job_lock"
		lockTTL := 10 * time.Minute

		maxRetries := 3
		retryDelay := 1 * time.Second
		var acquired bool
		var err error

		for i := range maxRetries {
			acquired, err = redis.SetNX(ctx, lockKey, "locked", lockTTL).Result()
			if err == nil && acquired {
				break
			}
			log.Printf("Lock acquisition attempt %d failed: %v", i+1, err)
			time.Sleep(retryDelay)
			retryDelay *= 2
		}

		if !acquired {
			log.Println("Failed to acquire lock after retries")
			return
		}

		defer func() {
			if _, err := redis.Del(ctx, lockKey).Result(); err != nil {
				log.Printf("Failed to release lock: %v", err)
			}
		}()

		log.Println("Lock acquired. Running billing job...")
		if err := billingService.GenerateBills(); err != nil {
			log.Printf("Error generating bills: %v", err)
		}
	})
	if err != nil {
		log.Fatalf("Failed to schedule billing job: %v", err)
	}

	c.Start()
}
