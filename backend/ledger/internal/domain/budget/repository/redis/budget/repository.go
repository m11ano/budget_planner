package budget

import (
	"log/slog"

	"github.com/redis/go-redis/v9"
)

type Repository struct {
	pkg         string
	logger      *slog.Logger
	redisClient *redis.Client
}

func NewRepository(logger *slog.Logger, redisClient *redis.Client) *Repository {
	return &Repository{
		pkg:         "Budget.repository.RedisBudget",
		logger:      logger,
		redisClient: redisClient,
	}
}
