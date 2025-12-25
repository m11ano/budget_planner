package budget

import (
	"context"
	"testing"
	"time"

	"cloud.google.com/go/civil"
	"github.com/google/uuid"
	"github.com/govalues/decimal"
	"github.com/m11ano/budget_planner/backend/ledger/internal/common/uctypes"
	"github.com/m11ano/budget_planner/backend/ledger/internal/domain/budget/entity"
	"github.com/m11ano/budget_planner/backend/ledger/internal/domain/budget/usecase"
	"github.com/stretchr/testify/require"
)

func TestBudgetUsecase_CreateBudgetByDTO_Table(t *testing.T) {
	t.Parallel()

	accountID := uuid.New()

	tests := []struct {
		name string
		in   usecase.CreateBudgetDataInput

		categoryErr error

		wantErr bool
	}{
		{
			name: "OK",
			in: usecase.CreateBudgetDataInput{
				AccountID:  accountID,
				Period:     civil.Date{Year: 2025, Month: 12, Day: 15},
				CategoryID: 10,
				Amount:     decimal.MustParse("100.129"),
			},
			categoryErr: nil,
			wantErr:     false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			s := newDependencies(t)
			defer finishDependencies(s)

			s.categoryRepo.FindOneByIDMock.Set(func(ctx context.Context, id uint64, _ *uctypes.QueryGetOneParams) (*entity.Category, error) {
				require.Equal(t, tt.in.CategoryID, id)
				if tt.categoryErr != nil {
					return nil, tt.categoryErr
				}
				return &entity.Category{ID: id}, nil
			})

			if !tt.wantErr {
				s.budgetRepo.FindListMock.Set(func(ctx context.Context, opt *usecase.BudgetListOptions, qp *uctypes.QueryGetListParams) ([]*entity.Budget, error) {
					require.NotNil(t, opt.FilterAccountID)
					require.Equal(t, tt.in.AccountID, *opt.FilterAccountID)

					require.NotNil(t, opt.FilterPeriod)
					require.Equal(t,
						civil.Date{Year: tt.in.Period.Year, Month: tt.in.Period.Month, Day: 1},
						*opt.FilterPeriod,
					)

					require.NotNil(t, opt.FilterCategoryID)
					require.Equal(t, tt.in.CategoryID, *opt.FilterCategoryID)

					require.NotNil(t, qp)
					require.Equal(t, uint64(1), qp.Limit)

					return nil, nil
				})

				s.budgetRepo.CreateMock.Set(func(ctx context.Context, b *entity.Budget) error {
					require.Equal(t, tt.in.AccountID, b.AccountID)
					require.Equal(t, tt.in.CategoryID, b.CategoryID)

					require.Equal(t, 1, b.Period.Day)

					require.Equal(t, decimal.MustParse("100.12"), b.Amount)

					return nil
				})
			}

			called := make(chan struct{}, 1)
			s.budgetCacheRepo.ClearForPrefixesMock.Set(func(ctx context.Context, prefixes ...string) error {
				select {
				case called <- struct{}{}:
				default:
				}
				return nil
			})

			got, err := s.uc.CreateBudgetByDTO(testCtx(), tt.in)

			if tt.wantErr {
				require.Error(t, err)
				require.Nil(t, got)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, got)
			require.NotNil(t, got.Budget)

			select {
			case <-called:
			case <-time.After(250 * time.Millisecond):
				t.Fatalf("expected budgetCacheRepo.ClearForPrefixes to be called")
			}
		})
	}
}

func TestBudgetUsecase_PatchBudgetByDTO_OK(t *testing.T) {
	t.Parallel()

	s := newDependencies(t)
	defer finishDependencies(s)

	id := uuid.New()
	accountID := uuid.New()

	now := time.Now().Truncate(time.Microsecond)
	b := &entity.Budget{
		ID:         id,
		AccountID:  accountID,
		CategoryID: 10,
		Period:     civil.Date{Year: 2025, Month: 12, Day: 1},
		Amount:     decimal.MustParse("10.00"),
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	newAmount := decimal.MustParse("55.559")

	s.budgetRepo.FindOneByIDMock.Set(func(ctx context.Context, gotID uuid.UUID, qp *uctypes.QueryGetOneParams) (*entity.Budget, error) {
		require.Equal(t, id, gotID)

		if qp == nil {
			return b, nil
		}

		require.True(t, qp.ForUpdate)
		return b, nil
	})

	s.budgetRepo.UpdateMock.Set(func(ctx context.Context, got *entity.Budget) error {
		require.Equal(t, decimal.MustParse("55.55"), got.Amount)
		return nil
	})

	called := make(chan struct{}, 1)
	s.budgetCacheRepo.ClearForPrefixesMock.Set(func(ctx context.Context, prefixes ...string) error {
		select {
		case called <- struct{}{}:
		default:
		}
		return nil
	})

	err := s.uc.PatchBudgetByDTO(testCtx(), id, usecase.PatchBudgetDataInput{
		Version: now.UnixMicro(),
		Amount:  &newAmount,
	}, true)
	require.NoError(t, err)

	select {
	case <-called:
	case <-time.After(250 * time.Millisecond):
		t.Fatalf("expected budgetCacheRepo.ClearForPrefixes to be called")
	}
}

func TestBudgetUsecase_DeleteBudgetByID_OK(t *testing.T) {
	t.Parallel()

	s := newDependencies(t)
	defer finishDependencies(s)

	id := uuid.New()
	accountID := uuid.New()

	b := &entity.Budget{
		ID:         id,
		AccountID:  accountID,
		CategoryID: 10,
		Period:     civil.Date{Year: 2025, Month: 12, Day: 1},
		Amount:     decimal.MustParse("10.00"),
	}

	s.budgetRepo.FindOneByIDMock.Set(func(ctx context.Context, gotID uuid.UUID, qp *uctypes.QueryGetOneParams) (*entity.Budget, error) {
		require.Equal(t, id, gotID)

		if qp == nil {
			return b, nil
		}

		require.True(t, qp.ForUpdate)
		return b, nil
	})

	s.budgetRepo.UpdateMock.Set(func(ctx context.Context, got *entity.Budget) error {
		require.NotNil(t, got.DeletedAt)
		return nil
	})

	called := make(chan struct{}, 1)
	s.budgetCacheRepo.ClearForPrefixesMock.Set(func(ctx context.Context, prefixes ...string) error {
		select {
		case called <- struct{}{}:
		default:
		}
		return nil
	})

	err := s.uc.DeleteBudgetByID(testCtx(), id)
	require.NoError(t, err)

	select {
	case <-called:
	case <-time.After(250 * time.Millisecond):
		t.Fatalf("expected budgetCacheRepo.ClearForPrefixes to be called")
	}
}
