package auth

import (
	"context"
	"errors"

	"github.com/google/uuid"
	appErrors "github.com/m11ano/budget_planner/backend/auth/internal/app/errors"
)

func (uc *UsecaseImpl) IsSessionConfirmed(
	ctx context.Context,
	id uuid.UUID,
) (bool, error) {
	const op = "IsSessionConfirmed"

	session, err := uc.sessionUC.FindOneByID(ctx, id, nil)
	if err != nil {
		if errors.Is(err, appErrors.ErrNotFound) {
			return false, nil
		}
		return false, appErrors.Chainf(err, "%s.%s", uc.pkg, op)
	}

	accountDTO, err := uc.accountUC.FindOneByID(ctx, session.AccountID, nil)
	if err != nil {
		return false, appErrors.Chainf(err, "%s.%s", uc.pkg, op)
	}

	if !accountDTO.Account.IsConfirmed || accountDTO.Account.IsBlocked {
		return false, nil
	}

	return true, nil
}
