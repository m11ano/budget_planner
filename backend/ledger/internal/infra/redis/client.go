package redis

import "github.com/redis/go-redis/v9"

type ClientConfig struct {
	Addr     string
	Password string
}

func NewRedisClient(cfg ClientConfig) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
	})
}
