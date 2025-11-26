package infrastructure

import (
	"context"
	"fmt"
	"log"

	"go-auth-service/config"

	"github.com/redis/go-redis/v9"
)

func NewRedisClient(cfg config.Config) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.RedisHost, cfg.RedisPort),
		Password: cfg.RedisPassword,
		DB:       0,
	})

	_, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		log.Printf("Warning: Failed to connect to Redis: %v. Continuing without Redis.", err)
		// We don't fatal here because Redis is optional/bonus, but better to have it working.
		// If strict requirement, we would fatal.
	}

	return rdb
}
