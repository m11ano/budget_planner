package providing

import (
	"github.com/m11ano/budget_planner/backend/ledger/internal/app/config"
	redisClient "github.com/m11ano/budget_planner/backend/ledger/internal/infra/redis"
	"github.com/redis/go-redis/v9"
)

func NewRedisClient(cfg config.Config) *redis.Client {
	return redisClient.NewRedisClient(redisClient.ClientConfig{
		Addr:         cfg.Redis.Addr,
		Password:     cfg.Redis.Password,
		DialTimeout:  cfg.Redis.DialTimeout,
		ReadTimeout:  cfg.Redis.ReadTimeout,
		WriteTimeout: cfg.Redis.WriteTimeout,
		PoolSize:     cfg.Redis.PoolSize,
		MinIdleConns: cfg.Redis.MinIdleConns,
		MaxRetries:   cfg.Redis.MaxRetries,
	})
}
