package budget

import (
	"context"

	"cloud.google.com/go/civil"
	"github.com/google/uuid"
	"github.com/m11ano/budget_planner/backend/ledger/internal/common/uctypes"
	"github.com/m11ano/budget_planner/backend/ledger/internal/domain/budget/entity"
	"github.com/m11ano/budget_planner/backend/ledger/internal/domain/budget/usecase"
	"github.com/samber/lo"

	appErrors "github.com/m11ano/budget_planner/backend/ledger/internal/app/errors"
)

func (uc *UsecaseImpl) FindOneByID(
	ctx context.Context,
	id uuid.UUID,
	queryParams *uctypes.QueryGetOneParams,
) (*usecase.BudgetDTO, error) {
	const op = "FindOneByID"

	item, err := uc.budgetRepo.FindOneByID(ctx, id, queryParams)
	if err != nil {
		return nil, appErrors.Chainf(err, "%s.%s", uc.pkg, op)
	}

	dto, err := uc.entitiesToDTO(ctx, []*entity.Budget{item})
	if err != nil {
		return nil, appErrors.Chainf(err, "%s.%s", uc.pkg, op)
	}

	if len(dto) == 0 {
		return nil, appErrors.Chainf(appErrors.ErrInternal, "%s.%s", uc.pkg, op)
	}

	return dto[0], nil
}

func (uc *UsecaseImpl) FindOneByParams(
	ctx context.Context,
	accountID uuid.UUID,
	period civil.Date,
	categoryID uint64,
) (*usecase.BudgetDTO, error) {
	const op = "FindOneByParams"

	period.Day = 1

	item, err := uc.budgetRepo.FindList(ctx, &usecase.BudgetListOptions{
		FilterAccountID:  &accountID,
		FilterPeriod:     &period,
		FilterCategoryID: &categoryID,
	}, &uctypes.QueryGetListParams{
		Limit: 1,
	})
	if err != nil {
		return nil, appErrors.Chainf(err, "%s.%s", uc.pkg, op)
	}

	if len(item) == 0 {
		return nil, appErrors.Chainf(appErrors.ErrNotFound, "%s.%s", uc.pkg, op)
	}

	dto, err := uc.entitiesToDTO(ctx, []*entity.Budget{item[0]})
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
	listOptions *usecase.BudgetListOptions,
	queryParams *uctypes.QueryGetListParams,
) ([]*usecase.BudgetDTO, error) {
	const op = "FindList"

	items, err := uc.budgetRepo.FindList(ctx, listOptions, queryParams)
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
	listOptions *usecase.BudgetListOptions,
	queryParams *uctypes.QueryGetListParams,
) ([]*usecase.BudgetDTO, uint64, error) {
	const op = "FindPagedList"

	items, total, err := uc.budgetRepo.FindPagedList(ctx, listOptions, queryParams)
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
	listOptions *usecase.BudgetListOptions,
	queryParams *uctypes.QueryGetListParams,
) (map[uuid.UUID]*usecase.BudgetDTO, error) {
	const op = "FindListInMap"

	items, err := uc.budgetRepo.FindList(ctx, listOptions, queryParams)
	if err != nil {
		return nil, appErrors.Chainf(err, "%s.%s", uc.pkg, op)
	}

	out, err := uc.entitiesToDTO(ctx, items)
	if err != nil {
		return nil, appErrors.Chainf(err, "%s.%s", uc.pkg, op)
	}

	result := lo.SliceToMap(out, func(item *usecase.BudgetDTO) (uuid.UUID, *usecase.BudgetDTO) {
		return item.Budget.ID, item
	})

	return result, nil
}
