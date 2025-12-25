package budget

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/m11ano/budget_planner/backend/auth/pkg/auth"
	appErrors "github.com/m11ano/budget_planner/backend/ledger/internal/app/errors"
	"github.com/m11ano/budget_planner/backend/ledger/internal/common/uctypes"
	"github.com/m11ano/budget_planner/backend/ledger/internal/domain/budget/entity"
	"github.com/m11ano/budget_planner/backend/ledger/internal/domain/budget/usecase"
	"github.com/m11ano/budget_planner/backend/ledger/internal/infra/loghandler"
	"github.com/m11ano/budget_planner/backend/ledger/pkg/pgclient"
	"github.com/samber/lo"
)

func (uc *UsecaseImpl) CreateBudgetByDTO(
	ctx context.Context,
	in usecase.CreateBudgetDataInput,
) (*usecase.BudgetDTO, error) {
	const op = "CreateBudgetByDTO"

	if auth.IsNeedToCheckRights(ctx) {
		authData := auth.GetAuthData(ctx)
		if authData == nil || authData.AccountID != in.AccountID {
			return nil, appErrors.ErrForbidden
		}
	}

	budget, err := entity.NewBudget(
		in.AccountID,
		in.Amount,
		in.Period,
		in.CategoryID,
	)
	if err != nil {
		return nil, appErrors.Chainf(err, "%s.%s", uc.pkg, op)
	}

	base := context.WithoutCancel(ctx)
	clrCtx, cancel := context.WithTimeout(base, time.Second*10)
	go func(ctx context.Context) {
		defer cancel()
		err := uc.budgetCacheRepo.ClearForPrefixes(
			ctx,
			buildKeyForFindList(&usecase.BudgetListOptions{
				FilterAccountID: &budget.AccountID,
			}, nil),
		)
		if err != nil {
			uc.logger.ErrorContext(loghandler.WithSource(ctx), "redis clear err", slog.Any("error", err))
		}
	}(clrCtx)

	err = uc.dbMasterClient.DoWithIsoLvl(ctx, pgclient.Serializable, func(ctx context.Context) error {
		_, err := uc.categoryRepo.FindOneByID(ctx, budget.CategoryID, nil)
		if err != nil {
			if errors.Is(err, appErrors.ErrNotFound) {
				return appErrors.Chainf(appErrors.ErrBadRequest.WithHints("category not found"), "%s.%s", uc.pkg, op)
			}
			return err
		}

		check, err := uc.budgetRepo.FindList(ctx, &usecase.BudgetListOptions{
			FilterAccountID:  &budget.AccountID,
			FilterPeriod:     &budget.Period,
			FilterCategoryID: &budget.CategoryID,
		}, &uctypes.QueryGetListParams{
			Limit: 1,
		})
		if err != nil {
			return err
		}

		if len(check) > 0 {
			return appErrors.Chainf(
				appErrors.ErrBadRequest.WithHints("budget for period and category already exists"),
				"%s.%s", uc.pkg, op)
		}

		err = uc.budgetRepo.Create(ctx, budget)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, appErrors.Chainf(err, "%s.%s", uc.pkg, op)
	}

	budgetDTO, err := uc.entitiesToDTO(ctx, []*entity.Budget{budget})
	if err != nil {
		return nil, appErrors.Chainf(err, "%s.%s", uc.pkg, op)
	}

	if len(budgetDTO) == 0 {
		uc.logger.ErrorContext(loghandler.WithSource(ctx), "unpredicted empty budget dto")
		return nil, appErrors.Chainf(appErrors.ErrInternal, "%s.%s", uc.pkg, op)
	}

	return budgetDTO[0], nil
}

func (uc *UsecaseImpl) PatchBudgetByDTO(
	ctx context.Context,
	id uuid.UUID,
	in usecase.PatchBudgetDataInput,
	skipVersionCheck bool,
) error {
	const op = "PatchBudgetByDTO"

	budget, err := uc.budgetRepo.FindOneByID(ctx, id, nil)
	if err != nil {
		return appErrors.Chainf(err, "%s.%s", uc.pkg, op)
	}

	base := context.WithoutCancel(ctx)
	clrCtx, cancel := context.WithTimeout(base, time.Second*10)
	go func(ctx context.Context) {
		defer cancel()
		err := uc.budgetCacheRepo.ClearForPrefixes(
			ctx,
			buildKeyForFindOneByID(id),
			buildKeyForFindList(&usecase.BudgetListOptions{
				FilterAccountID: &budget.AccountID,
			}, nil),
		)
		if err != nil {
			uc.logger.ErrorContext(loghandler.WithSource(ctx), "redis clear err", slog.Any("error", err))
		}
	}(clrCtx)

	err = uc.dbMasterClient.DoWithIsoLvl(ctx, pgclient.Serializable, func(ctx context.Context) error {
		budget, err := uc.budgetRepo.FindOneByID(ctx, id, &uctypes.QueryGetOneParams{
			ForUpdate: true,
		})
		if err != nil {
			return err
		}

		if auth.IsNeedToCheckRights(ctx) {
			authData := auth.GetAuthData(ctx)
			if authData == nil || authData.AccountID != budget.AccountID {
				return appErrors.ErrForbidden
			}
		}

		if !skipVersionCheck && budget.Version() != in.Version {
			return appErrors.ErrVersionConflict.
				WithDetail("last_version", false, budget.Version()).
				WithDetail("last_updated_at", false, budget.UpdatedAt)
		}

		if in.Amount != nil {
			err = budget.SetAmount(*in.Amount)
			if err != nil {
				return err
			}
		}

		checkFilter := &usecase.BudgetListOptions{
			FilterAccountID:  &budget.AccountID,
			FilterPeriod:     &budget.Period,
			FilterCategoryID: &budget.CategoryID,
		}
		needCheck := false

		if in.CategoryID != nil && *in.CategoryID != budget.CategoryID {
			budget.CategoryID = *in.CategoryID

			needCheck = true
			checkFilter.FilterCategoryID = &budget.CategoryID
		}

		if in.Period != nil &&
			(in.Period.Day != budget.Period.Day ||
				in.Period.Month != budget.Period.Month ||
				in.Period.Year != budget.Period.Year) {

			err := budget.SetPeriod(*in.Period)
			if err != nil {
				return err
			}

			needCheck = true
			checkFilter.FilterPeriod = &budget.Period
		}

		if needCheck {
			check, err := uc.budgetRepo.FindList(ctx, checkFilter, &uctypes.QueryGetListParams{
				Limit: 1,
			})
			if err != nil {
				return err
			}

			if len(check) > 0 {
				return appErrors.Chainf(
					appErrors.ErrBadRequest.WithHints("budget for period and category already exists"),
					"%s.%s", uc.pkg, op)
			}
		}

		err = uc.budgetRepo.Update(ctx, budget)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return appErrors.Chainf(err, "%s.%s", uc.pkg, op)
	}

	return nil
}

func (uc *UsecaseImpl) DeleteBudgetByID(ctx context.Context, id uuid.UUID) error {
	const op = "DeleteBudgetByID"

	budget, err := uc.budgetRepo.FindOneByID(ctx, id, nil)
	if err != nil {
		return appErrors.Chainf(err, "%s.%s", uc.pkg, op)
	}

	base := context.WithoutCancel(ctx)
	clrCtx, cancel := context.WithTimeout(base, time.Second*10)
	go func(ctx context.Context) {
		defer cancel()
		err := uc.budgetCacheRepo.ClearForPrefixes(
			ctx,
			buildKeyForFindOneByID(id),
			buildKeyForFindList(&usecase.BudgetListOptions{
				FilterAccountID: &budget.AccountID,
			}, nil),
		)
		if err != nil {
			uc.logger.ErrorContext(loghandler.WithSource(ctx), "redis clear err", slog.Any("error", err))
		}
	}(clrCtx)

	err = uc.dbMasterClient.Do(ctx, func(ctx context.Context) error {
		budget, err := uc.budgetRepo.FindOneByID(ctx, id, &uctypes.QueryGetOneParams{
			ForUpdate: true,
		})
		if err != nil {
			return err
		}

		if auth.IsNeedToCheckRights(ctx) {
			authData := auth.GetAuthData(ctx)
			if authData == nil || authData.AccountID != budget.AccountID {
				return appErrors.ErrForbidden
			}
		}

		budget.DeletedAt = lo.ToPtr(time.Now())

		err = uc.budgetRepo.Update(ctx, budget)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return appErrors.Chainf(err, "%s.%s", uc.pkg, op)
	}

	return nil
}
