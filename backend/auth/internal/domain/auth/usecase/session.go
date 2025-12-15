package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/m11ano/budget_planner/backend/auth/internal/common/uctypes"
	"github.com/m11ano/budget_planner/backend/auth/internal/domain/auth/entity"
)

type SessionListOptions struct {
	FilterAccountID *uuid.UUID
}

type SessionUsecase interface {
	Create(
		ctx context.Context,
		item *entity.Session,
	) (err error)

	Update(
		ctx context.Context,
		item *entity.Session,
	) (err error)

	FindOneByID(
		ctx context.Context,
		id uuid.UUID,
		queryParams *uctypes.QueryGetOneParams,
	) (session *entity.Session, err error)

	FindList(
		ctx context.Context,
		listOptions *SessionListOptions,
		queryParams *uctypes.QueryGetListParams,
	) (items []*entity.Session, err error)

	RevokeSessionsByAccountID(ctx context.Context, accountID uuid.UUID) (resErr error)

	RevokeSessionByID(ctx context.Context, ID uuid.UUID) (resErr error)
}

type SessionRepository interface {
	Create(
		ctx context.Context,
		item *entity.Session,
	) (err error)

	Update(
		ctx context.Context,
		item *entity.Session,
	) (err error)

	FindOneByID(
		ctx context.Context,
		id uuid.UUID,
		queryParams *uctypes.QueryGetOneParams,
	) (session *entity.Session, err error)

	FindList(
		ctx context.Context,
		listOptions *SessionListOptions,
		queryParams *uctypes.QueryGetListParams,
	) (items []*entity.Session, err error)
}
