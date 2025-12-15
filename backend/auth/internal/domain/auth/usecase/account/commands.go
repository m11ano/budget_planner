package account

import (
	"context"
	"errors"
	"net"

	"github.com/google/uuid"
	appErrors "github.com/m11ano/budget_planner/backend/auth/internal/app/errors"
	"github.com/m11ano/budget_planner/backend/auth/internal/common/uctypes"
	"github.com/m11ano/budget_planner/backend/auth/internal/domain/auth/entity"
	"github.com/m11ano/budget_planner/backend/auth/internal/domain/auth/usecase"
	"github.com/m11ano/budget_planner/backend/auth/internal/infra/loghandler"
)

func (uc *UsecaseImpl) CreateAccountByDTO(
	ctx context.Context,
	in usecase.CreateAccountDataInput,
	requestIP *net.IP,
) (*usecase.AccountDTO, error) {
	const op = "CreateAccountByDTO"

	account, err := entity.NewAccount(
		in.Email,
		in.Password,
		in.SkipPasswordCheck,
		in.IsConfirmed,
	)
	if err != nil {
		return nil, appErrors.Chainf(err, "%s.%s", uc.pkg, op)
	}

	err = account.SetProfileName(in.ProfileName)
	if err != nil {
		return nil, appErrors.Chainf(err, "%s.%s", uc.pkg, op)
	}

	err = account.SetProfileSurname(in.ProfileSurname)
	if err != nil {
		return nil, appErrors.Chainf(err, "%s.%s", uc.pkg, op)
	}

	err = uc.dbMasterClient.Do(ctx, func(ctx context.Context) error {
		_, err := uc.accountRepo.FindOneByEmail(ctx, account.Email, nil)
		if err == nil {
			return appErrors.ErrUniqueViolation.WithDetail("field", false, "email").WithHints("email is already in use")
		} else if !errors.Is(err, appErrors.ErrNotFound) {
			return err
		}

		err = uc.accountRepo.Create(ctx, account)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, appErrors.Chainf(err, "%s.%s", uc.pkg, op)
	}

	accountDTO, err := uc.entitiesToDTO(ctx, []*entity.Account{account})
	if err != nil {
		return nil, appErrors.Chainf(err, "%s.%s", uc.pkg, op)
	}

	if len(accountDTO) == 0 {
		uc.logger.ErrorContext(loghandler.WithSource(ctx), "unpredicted empty account dto")
		return nil, appErrors.Chainf(appErrors.ErrInternal, "%s.%s", uc.pkg, op)
	}

	return accountDTO[0], nil
}

func (uc *UsecaseImpl) PatchAccountByDTO(
	ctx context.Context,
	id uuid.UUID,
	in usecase.PatchAccountDataInput,
	skipVersionCheck bool,
) error {
	const op = "PatchAccountByDTO"

	err := uc.dbMasterClient.Do(ctx, func(ctx context.Context) error {
		account, err := uc.accountRepo.FindOneByID(ctx, id, &uctypes.QueryGetOneParams{
			ForUpdate: true,
		})
		if err != nil {
			return err
		}

		if !skipVersionCheck && account.Version() != in.Version {
			return appErrors.ErrVersionConflict.
				WithDetail("last_version", false, account.Version()).
				WithDetail("last_updated_at", false, account.UpdatedAt)
		}

		if in.Email != nil {
			err = account.SetEmail(*in.Email)
			if err != nil {
				return err
			}

			checkAccount, err := uc.accountRepo.FindOneByEmail(ctx, account.Email, nil)
			if err == nil && checkAccount.ID != account.ID {
				return appErrors.ErrUniqueViolation.WithDetail("field", false, "email")
			} else if err != nil && !errors.Is(err, appErrors.ErrNotFound) {
				return err
			}
		}

		if in.Password != nil {
			err = account.SetPassword(*in.Password, in.SkipPasswordCheck)
			if err != nil {
				return err
			}
		}

		if in.IsBlocked != nil {
			account.IsBlocked = *in.IsBlocked
		}

		if in.ProfileName != nil {
			err := account.SetProfileName(*in.ProfileName)
			if err != nil {
				return err
			}
		}

		if in.ProfileSurname != nil {
			err := account.SetProfileSurname(*in.ProfileSurname)
			if err != nil {
				return err
			}
		}

		err = uc.accountRepo.Update(ctx, account)
		if err != nil {
			return err
		}

		if in.Password != nil {
			err = uc.sessionUC.RevokeSessionsByAccountID(ctx, id)
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

func (uc *UsecaseImpl) UpdateAccount(ctx context.Context, item *entity.Account) error {
	const op = "UpdateAccount"

	err := uc.accountRepo.Update(ctx, item)
	if err != nil {
		return appErrors.Chainf(err, "%s.%s", uc.pkg, op)
	}

	return nil
}
