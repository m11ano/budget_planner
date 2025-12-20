package budget

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/m11ano/budget_planner/backend/auth/pkg/auth"
	"github.com/m11ano/budget_planner/backend/ledger/internal/common/uctypes"
	"github.com/m11ano/budget_planner/backend/ledger/internal/domain/budget/entity"
	"github.com/m11ano/budget_planner/backend/ledger/internal/domain/budget/usecase"
	"github.com/m11ano/budget_planner/backend/ledger/internal/infra/loghandler"

	appErrors "github.com/m11ano/budget_planner/backend/ledger/internal/app/errors"
)

type FindOneByIDSFResult struct {
	Item     *entity.Budget
	HitCache bool
}

var budgetCacheTTL = time.Second * 30

func (uc *UsecaseImpl) FindOneByID(
	ctx context.Context,
	id uuid.UUID,
	queryParams *uctypes.QueryGetOneParams,
) (*usecase.BudgetDTO, bool, error) {
	const op = "FindOneByID"

	key := buildKeyForFindOneByID(id)

	result, err, _ := uc.sfGroup.Do(key, func() (any, error) {
		if queryParams == nil || !queryParams.SkipCache {
			cacheItem, err := uc.budgetCacheRepo.GetBudget(ctx, key)
			if err == nil {
				uc.logger.InfoContext(ctx, "Budget::FindOneByID cache hit", slog.Any("key", key))

				return FindOneByIDSFResult{
					Item:     cacheItem,
					HitCache: true,
				}, nil
			} else if !errors.Is(err, appErrors.ErrNotFound) {
				uc.logger.ErrorContext(loghandler.WithSource(ctx), "redis get err", slog.Any("error", err))
			}

			uc.logger.InfoContext(ctx, "Budget::FindOneByID cache miss", slog.Any("key", key))
		}

		item, err := uc.budgetRepo.FindOneByID(ctx, id, queryParams)
		if err != nil {
			return nil, appErrors.Chainf(err, "%s.%s", uc.pkg, op)
		}

		err = uc.budgetCacheRepo.SaveBudget(ctx, key, item, &budgetCacheTTL)
		if err != nil {
			uc.logger.ErrorContext(loghandler.WithSource(ctx), "redis save err", slog.Any("error", err))
		}

		return FindOneByIDSFResult{
			Item:     item,
			HitCache: false,
		}, nil
	})
	if err != nil {
		return nil, false, appErrors.Chainf(err, "%s.%s", uc.pkg, op)
	}

	sfResult, ok := result.(FindOneByIDSFResult)
	if !ok {
		return nil, false, appErrors.Chainf(appErrors.ErrInternal, "%s.%s", uc.pkg, op)
	}

	if auth.IsNeedToCheckRights(ctx) {
		authData := auth.GetAuthData(ctx)
		if authData == nil || authData.AccountID != sfResult.Item.AccountID {
			return nil, sfResult.HitCache, appErrors.ErrForbidden
		}
	}

	dto, err := uc.entitiesToDTO(ctx, []*entity.Budget{sfResult.Item})
	if err != nil {
		return nil, sfResult.HitCache, appErrors.Chainf(err, "%s.%s", uc.pkg, op)
	}

	if len(dto) == 0 {
		return nil, sfResult.HitCache, appErrors.Chainf(appErrors.ErrInternal, "%s.%s", uc.pkg, op)
	}

	return dto[0], sfResult.HitCache, nil
}

type FindListSFResult struct {
	Items    []*entity.Budget
	HitCache bool
}

