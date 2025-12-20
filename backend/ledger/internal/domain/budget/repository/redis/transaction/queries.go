package transaction

import (
	"context"
	"encoding/json"

	appErrors "github.com/m11ano/budget_planner/backend/ledger/internal/app/errors"
	"github.com/m11ano/budget_planner/backend/ledger/internal/domain/budget/entity"
	redisRepo "github.com/m11ano/budget_planner/backend/ledger/internal/domain/budget/repository/redis"
	"github.com/redis/go-redis/v9"
)

func (r *Repository) GetReports(ctx context.Context, key string) ([]*entity.ReportItem, error) {
	const op = "GetReports"

	payload, err := r.redisClient.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, appErrors.Chainf(appErrors.ErrNotFound, "%s.%s", r.pkg, op)
		}
		return nil, appErrors.Chainf(appErrors.ErrInternal.WithWrap(err), "%s.%s", r.pkg, op)
	}

	dbItems := make([]*redisRepo.ReportItemModel, 0)

	err = json.Unmarshal([]byte(payload), &dbItems)
	if err != nil {
		return nil, appErrors.Chainf(appErrors.ErrInternal.WithWrap(err), "%s.%s", r.pkg, op)
	}

	items := make([]*entity.ReportItem, 0, len(dbItems))
	for _, it := range dbItems {
		items = append(items, it.ToEntity())
	}

	return items, nil
}
