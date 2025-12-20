package budget

import (
	"context"
	"encoding/json"
	"time"

	appErrors "github.com/m11ano/budget_planner/backend/ledger/internal/app/errors"
	"github.com/m11ano/budget_planner/backend/ledger/internal/domain/budget/entity"
	redisRepo "github.com/m11ano/budget_planner/backend/ledger/internal/domain/budget/repository/redis"
)

func (r *Repository) SaveBudgetsList(ctx context.Context, key string, items []*entity.Budget, ttl *time.Duration) error {
	const op = "SaveBudgetsList"

	dbItems := make([]*redisRepo.BudgetModel, 0, len(items))
	for _, it := range items {
		dbItems = append(dbItems, redisRepo.MapBudgetEntityToDBModel(it))
	}

	payload, err := json.Marshal(dbItems)
	if err != nil {
		return appErrors.Chainf(err, "%s.%s", r.pkg, op)
	}

	var payloadTTL time.Duration
	if ttl != nil {
		payloadTTL = *ttl
	}

	cmd := r.redisClient.Set(ctx, key, payload, payloadTTL)
	if cmd.Err() != nil {
		return appErrors.Chainf(appErrors.ErrInternal.WithWrap(cmd.Err()), "%s.%s", r.pkg, op)
	}

	return nil
}

func (r *Repository) SaveBudgetsPagedList(
	ctx context.Context,
	key string,
	items []*entity.Budget,
	total uint64,
	ttl *time.Duration,
) error {
	const op = "SaveBudgetsPagedList"

	dbPayload := &redisRepo.BudgetPagedList{
		Items: make([]*redisRepo.BudgetModel, 0, len(items)),
		Total: total,
	}
	for _, it := range items {
		dbPayload.Items = append(dbPayload.Items, redisRepo.MapBudgetEntityToDBModel(it))
	}

	payload, err := json.Marshal(dbPayload)
	if err != nil {
		return appErrors.Chainf(err, "%s.%s", r.pkg, op)
	}

	var payloadTTL time.Duration
	if ttl != nil {
		payloadTTL = *ttl
	}

	cmd := r.redisClient.Set(ctx, key, payload, payloadTTL)
	if cmd.Err() != nil {
		return appErrors.Chainf(appErrors.ErrInternal.WithWrap(cmd.Err()), "%s.%s", r.pkg, op)
	}

	return nil
}

func (r *Repository) SaveBudget(ctx context.Context, key string, item *entity.Budget, ttl *time.Duration) error {
	const op = "SaveBudget"

	dbItem := redisRepo.MapBudgetEntityToDBModel(item)

	payload, err := json.Marshal(dbItem)
	if err != nil {
		return appErrors.Chainf(err, "%s.%s", r.pkg, op)
	}

	var payloadTTL time.Duration
	if ttl != nil {
		payloadTTL = *ttl
	}

	cmd := r.redisClient.Set(ctx, key, payload, payloadTTL)
	if cmd.Err() != nil {
		return appErrors.Chainf(appErrors.ErrInternal.WithWrap(cmd.Err()), "%s.%s", r.pkg, op)
	}

	return nil
}