func (uc *UsecaseImpl) FindList(
	ctx context.Context,
	listOptions *usecase.BudgetListOptions,
	queryParams *uctypes.QueryGetListParams,
) ([]*usecase.BudgetDTO, bool, error) {
	const op = "FindList"

	key := buildKeyForFindList(listOptions, queryParams)

	result, err, _ := uc.sfGroup.Do(key, func() (any, error) {
		if queryParams == nil || !queryParams.SkipCache {
			cacheItems, err := uc.budgetCacheRepo.GetBudgetsList(ctx, key)
			if err == nil {
				uc.logger.InfoContext(ctx, "Budget::GetBudgetsList cache hit", slog.Any("key", key))

				return FindListSFResult{
					Items:    cacheItems,
					HitCache: true,
				}, nil
			} else if !errors.Is(err, appErrors.ErrNotFound) {
				uc.logger.ErrorContext(loghandler.WithSource(ctx), "redis get err", slog.Any("error", err))
			}

			uc.logger.InfoContext(ctx, "Budget::GetBudgetsList cache miss", slog.Any("key", key))
		}

		items, err := uc.budgetRepo.FindList(ctx, listOptions, queryParams)
		if err != nil {
			return nil, appErrors.Chainf(err, "%s.%s", uc.pkg, op)
		}

		err = uc.budgetCacheRepo.SaveBudgetsList(ctx, key, items, &budgetCacheTTL)
		if err != nil {
			uc.logger.ErrorContext(loghandler.WithSource(ctx), "redis save err", slog.Any("error", err))
		}

		return FindListSFResult{
			Items:    items,
			HitCache: false,
		}, nil
	})
	if err != nil {
		return nil, false, appErrors.Chainf(err, "%s.%s", uc.pkg, op)
	}

	sfResult, ok := result.(FindListSFResult)
	if !ok {
		return nil, false, appErrors.Chainf(appErrors.ErrInternal, "%s.%s", uc.pkg, op)
	}

	out, err := uc.entitiesToDTO(ctx, sfResult.Items)
	if err != nil {
		return nil, sfResult.HitCache, appErrors.Chainf(err, "%s.%s", uc.pkg, op)
	}

	return out, sfResult.HitCache, nil
}

type FindPagedListSFResult struct {
	Items    []*entity.Budget
	Total    uint64
	HitCache bool
}

func (uc *UsecaseImpl) FindPagedList(
	ctx context.Context,
	listOptions *usecase.BudgetListOptions,
	queryParams *uctypes.QueryGetListParams,
) ([]*usecase.BudgetDTO, uint64, bool, error) {
	const op = "FindPagedList"

	key := buildKeyForFindPagedList(listOptions, queryParams)

	result, err, _ := uc.sfGroup.Do(key, func() (any, error) {
		if queryParams == nil || !queryParams.SkipCache {
			cacheItems, total, err := uc.budgetCacheRepo.GetBudgetsPagedList(ctx, key)
			if err == nil {
				uc.logger.InfoContext(ctx, "Budget::GetBudgetsPagedList cache hit", slog.Any("key", key))

				return FindPagedListSFResult{
					Items:    cacheItems,
					Total:    total,
					HitCache: true,
				}, nil
			} else if !errors.Is(err, appErrors.ErrNotFound) {
				uc.logger.ErrorContext(loghandler.WithSource(ctx), "redis get err", slog.Any("error", err))
			}

			uc.logger.InfoContext(ctx, "Budget::GetBudgetsPagedList cache miss", slog.Any("key", key))
		}

		items, total, err := uc.budgetRepo.FindPagedList(ctx, listOptions, queryParams)
		if err != nil {
			return nil, appErrors.Chainf(err, "%s.%s", uc.pkg, op)
		}

		err = uc.budgetCacheRepo.SaveBudgetsPagedList(ctx, key, items, total, &budgetCacheTTL)
		if err != nil {
			uc.logger.ErrorContext(loghandler.WithSource(ctx), "redis save err", slog.Any("error", err))
		}

		return FindPagedListSFResult{
			Items:    items,
			Total:    total,
			HitCache: false,
		}, nil
	})
	if err != nil {
		return nil, 0, false, appErrors.Chainf(err, "%s.%s", uc.pkg, op)
	}

	sfResult, ok := result.(FindPagedListSFResult)
	if !ok {
		return nil, 0, false, appErrors.Chainf(appErrors.ErrInternal, "%s.%s", uc.pkg, op)
	}

	out, err := uc.entitiesToDTO(ctx, sfResult.Items)
	if err != nil {
		return nil, sfResult.Total, sfResult.HitCache, appErrors.Chainf(err, "%s.%s", uc.pkg, op)
	}

	return out, sfResult.Total, sfResult.HitCache, nil
}
