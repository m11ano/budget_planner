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

type TransactionListOptions struct {
	FilterAccountID      *uuid.UUID
	FilterOccurredOnFrom *civil.Date
	FilterOccurredOnTo   *civil.Date
	Sort                 []uctypes.SortOption[TransactionListOptionsSortField]
}

type TransactionListOptionsSortField string

const (
	TransactionListOptionsSortFieldOccurredOn TransactionListOptionsSortField = "occurred_on"
	TransactionListOptionsSortFieldCreatedAt  TransactionListOptionsSortField = "created_at"
)

type TransactionDTO struct {
	Transaction *entity.Transaction
}

type CreateTransactionDataInput struct {
	AccountID   uuid.UUID
	IsIncome    bool
	Amount      decimal.Decimal
	OccurredOn  civil.Date
	CategoryID  uint64
	Description string
}

type PatchTransactionDataInput struct {
	Version int64

	Amount      *decimal.Decimal
	OccurredOn  *civil.Date
	CategoryID  *uint64
	Description *string
}

type CountReportItemsQueryFilter struct {
	AccountID  uuid.UUID
	DateFrom   *civil.Date
	DateTo     *civil.Date
	CategoryID *uint64
	ExcludeIDs []uuid.UUID
}

type TransactionUsecase interface {
	FindOneByID(
		ctx context.Context,
		id uuid.UUID,
		queryParams *uctypes.QueryGetOneParams,
	) (resItem *TransactionDTO, resErr error)

	FindList(
		ctx context.Context,
		listOptions *TransactionListOptions,
		queryParams *uctypes.QueryGetListParams,
	) (resItems []*TransactionDTO, resErr error)

	FindPagedList(
		ctx context.Context,
		listOptions *TransactionListOptions,
		queryParams *uctypes.QueryGetListParams,
	) (resItems []*TransactionDTO, total uint64, resErr error)

	FindListInMap(
		ctx context.Context,
		listOptions *TransactionListOptions,
		queryParams *uctypes.QueryGetListParams,
	) (resItems map[uuid.UUID]*TransactionDTO, resErr error)

	CreateTransactionByDTO(
		ctx context.Context,
		in CreateTransactionDataInput,
	) (resTransactionDTO *TransactionDTO, resErr error)

	PatchTransactionByDTO(
		ctx context.Context,
		id uuid.UUID,
		in PatchTransactionDataInput,
		skipVersionCheck bool,
	) (resErr error)

	DeleteTransactionByID(
		ctx context.Context,
		id uuid.UUID,
	) (resErr error)

	CountReportItems(
		ctx context.Context,
		queryFilter CountReportItemsQueryFilter,
	) (items []*entity.ReportItem, cacheHit bool, err error)
}

type TransactionRepository interface {
	FindOneByID(
		ctx context.Context,
		id uuid.UUID,
		queryParams *uctypes.QueryGetOneParams,
	) (transaction *entity.Transaction, err error)

	FindList(
		ctx context.Context,
		listOptions *TransactionListOptions,
		queryParams *uctypes.QueryGetListParams,
	) (items []*entity.Transaction, err error)

	FindPagedList(
		ctx context.Context,
		listOptions *TransactionListOptions,
		queryParams *uctypes.QueryGetListParams,
	) (items []*entity.Transaction, total uint64, err error)

	Create(ctx context.Context, item *entity.Transaction) (err error)

	Update(ctx context.Context, item *entity.Transaction) (err error)

	CountReportItems(
		ctx context.Context,
		queryFilter CountReportItemsQueryFilter,
	) (items []*entity.AccountTransactionReportItem, err error)
}

type TransactionCacheRepository interface {
	SaveReports(ctx context.Context, key string, items []*entity.ReportItem, ttl *time.Duration) (err error)

	GetReports(ctx context.Context, key string) (items []*entity.ReportItem, err error)
}
