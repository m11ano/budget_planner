package transaction

import (
	"context"
	"encoding/json"
	"time"

	appErrors "github.com/m11ano/budget_planner/backend/ledger/internal/app/errors"
	"github.com/m11ano/budget_planner/backend/ledger/internal/domain/budget/entity"
	redisRepo "github.com/m11ano/budget_planner/backend/ledger/internal/domain/budget/repository/redis"
)

func (r *Repository) SaveReports(ctx context.Context, key string, items []*entity.ReportItem, ttl time.Duration) error {
	const op = "SaveReports"

	dbItems := make([]*redisRepo.ReportItemModel, 0, len(items))
	for _, it := range items {
		dbItems = append(dbItems, redisRepo.MapTransactionEntityToDBModel(it))
	}

	payload, err := json.Marshal(dbItems)
	if err != nil {
		return appErrors.Chainf(err, "%s.%s", r.pkg, op)
	}

	cmd := r.redisClient.Set(ctx, key, payload, ttl)
	if cmd.Err() != nil {
		return appErrors.Chainf(appErrors.ErrInternal.WithWrap(cmd.Err()), "%s.%s", r.pkg, op)
	}

	return nil
}

func (r *Repository) ClearForPrefixes(ctx context.Context, prefixes ...string) error {
	const op = "ClearForPrefixes"

	if len(prefixes) == 0 {
		return nil
	}

	const scanCount int64 = 1000
	const delBatch = 1000

	for _, prefix := range prefixes {
		if prefix == "" {
			continue
		}

		pattern := prefix + "*"
		var cursor uint64

		keysBatch := make([]string, 0, delBatch)

		for {
			if err := ctx.Err(); err != nil {
				return appErrors.Chainf(appErrors.ErrInternal.WithWrap(err), "%s.%s prefix=%s", r.pkg, op, pattern)
			}

			keys, nextCursor, err := r.redisClient.Scan(ctx, cursor, pattern, scanCount).Result()
			if err != nil {
				return appErrors.Chainf(appErrors.ErrInternal.WithWrap(err), "%s.%s prefix=%s", r.pkg, op, pattern)
			}

			for _, k := range keys {
				keysBatch = append(keysBatch, k)
				if len(keysBatch) >= delBatch {
					if err := r.redisClient.Unlink(ctx, keysBatch...).Err(); err != nil {
						return appErrors.Chainf(appErrors.ErrInternal.WithWrap(err), "%s.%s prefix=%s", r.pkg, op, pattern)
					}
					keysBatch = keysBatch[:0]
				}
			}

			cursor = nextCursor
			if cursor == 0 {
				break
			}
		}

		if len(keysBatch) > 0 {
			if err := r.redisClient.Unlink(ctx, keysBatch...).Err(); err != nil {
				return appErrors.Chainf(appErrors.ErrInternal.WithWrap(err), "%s.%s prefix=%s", r.pkg, op, pattern)
			}
		}
	}

	return nil
}
