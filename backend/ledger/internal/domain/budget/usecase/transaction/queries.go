package transaction

import (
	"context"
	"errors"
	"log/slog"
	"slices"
	"time"

	"cloud.google.com/go/civil"
	"github.com/google/uuid"
	"github.com/m11ano/budget_planner/backend/auth/pkg/auth"
	"github.com/m11ano/budget_planner/backend/ledger/internal/common/uctypes"
	"github.com/m11ano/budget_planner/backend/ledger/internal/domain/budget/entity"
	"github.com/m11ano/budget_planner/backend/ledger/internal/domain/budget/usecase"
	"github.com/m11ano/budget_planner/backend/ledger/internal/infra/loghandler"
	"github.com/m11ano/budget_planner/backend/ledger/pkg/pgclient"
	"github.com/samber/lo"

	appErrors "github.com/m11ano/budget_planner/backend/ledger/internal/app/errors"
)

func (uc *UsecaseImpl) FindOneByID(
	ctx context.Context,
	id uuid.UUID,
	queryParams *uctypes.QueryGetOneParams,
) (*usecase.TransactionDTO, error) {
	const op = "FindOneByID"

	item, err := uc.transactionRepo.FindOneByID(ctx, id, queryParams)
	if err != nil {
		return nil, appErrors.Chainf(err, "%s.%s", uc.pkg, op)
	}

	if auth.IsNeedToCheckRights(ctx) {
		authData := auth.GetAuthData(ctx)
		if authData == nil || authData.AccountID != item.AccountID {
			return nil, appErrors.ErrForbidden
		}
	}

	dto, err := uc.entitiesToDTO(ctx, []*entity.Transaction{item})
	if err != nil {
		return nil, appErrors.Chainf(err, "%s.%s", uc.pkg, op)
	}

	if len(dto) == 0 {
		return nil, appErrors.Chainf(appErrors.ErrInternal, "%s.%s", uc.pkg, op)
	}

	return dto[0], nil
}

func (uc *UsecaseImpl) FindList(
	ctx context.Context,
	listOptions *usecase.TransactionListOptions,
	queryParams *uctypes.QueryGetListParams,
) ([]*usecase.TransactionDTO, error) {
	const op = "FindList"

	items, err := uc.transactionRepo.FindList(ctx, listOptions, queryParams)
	if err != nil {
		return nil, appErrors.Chainf(err, "%s.%s", uc.pkg, op)
	}

	out, err := uc.entitiesToDTO(ctx, items)
	if err != nil {
		return nil, appErrors.Chainf(err, "%s.%s", uc.pkg, op)
	}

	return out, nil
}

func (uc *UsecaseImpl) FindPagedList(
	ctx context.Context,
	listOptions *usecase.TransactionListOptions,
	queryParams *uctypes.QueryGetListParams,
) ([]*usecase.TransactionDTO, uint64, error) {
	const op = "FindPagedList"

	items, total, err := uc.transactionRepo.FindPagedList(ctx, listOptions, queryParams)
	if err != nil {
		return nil, 0, appErrors.Chainf(err, "%s.%s", uc.pkg, op)
	}

	out, err := uc.entitiesToDTO(ctx, items)
	if err != nil {
		return nil, 0, appErrors.Chainf(err, "%s.%s", uc.pkg, op)
	}

	return out, total, nil
}

func (uc *UsecaseImpl) FindListInMap(
	ctx context.Context,
	listOptions *usecase.TransactionListOptions,
	queryParams *uctypes.QueryGetListParams,
) (map[uuid.UUID]*usecase.TransactionDTO, error) {
	const op = "FindListInMap"

	items, err := uc.transactionRepo.FindList(ctx, listOptions, queryParams)
	if err != nil {
		return nil, appErrors.Chainf(err, "%s.%s", uc.pkg, op)
	}

	out, err := uc.entitiesToDTO(ctx, items)
	if err != nil {
		return nil, appErrors.Chainf(err, "%s.%s", uc.pkg, op)
	}

	result := lo.SliceToMap(out, func(item *usecase.TransactionDTO) (uuid.UUID, *usecase.TransactionDTO) {
		return item.Transaction.ID, item
	})

	return result, nil
}

type CountReportItemsSFResult struct {
	Items    []*entity.ReportItem
	HitCache bool
}

var reportsCacheTTL = time.Second * 30

