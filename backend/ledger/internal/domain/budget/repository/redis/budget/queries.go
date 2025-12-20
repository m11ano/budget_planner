package budget

import (
	"context"
	"encoding/json"

	appErrors "github.com/m11ano/budget_planner/backend/ledger/internal/app/errors"
	"github.com/m11ano/budget_planner/backend/ledger/internal/domain/budget/entity"
	redisRepo "github.com/m11ano/budget_planner/backend/ledger/internal/domain/budget/repository/redis"
	"github.com/redis/go-redis/v9"
)

func (r *Repository) GetBudgetsList(ctx context.Context, key string) ([]*entity.Budget, error) {
	const op = "GetBudgetsList"

	payload, err := r.redisClient.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, appErrors.Chainf(appErrors.ErrNotFound, "%s.%s", r.pkg, op)
		}
		return nil, appErrors.Chainf(appErrors.ErrInternal.WithWrap(err), "%s.%s", r.pkg, op)
	}

	dbItems := make([]*redisRepo.BudgetModel, 0)

	err = json.Unmarshal([]byte(payload), &dbItems)
	if err != nil {
		return nil, appErrors.Chainf(appErrors.ErrInternal.WithWrap(err), "%s.%s", r.pkg, op)
	}

	items := make([]*entity.Budget, 0, len(dbItems))
	for _, it := range dbItems {
		items = append(items, it.ToEntity())
	}

	return items, nil
}

func (r *Repository) GetBudgetsPagedList(ctx context.Context, key string) ([]*entity.Budget, uint64, error) {
	const op = "GetBudgetsPagedList"

	payload, err := r.redisClient.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, 0, appErrors.Chainf(appErrors.ErrNotFound, "%s.%s", r.pkg, op)
		}
		return nil, 0, appErrors.Chainf(appErrors.ErrInternal.WithWrap(err), "%s.%s", r.pkg, op)
	}

	dbPayload := &redisRepo.BudgetPagedList{
		Items: make([]*redisRepo.BudgetModel, 0),
	}

	err = json.Unmarshal([]byte(payload), &dbPayload)
	if err != nil {
		return nil, 0, appErrors.Chainf(appErrors.ErrInternal.WithWrap(err), "%s.%s", r.pkg, op)
	}

	items := make([]*entity.Budget, 0, len(dbPayload.Items))
	for _, it := range dbPayload.Items {
		items = append(items, it.ToEntity())
	}

	return items, dbPayload.Total, nil
}

func (r *Repository) GetBudget(ctx context.Context, key string) (*entity.Budget, error) {
	const op = "GetBudget"

	payload, err := r.redisClient.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, appErrors.Chainf(appErrors.ErrNotFound, "%s.%s", r.pkg, op)
		}
		return nil, appErrors.Chainf(appErrors.ErrInternal.WithWrap(err), "%s.%s", r.pkg, op)
	}

	dbItem := &redisRepo.BudgetModel{}

	err = json.Unmarshal([]byte(payload), dbItem)
	if err != nil {
		return nil, appErrors.Chainf(appErrors.ErrInternal.WithWrap(err), "%s.%s", r.pkg, op)
	}

	return dbItem.ToEntity(), nil
}
