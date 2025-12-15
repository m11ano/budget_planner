package usecase

import (
	"context"
	"net"

	"github.com/google/uuid"
	"github.com/m11ano/budget_planner/backend/auth/internal/common/uctypes"
	"github.com/m11ano/budget_planner/backend/auth/internal/domain/auth/entity"
)

type AccountListOptions struct {
	FilterIDs *[]uuid.UUID
	Sort      []uctypes.SortOption[AccountListOptionsSortField]
}

type AccountListOptionsSortField string

const (
	AccountListOptionsSortFieldCreatedAt AccountListOptionsSortField = "created_at"
)

type AccountDTO struct {
	Account *entity.Account
}

type CreateAccountDataInput struct {
	Email             string
	Password          string
	SkipPasswordCheck bool
	IsConfirmed       bool
	IsBlocked         bool
	ProfileName       string
	ProfileSurname    string
}

type PatchAccountDataInput struct {
	Version int64

	Email             *string
	Password          *string
	SkipPasswordCheck bool
	IsConfirmed       *bool
	IsBlocked         *bool
	ProfileName       *string
	ProfileSurname    *string
}

type AccountUsecase interface {
	FindOneByEmail(
		ctx context.Context,
		email string,
		queryParams *uctypes.QueryGetOneParams,
	) (resAccount *AccountDTO, resErr error)

	FindOneByID(
		ctx context.Context,
		id uuid.UUID,
		queryParams *uctypes.QueryGetOneParams,
	) (resAccount *AccountDTO, resErr error)

	FindList(
		ctx context.Context,
		listOptions *AccountListOptions,
		queryParams *uctypes.QueryGetListParams,
	) (resItems []*AccountDTO, resErr error)

	FindPagedList(
		ctx context.Context,
		listOptions *AccountListOptions,
		queryParams *uctypes.QueryGetListParams,
	) (resItems []*AccountDTO, total uint64, resErr error)

	FindListInMap(
		ctx context.Context,
		listOptions *AccountListOptions,
		queryParams *uctypes.QueryGetListParams,
	) (resItems map[uuid.UUID]*AccountDTO, resErr error)

	CreateAccountByDTO(
		ctx context.Context,
		in CreateAccountDataInput,
		requestIP *net.IP,
	) (resAccountDTO *AccountDTO, resErr error)

	PatchAccountByDTO(
		ctx context.Context,
		id uuid.UUID,
		in PatchAccountDataInput,
		skipVersionCheck bool,
	) (resErr error)

	UpdateAccount(ctx context.Context, item *entity.Account) (resErr error)
}

type AccountRepository interface {
	FindOneByEmail(
		ctx context.Context,
		email string,
		queryParams *uctypes.QueryGetOneParams,
	) (account *entity.Account, err error)

	FindOneByID(
		ctx context.Context,
		id uuid.UUID,
		queryParams *uctypes.QueryGetOneParams,
	) (account *entity.Account, err error)

	FindList(
		ctx context.Context,
		listOptions *AccountListOptions,
		queryParams *uctypes.QueryGetListParams,
	) (items []*entity.Account, err error)

	FindPagedList(
		ctx context.Context,
		listOptions *AccountListOptions,
		queryParams *uctypes.QueryGetListParams,
	) (items []*entity.Account, total uint64, err error)

	Create(ctx context.Context, item *entity.Account) (err error)

	Update(ctx context.Context, item *entity.Account) (err error)
}
