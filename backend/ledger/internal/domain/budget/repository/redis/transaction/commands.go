package transaction

import (
	"context"
	"encoding/json"
	"time"

	appErrors "github.com/m11ano/budget_planner/backend/ledger/internal/app/errors"
	"github.com/m11ano/budget_planner/backend/ledger/internal/domain/budget/entity"
	redisRepo "github.com/m11ano/budget_planner/backend/ledger/internal/domain/budget/repository/redis"
)

func (r *Repository) SaveReports(ctx context.Context, key string, items []*entity.ReportItem, ttl *time.Duration) error {
	const op = "SaveReports"

	dbItems := make([]*redisRepo.ReportItemModel, 0, len(items))
	for _, it := range items {
		dbItems = append(dbItems, redisRepo.MapTransactionEntityToDBModel(it))
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
