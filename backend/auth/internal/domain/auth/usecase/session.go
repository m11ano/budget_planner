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

//go:generate minimock -i github.com/m11ano/budget_planner/backend/auth/internal/domain/auth/usecase.SessionUsecase -o mocks/session_usecase.go
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

//go:generate minimock -i github.com/m11ano/budget_planner/backend/auth/internal/domain/auth/usecase.SessionRepository -o mocks/session_repository.go
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
