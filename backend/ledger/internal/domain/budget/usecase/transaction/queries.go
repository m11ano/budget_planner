package transaction

import (
	"context"

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
) (*usecase.TransactionDTO, error) {
	const op = "FindOneByID"

	item, err := uc.transactionRepo.FindOneByID(ctx, id, queryParams)
	if err != nil {
		return nil, appErrors.Chainf(err, "%s.%s", uc.pkg, op)
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
