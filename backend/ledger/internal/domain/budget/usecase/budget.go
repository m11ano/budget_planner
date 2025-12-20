package usecase

import (
	"context"
	"time"

	"cloud.google.com/go/civil"
	"github.com/google/uuid"
	"github.com/govalues/decimal"
	"github.com/m11ano/budget_planner/backend/ledger/internal/common/uctypes"
	"github.com/m11ano/budget_planner/backend/ledger/internal/domain/budget/entity"
)

type BudgetListOptions struct {
	FilterAccountID  *uuid.UUID
	FilterPeriod     *civil.Date
	FilterPeriodFrom *civil.Date
	FilterPeriodTo   *civil.Date
	FilterCategoryID *uint64
	Sort             []uctypes.SortOption[BudgetListOptionsSortField]
}

type BudgetListOptionsSortField string

const (
	BudgetListOptionsSortFieldPeriod    BudgetListOptionsSortField = "period"
	BudgetListOptionsSortFieldCreatedAt BudgetListOptionsSortField = "created_at"
)

type BudgetDTO struct {
	Budget *entity.Budget
}

type CreateBudgetDataInput struct {
	AccountID  uuid.UUID
	Period     civil.Date
	CategoryID uint64
	Amount     decimal.Decimal
}

type PatchBudgetDataInput struct {
	Version int64

	Amount     *decimal.Decimal
	Period     *civil.Date
	CategoryID *uint64
}

type BudgetUsecase interface {
	FindOneByID(
		ctx context.Context,
		id uuid.UUID,
		queryParams *uctypes.QueryGetOneParams,
	) (resItem *BudgetDTO, cacheHit bool, resErr error)

	FindList(
		ctx context.Context,
		listOptions *BudgetListOptions,
		queryParams *uctypes.QueryGetListParams,
	) (resItems []*BudgetDTO, cacheHit bool, resErr error)

	FindPagedList(
		ctx context.Context,
		listOptions *BudgetListOptions,
		queryParams *uctypes.QueryGetListParams,
	) (resItems []*BudgetDTO, total uint64, cacheHit bool, resErr error)

	CreateBudgetByDTO(
		ctx context.Context,
		in CreateBudgetDataInput,
	) (resBudgetDTO *BudgetDTO, resErr error)

	PatchBudgetByDTO(
		ctx context.Context,
		id uuid.UUID,
		in PatchBudgetDataInput,
		skipVersionCheck bool,
	) (resErr error)

	DeleteBudgetByID(
		ctx context.Context,
		id uuid.UUID,
	) (resErr error)
}

type BudgetRepository interface {
	FindOneByID(
		ctx context.Context,
		id uuid.UUID,
		queryParams *uctypes.QueryGetOneParams,
	) (budget *entity.Budget, err error)

	FindList(
		ctx context.Context,
		listOptions *BudgetListOptions,
		queryParams *uctypes.QueryGetListParams,
	) (items []*entity.Budget, err error)

	FindPagedList(
		ctx context.Context,
		listOptions *BudgetListOptions,
		queryParams *uctypes.QueryGetListParams,
	) (items []*entity.Budget, total uint64, err error)

	Create(ctx context.Context, item *entity.Budget) (err error)

	Update(ctx context.Context, item *entity.Budget) (err error)
}

type BudgetCacheRepository interface {
	SaveBudgetsList(ctx context.Context, key string, items []*entity.Budget, ttl *time.Duration) (err error)

	SaveBudgetsPagedList(ctx context.Context, key string, items []*entity.Budget, total uint64, ttl *time.Duration) (err error)

	SaveBudget(ctx context.Context, key string, item *entity.Budget, ttl *time.Duration) (err error)

	GetBudgetsList(ctx context.Context, key string) (items []*entity.Budget, err error)

	GetBudgetsPagedList(ctx context.Context, key string) (items []*entity.Budget, total uint64, err error)

	GetBudget(ctx context.Context, key string) (item *entity.Budget, err error)
}