func (uc *UsecaseImpl) CountReportItems(
	ctx context.Context,
	queryFilter usecase.CountReportItemsQueryFilter,
) ([]*entity.ReportItem, bool, error) {
	const op = "CountReportItems"

	key := buildKeyForCountReportItems(queryFilter)

	result, err, _ := uc.sfGroup.Do(key, func() (any, error) {
		cacheItems, err := uc.transactionCacheRepo.GetReports(ctx, key)
		if err == nil {
			uc.logger.InfoContext(ctx, "CountReportItems cache hit", slog.Any("key", key))

			return CountReportItemsSFResult{
				Items:    cacheItems,
				HitCache: true,
			}, nil
		} else if !errors.Is(err, appErrors.ErrNotFound) {
			uc.logger.ErrorContext(loghandler.WithSource(ctx), "redis get err", slog.Any("error", err))
		}

		uc.logger.InfoContext(ctx, "CountReportItems cache miss", slog.Any("key", key))

		items, err := uc.countReportItems(ctx, queryFilter)
		if err != nil {
			return nil, err
		}

		err = uc.transactionCacheRepo.SaveReports(ctx, key, items, &reportsCacheTTL)
		if err != nil {
			uc.logger.ErrorContext(loghandler.WithSource(ctx), "redis save err", slog.Any("error", err))
		}

		return CountReportItemsSFResult{
			Items:    items,
			HitCache: false,
		}, nil
	})
	if err != nil {
		return nil, false, appErrors.Chainf(err, "%s.%s", uc.pkg, op)
	}

	sfResult, ok := result.(CountReportItemsSFResult)
	if !ok {
		return nil, false, appErrors.Chainf(appErrors.ErrInternal, "%s.%s", uc.pkg, op)
	}

	return sfResult.Items, sfResult.HitCache, nil
}

func (uc *UsecaseImpl) countReportItems(
	ctx context.Context,
	queryFilter usecase.CountReportItemsQueryFilter,
) ([]*entity.ReportItem, error) {
	const op = "CountReportItems"

	out := make([]*entity.ReportItem, 0)

	err := uc.dbMasterClient.DoWithIsoLvl(ctx, pgclient.RepeatableRead, func(ctx context.Context) error {
		var err error

		categories, err := uc.categoryRepo.FindList(ctx, &usecase.CategoryListOptions{
			Sort: []uctypes.SortOption[usecase.CategoryListOptionsSortField]{
				{
					Field:  usecase.CategoryListOptionsSortFieldID,
					IsDesc: false,
				},
			},
		}, nil)
		if err != nil {
			return err
		}

		if queryFilter.CategoryID != nil {
			categories = lo.Filter(categories, func(item *entity.Category, _ int) bool {
				return item.ID == *queryFilter.CategoryID
			})
		}

		txReportItems, err := uc.transactionRepo.CountReportItems(ctx, queryFilter)
		if err != nil {
			return err
		}

		slices.SortFunc(txReportItems, func(a, b *entity.AccountTransactionReportItem) int {
			return a.Period.Compare(b.Period)
		})

		var dateStart civil.Date
		if queryFilter.DateFrom != nil {
			dateStart = *queryFilter.DateFrom
		} else {
			if len(txReportItems) == 0 {
				return nil
			}

			dateStart = txReportItems[0].Period
		}

		var dateEnd civil.Date
		if queryFilter.DateTo != nil {
			dateEnd = *queryFilter.DateTo
		} else {
			if len(txReportItems) == 0 {
				return nil
			}

			dateEnd = txReportItems[len(txReportItems)-1].Period.AddMonths(1).AddDays(-1)
		}

		periodStart := dateStart
		periodStart.Day = 1

		periodEnd := dateEnd
		periodEnd.Day = 1

		budgets, err := uc.budgetRepo.FindList(ctx, &usecase.BudgetListOptions{
			FilterAccountID:  &queryFilter.AccountID,
			FilterPeriodFrom: &periodStart,
			FilterPeriodTo:   &periodEnd,
		}, nil)
		if err != nil {
			return err
		}

		for p := periodStart; p.Compare(periodEnd) <= 0; p = p.AddMonths(1) {
			item := &entity.ReportItem{
				AccountID: queryFilter.AccountID,
				Items:     make([]*entity.AccountTransactionReportItem, 0),
			}

			if p.Compare(periodStart) == 0 {
				item.DateFrom = dateStart
			} else {
				item.DateFrom = p
			}

			if p.Compare(periodEnd) == 0 {
				item.DateTo = dateEnd
			} else {
				item.DateTo = p.AddMonths(1).AddDays(-1)
			}

			for _, category := range categories {
				itemInReports, ok := lo.Find(txReportItems, func(item *entity.AccountTransactionReportItem) bool {
					return item.Period.Compare(p) == 0 && item.CategoryID == category.ID
				})

				if ok {
					item.Items = append(item.Items, itemInReports)
				} else {
					repItem := &entity.AccountTransactionReportItem{
						Period:     p,
						CategoryID: category.ID,
					}

					budgetForItem, ok := lo.Find(budgets, func(budget *entity.Budget) bool {
						return budget.CategoryID == category.ID && budget.Period.Compare(repItem.Period) == 0
					})
					if ok {
						repItem.BudgetID = &budgetForItem.ID
						repItem.BudgetAmount = &budgetForItem.Amount
					}

					item.Items = append(item.Items, repItem)
				}
			}

			out = append(out, item)
		}

		return nil
	})
	if err != nil {
		return nil, appErrors.Chainf(err, "%s.%s", uc.pkg, op)
	}

	return out, nil
}
