package handlers

import (
	"context"
	"database/sql"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

type HealthHandler struct {
	DB    *sql.DB
	Redis *redis.Client
}

func NewHealthHandler(db *sql.DB, redis *redis.Client) *HealthHandler {
	return &HealthHandler{DB: db, Redis: redis}
}

func (h *HealthHandler) Check(c *gin.Context) {
	status := "ok"
	dbStatus := "ok"
	redisStatus := "ok"

	if err := h.DB.Ping(); err != nil {
		dbStatus = "error"
		status = "degraded"
	}

	connTimeout, _ := time.ParseDuration(os.Getenv("REDIS_CONN_TIMEOUT"))
	if connTimeout == 0 {
		connTimeout = 5 * time.Second
	}
	ctx, cancel := context.WithTimeout(context.Background(), connTimeout)
	defer cancel()

	if _, err := h.Redis.Ping(ctx).Result(); err != nil {
		redisStatus = "error"
		status = "degraded"
	}

	c.JSON(http.StatusOK, gin.H{
		"status":    status,
		"timestamp": time.Now().Format(time.RFC3339),
		"components": gin.H{
			"database": dbStatus,
			"redis":    redisStatus,
		},
	})
}
