package usecase

import (
	"context"

	"github.com/m11ano/budget_planner/backend/ledger/internal/common/uctypes"
	"github.com/m11ano/budget_planner/backend/ledger/internal/domain/budget/entity"
)

type CategoryListOptions struct {
	FilterIDs *[]uint64
	Sort      []uctypes.SortOption[CategoryListOptionsSortField]
}

type CategoryListOptionsSortField string

const (
	CategoryListOptionsSortFieldID CategoryListOptionsSortField = "id"
)

type CategoryDTO struct {
	Category *entity.Category
}

type CategoryUsecase interface {
	FindOneByID(
		ctx context.Context,
		id uint64,
		queryParams *uctypes.QueryGetOneParams,
	) (resItem *CategoryDTO, resErr error)

	FindList(
		ctx context.Context,
		listOptions *CategoryListOptions,
		queryParams *uctypes.QueryGetListParams,
	) (resItems []*CategoryDTO, resErr error)
}

type CategoryRepository interface {
	FindOneByID(
		ctx context.Context,
		id uint64,
		queryParams *uctypes.QueryGetOneParams,
	) (transaction *entity.Category, err error)

	FindList(
		ctx context.Context,
		listOptions *CategoryListOptions,
		queryParams *uctypes.QueryGetListParams,
	) (items []*entity.Category, err error)
}
