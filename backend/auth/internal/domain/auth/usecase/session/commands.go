package session

import (
	"context"

	"github.com/google/uuid"
	appErrors "github.com/m11ano/budget_planner/backend/auth/internal/app/errors"
	"github.com/m11ano/budget_planner/backend/auth/internal/common/uctypes"
	"github.com/m11ano/budget_planner/backend/auth/internal/domain/auth/entity"
	"github.com/m11ano/budget_planner/backend/auth/internal/domain/auth/usecase"
)

func (uc *UsecaseImpl) Create(
	ctx context.Context,
	item *entity.Session,
) error {
	const op = "Create"

	err := uc.repo.Create(ctx, item)
	if err != nil {
		return appErrors.Chainf(err, "%s.%s", uc.pkg, op)
	}

	return nil
}

func (uc *UsecaseImpl) Update(
	ctx context.Context,
	item *entity.Session,
) error {
	const op = "Update"

	err := uc.repo.Update(ctx, item)
	if err != nil {
		return appErrors.Chainf(err, "%s.%s", uc.pkg, op)
	}

	return nil
}

func (uc *UsecaseImpl) RevokeSessionsByAccountID(ctx context.Context, accountID uuid.UUID) error {
	const op = "RevokeSessionsByAccountID"

	err := uc.dbMasterClient.Do(ctx, func(ctx context.Context) error {
		sessions, err := uc.repo.FindList(ctx, &usecase.SessionListOptions{
			FilterAccountID: &accountID,
		}, &uctypes.QueryGetListParams{
			ForUpdate: true,
		})
		if err != nil {
			return err
		}

		for _, session := range sessions {
			err := uc.delete(ctx, session)
			if err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		return appErrors.Chainf(err, "%s.%s", uc.pkg, op)
	}

	return nil
}

func (uc *UsecaseImpl) RevokeSessionByID(ctx context.Context, ID uuid.UUID) error {
	const op = "RevokeSessionByID"

	err := uc.dbMasterClient.Do(ctx, func(ctx context.Context) error {
		session, err := uc.repo.FindOneByID(ctx, ID, &uctypes.QueryGetOneParams{
			ForUpdate: true,
		})
		if err != nil {
			return err
		}

		err = uc.delete(ctx, session)
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
