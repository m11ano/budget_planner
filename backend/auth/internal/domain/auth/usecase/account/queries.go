package account

import (
	"context"

	"github.com/google/uuid"
	"github.com/m11ano/budget_planner/backend/auth/internal/common/uctypes"
	"github.com/m11ano/budget_planner/backend/auth/internal/domain/auth/entity"
	"github.com/m11ano/budget_planner/backend/auth/internal/domain/auth/usecase"
	"github.com/samber/lo"

	appErrors "github.com/m11ano/budget_planner/backend/auth/internal/app/errors"
)

func (uc *UsecaseImpl) FindOneByEmail(
	ctx context.Context,
	email string,
	queryParams *uctypes.QueryGetOneParams,
) (*usecase.AccountDTO, error) {
	const op = "FindOneByEmail"

	item, err := uc.accountRepo.FindOneByEmail(ctx, email, queryParams)
	if err != nil {
		return nil, appErrors.Chainf(err, "%s.%s", uc.pkg, op)
	}

	dto, err := uc.entitiesToDTO(ctx, []*entity.Account{item})
	if err != nil {
		return nil, appErrors.Chainf(err, "%s.%s", uc.pkg, op)
	}

	if len(dto) == 0 {
		return nil, appErrors.Chainf(appErrors.ErrInternal, "%s.%s", uc.pkg, op)
	}

	return dto[0], nil
}

func (uc *UsecaseImpl) FindOneByID(
	ctx context.Context,
	id uuid.UUID,
	queryParams *uctypes.QueryGetOneParams,
) (*usecase.AccountDTO, error) {
	const op = "FindOneByID"

	item, err := uc.accountRepo.FindOneByID(ctx, id, queryParams)
	if err != nil {
		return nil, appErrors.Chainf(err, "%s.%s", uc.pkg, op)
	}

	dto, err := uc.entitiesToDTO(ctx, []*entity.Account{item})
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
	listOptions *usecase.AccountListOptions,
	queryParams *uctypes.QueryGetListParams,
) ([]*usecase.AccountDTO, error) {
	const op = "FindList"

	items, err := uc.accountRepo.FindList(ctx, listOptions, queryParams)
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
	listOptions *usecase.AccountListOptions,
	queryParams *uctypes.QueryGetListParams,
) ([]*usecase.AccountDTO, uint64, error) {
	const op = "FindPagedList"

	items, total, err := uc.accountRepo.FindPagedList(ctx, listOptions, queryParams)
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
	listOptions *usecase.AccountListOptions,
	queryParams *uctypes.QueryGetListParams,
) (map[uuid.UUID]*usecase.AccountDTO, error) {
	const op = "FindListInMap"

	items, err := uc.accountRepo.FindList(ctx, listOptions, queryParams)
	if err != nil {
		return nil, appErrors.Chainf(err, "%s.%s", uc.pkg, op)
	}

	out, err := uc.entitiesToDTO(ctx, items)
	if err != nil {
		return nil, appErrors.Chainf(err, "%s.%s", uc.pkg, op)
	}

	result := lo.SliceToMap(out, func(item *usecase.AccountDTO) (uuid.UUID, *usecase.AccountDTO) {
		return item.Account.ID, item
	})

	return result, nil
}
