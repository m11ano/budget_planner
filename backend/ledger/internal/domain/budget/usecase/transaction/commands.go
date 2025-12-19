package transaction

import (
	"context"
	"errors"
	"fmt"
	"time"

	"cloud.google.com/go/civil"
	"github.com/google/uuid"
	"github.com/m11ano/budget_planner/backend/auth/pkg/auth"
	appErrors "github.com/m11ano/budget_planner/backend/ledger/internal/app/errors"
	"github.com/m11ano/budget_planner/backend/ledger/internal/common/uctypes"
	"github.com/m11ano/budget_planner/backend/ledger/internal/domain/budget/entity"
	"github.com/m11ano/budget_planner/backend/ledger/internal/domain/budget/usecase"
	"github.com/m11ano/budget_planner/backend/ledger/internal/infra/loghandler"
	"github.com/m11ano/budget_planner/backend/ledger/pkg/pgclient"
	"github.com/samber/lo"
)

func (uc *UsecaseImpl) CreateTransactionByDTO(
	ctx context.Context,
	in usecase.CreateTransactionDataInput,
) (*usecase.TransactionDTO, error) {
	const op = "CreateTransactionByDTO"

	if auth.IsNeedToCheckRights(ctx) {
		authData := auth.GetAuthData(ctx)
		if authData == nil || authData.AccountID != in.AccountID {
			return nil, appErrors.ErrForbidden
		}
	}

	transaction, err := entity.NewTransaction(
		in.AccountID,
		in.IsIncome,
		in.Amount,
		in.OccurredOn,
		in.CategoryID,
	)
	if err != nil {
		return nil, appErrors.Chainf(err, "%s.%s", uc.pkg, op)
	}

	err = transaction.SetDescription(in.Description)
	if err != nil {
		return nil, appErrors.Chainf(err, "%s.%s", uc.pkg, op)
	}

	err = uc.dbMasterClient.DoWithIsoLvl(ctx, pgclient.Serializable, func(ctx context.Context) error {
		_, err := uc.categoryRepo.FindOneByID(ctx, transaction.CategoryID, nil)
		if err != nil {
			if errors.Is(err, appErrors.ErrNotFound) {
				return appErrors.Chainf(appErrors.ErrBadRequest.WithHints("category not found"), "%s.%s", uc.pkg, op)
			}
			return err
		}

		if !in.IsIncome {
			period := civil.Date{
				Year:  transaction.OccurredOn.Year,
				Month: transaction.OccurredOn.Month,
				Day:   1,
			}

			budgetCheck, err := uc.budgetRepo.FindList(ctx, &usecase.BudgetListOptions{
				FilterAccountID:  &transaction.AccountID,
				FilterPeriod:     &period,
				FilterCategoryID: &transaction.CategoryID,
			}, &uctypes.QueryGetListParams{
				Limit: 1,
			})
			if err != nil {
				return err
			}

			if len(budgetCheck) > 0 {
				budget := budgetCheck[0]

				periodEnd := period.AddMonths(1).AddDays(-1)

				reports, err := uc.transactionRepo.CountReportItems(ctx, &usecase.CountReportItemsFilter{
					AccountID:  &transaction.AccountID,
					PeriodFrom: &period,
					PeriodTo:   &periodEnd,
					CategoryID: &transaction.CategoryID,
				})
				if err != nil {
					return err
				}

				balance := transaction.Amount
				for _, report := range reports {
					balance, err = balance.Add(report.Sum)
					if err != nil {
						return err
					}
				}

				if balance.Cmp(budget.Amount.Neg()) == -1 {
					return appErrors.Chainf(appErrors.ErrBadRequest.WithHints("budget limit exceeded"), "%s.%s", uc.pkg, op)
				}
			}
		}

		err = uc.transactionRepo.Create(ctx, transaction)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, appErrors.Chainf(err, "%s.%s", uc.pkg, op)
	}

	transactionDTO, err := uc.entitiesToDTO(ctx, []*entity.Transaction{transaction})
	if err != nil {
		return nil, appErrors.Chainf(err, "%s.%s", uc.pkg, op)
	}

	if len(transactionDTO) == 0 {
		uc.logger.ErrorContext(loghandler.WithSource(ctx), "unpredicted empty transaction dto")
		return nil, appErrors.Chainf(appErrors.ErrInternal, "%s.%s", uc.pkg, op)
	}

	return transactionDTO[0], nil
}

func (uc *UsecaseImpl) PatchTransactionByDTO(
	ctx context.Context,
	id uuid.UUID,
	in usecase.PatchTransactionDataInput,
	skipVersionCheck bool,
) error {
	const op = "PatchTransactionByDTO"

	err := uc.dbMasterClient.DoWithIsoLvl(ctx, pgclient.Serializable, func(ctx context.Context) error {
		transaction, err := uc.transactionRepo.FindOneByID(ctx, id, &uctypes.QueryGetOneParams{
			ForUpdate: true,
		})
		if err != nil {
			return err
		}

		if auth.IsNeedToCheckRights(ctx) {
			authData := auth.GetAuthData(ctx)
			if authData == nil || authData.AccountID != transaction.AccountID {
				return appErrors.ErrForbidden
			}
		}

		if !skipVersionCheck && transaction.Version() != in.Version {
			return appErrors.ErrVersionConflict.
				WithDetail("last_version", false, transaction.Version()).
				WithDetail("last_updated_at", false, transaction.UpdatedAt)
		}

		// categoryChanged := false
		if in.CategoryID != nil {
			_, err := uc.categoryRepo.FindOneByID(ctx, *in.CategoryID, nil)
			if err != nil {
				if errors.Is(err, appErrors.ErrNotFound) {
					return appErrors.Chainf(appErrors.ErrBadRequest.WithHints("category not found"), "%s.%s", uc.pkg, op)
				}
				return err
			}

			// if transaction.CategoryID != *in.CategoryID {
			// 	categoryChanged = true
			// }

			transaction.CategoryID = *in.CategoryID
		}

		// periodChanged := false
		if in.OccurredOn != nil {
			// if transaction.OccurredOn.Month != in.OccurredOn.Month || transaction.OccurredOn.Year != in.OccurredOn.Year {
			// 	periodChanged = true
			// }

			err := transaction.SetOccuredOn(*in.OccurredOn)
			if err != nil {
				return err
			}
		}

		if in.Amount != nil {
			if !transaction.IsIncome {
				period := civil.Date{
					Year:  transaction.OccurredOn.Year,
					Month: transaction.OccurredOn.Month,
					Day:   1,
				}

				budgetCheck, err := uc.budgetRepo.FindList(ctx, &usecase.BudgetListOptions{
					FilterAccountID:  &transaction.AccountID,
					FilterPeriod:     &period,
					FilterCategoryID: &transaction.CategoryID,
				}, &uctypes.QueryGetListParams{
					Limit: 1,
				})
				if err != nil {
					return err
				}

				if len(budgetCheck) > 0 {
					budget := budgetCheck[0]

					periodEnd := period.AddMonths(1).AddDays(-1)

					reports, err := uc.transactionRepo.CountReportItems(ctx, &usecase.CountReportItemsFilter{
						AccountID:  &transaction.AccountID,
						PeriodFrom: &period,
						PeriodTo:   &periodEnd,
						CategoryID: &transaction.CategoryID,
						ExcludeIDs: []uuid.UUID{transaction.ID},
					})
					if err != nil {
						return err
					}

					balance := *in.Amount
					for _, report := range reports {
						fmt.Printf("%#v", report)
						balance, err = balance.Add(report.Sum)
						if err != nil {
							return err
						}
					}

					// if !categoryChanged && !periodChanged {
					// 	balance, err = balance.Add(transaction.Amount.Neg())
					// 	if err != nil {
					// 		return err
					// 	}
					// }

					if balance.Cmp(budget.Amount.Neg()) == -1 {
						return appErrors.Chainf(appErrors.ErrBadRequest.WithHints("budget limit exceeded"), "%s.%s", uc.pkg, op)
					}
				}
			}

			err = transaction.SetAmount(*in.Amount)
			if err != nil {
				return err
			}
		}

		if in.Description != nil {
			err = transaction.SetDescription(*in.Description)
			if err != nil {
				return err
			}
		}

		err = uc.transactionRepo.Update(ctx, transaction)
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

func (uc *UsecaseImpl) DeleteTransactionByID(ctx context.Context, id uuid.UUID) error {
	const op = "DeleteTransactionByID"

	err := uc.dbMasterClient.Do(ctx, func(ctx context.Context) error {
		transaction, err := uc.transactionRepo.FindOneByID(ctx, id, &uctypes.QueryGetOneParams{
			ForUpdate: true,
		})
		if err != nil {
			return err
		}

		if auth.IsNeedToCheckRights(ctx) {
			authData := auth.GetAuthData(ctx)
			if authData == nil || authData.AccountID != transaction.AccountID {
				return appErrors.ErrForbidden
			}
		}

		transaction.DeletedAt = lo.ToPtr(time.Now())

		err = uc.transactionRepo.Update(ctx, transaction)
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
