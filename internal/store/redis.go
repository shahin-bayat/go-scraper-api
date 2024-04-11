package store

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
	"github.com/shahin-bayat/scraper-api/internal/config"
)

func NewRedisStore(cfg *config.RedisConfig) (*redis.Client, error) {
	connStr := fmt.Sprintf("redis://%s:%s@%s:%s/%s", cfg.RedisUser, cfg.RedisPassword, cfg.RedisHost, cfg.RedisPort, cfg.RedisDB)

	if cfg.RedisInternalUrl != "" {
		connStr = cfg.RedisInternalUrl
	}
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
