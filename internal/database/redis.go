package database

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
)

func InitRedis() (*redis.Client, error) {
	addr := os.Getenv("REDIS_ADDR")
	if addr == "" {
		addr = "localhost:6379"
	}

	fmt.Printf("addr: %v\n", addr)

	password := os.Getenv("REDIS_PASSWORD")
	db, err := strconv.Atoi(os.Getenv("REDIS_DB"))
	if err != nil {
		return nil, err
	}

	poolSize, err := strconv.Atoi(os.Getenv("REDIS_POOL_SIZE"))
	if err != nil {
		return nil, err
	}
	if poolSize == 0 {
		poolSize = 10
	}

	connTimeout, err := time.ParseDuration(os.Getenv("REDIS_CONN_TIMEOUT"))
	if err != nil {
		return nil, err
	}
	if connTimeout == 0 {
		connTimeout = 5 * time.Second
	}

	maxRetries, _ := strconv.Atoi(os.Getenv("REDIS_MAX_RETRIES"))

	redis := redis.NewClient(&redis.Options{
		Addr:         addr,
		Password:     password,
		DB:           db,
		PoolSize:     poolSize,
		DialTimeout:  connTimeout,
		MaxRetries:   maxRetries,
		MinIdleConns: 2,
	})

	ctx, cancel := context.WithTimeout(context.Background(), connTimeout)
	defer cancel()

	if _, err := redis.Ping(ctx).Result(); err != nil {
		return nil, fmt.Errorf("redis connection failed: %v", err)
	}

	return redis, nil
}
