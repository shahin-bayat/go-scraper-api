package store

import (
	"context"
	"fmt"
	"os"

	"github.com/redis/go-redis/v9"
)

var (
	redisUser     = os.Getenv("REDIS_USER")
	redisPassword = os.Getenv("REDIS_PASSWORD")
	redisHost     = os.Getenv("REDIS_HOST")
	redisPort     = os.Getenv("REDIS_PORT")
	redisDB       = os.Getenv("REDIS_DB")
)

func NewRedisStore() (*redis.Client, error) {

	connStr := fmt.Sprintf("redis://%s:%s@%s:%s/%s", redisUser, redisPassword, redisHost, redisPort, redisDB)

	if os.Getenv("REDIS_INTERNAL_URL") != "" {
		connStr = os.Getenv("REDIS_INTERNAL_URL")
	}
	fmt.Printf("connStr: %s\n", connStr)
	opt, err := redis.ParseURL(connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse redis url: %w", err)
	}

	client := redis.NewClient(opt)
	if err := client.Ping(context.Background()).Err(); err != nil {
		return nil, fmt.Errorf("failed to ping redis: %w", err)
	}
	return client, nil
}
