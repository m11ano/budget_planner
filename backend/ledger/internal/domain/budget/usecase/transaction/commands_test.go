package transaction

import (
	"context"
	"sync"
	"sync/atomic"
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

func TestTransactionUsecase_CreateTransactionByDTO_OK_income(t *testing.T) {
	t.Parallel()

	s := newDependencies(t)
	defer finishDependencies(s)

	accID := uuid.New()
	catID := uint64(10)

	s.categoryRepo.FindOneByIDMock.Set(func(ctx context.Context, id uint64, _ *uctypes.QueryGetOneParams) (*entity.Category, error) {
		require.Equal(t, catID, id)
		return &entity.Category{ID: id}, nil
	})

	s.transactionRepo.CreateMock.Set(func(ctx context.Context, item *entity.Transaction) error {
		require.Equal(t, accID, item.AccountID)
		require.True(t, item.IsIncome)
		require.Equal(t, catID, item.CategoryID)

		require.Equal(t, decimal.MustParse("100.12"), item.Amount)
		require.Equal(t, "hello", item.Description)
		require.Equal(t, civil.Date{Year: 2025, Month: 12, Day: 15}, item.OccurredOn)
		return nil
	})

	s.categoryRepo.FindListMock.Set(func(ctx context.Context, opt *usecase.CategoryListOptions, _ *uctypes.QueryGetListParams) ([]*entity.Category, error) {
		return []*entity.Category{{ID: catID}}, nil
	})

	got, err := s.uc.CreateTransactionByDTO(testCtx(), usecase.CreateTransactionDataInput{
		AccountID:   accID,
		IsIncome:    true,
		Amount:      decimal.MustParse("100.129"),
		OccurredOn:  civil.Date{Year: 2025, Month: 12, Day: 15},
		CategoryID:  catID,
		Description: "   hello   ",
	}, true)
	require.NoError(t, err)
	require.NotNil(t, got)
	require.NotNil(t, got.Transaction)
	require.Equal(t, accID, got.Transaction.AccountID)
	require.True(t, got.Transaction.IsIncome)
	require.Equal(t, decimal.MustParse("100.12"), got.Transaction.Amount)
	require.Equal(t, "hello", got.Transaction.Description)
}

func TestTransactionUsecase_CreateTransactionByDTO_OK_outcome_no_budget(t *testing.T) {
	t.Parallel()

	s := newDependencies(t)
	defer finishDependencies(s)

	accID := uuid.New()
	catID := uint64(7)

	s.categoryRepo.FindOneByIDMock.Set(func(ctx context.Context, id uint64, _ *uctypes.QueryGetOneParams) (*entity.Category, error) {
		require.Equal(t, catID, id)
		return &entity.Category{ID: id}, nil
	})

	s.budgetRepo.FindListMock.Set(func(ctx context.Context, opt *usecase.BudgetListOptions, qp *uctypes.QueryGetListParams) ([]*entity.Budget, error) {
		require.NotNil(t, opt)
		require.NotNil(t, opt.FilterAccountID)
		require.Equal(t, accID, *opt.FilterAccountID)

		require.NotNil(t, opt.FilterPeriod)
		require.Equal(t, civil.Date{Year: 2025, Month: 12, Day: 1}, *opt.FilterPeriod)

		require.NotNil(t, opt.FilterCategoryID)
		require.Equal(t, catID, *opt.FilterCategoryID)

		require.NotNil(t, qp)
		require.Equal(t, uint64(1), qp.Limit)

		return nil, nil
	})

	s.transactionRepo.CreateMock.Set(func(ctx context.Context, item *entity.Transaction) error {
		require.Equal(t, accID, item.AccountID)
		require.False(t, item.IsIncome)
		require.Equal(t, catID, item.CategoryID)
		require.Equal(t, decimal.MustParse("-12.34"), item.Amount)
		require.Equal(t, "coffee", item.Description)
		return nil
	})

	s.categoryRepo.FindListMock.Set(func(ctx context.Context, _ *usecase.CategoryListOptions, _ *uctypes.QueryGetListParams) ([]*entity.Category, error) {
		return []*entity.Category{{ID: catID}}, nil
	})

	got, err := s.uc.CreateTransactionByDTO(testCtx(), usecase.CreateTransactionDataInput{
		AccountID:   accID,
		IsIncome:    false,
		Amount:      decimal.MustParse("-12.349"),
		OccurredOn:  civil.Date{Year: 2025, Month: 12, Day: 20},
		CategoryID:  catID,
		Description: " coffee ",
	}, true)
	require.NoError(t, err)
	require.NotNil(t, got)
	require.NotNil(t, got.Transaction)
	require.Equal(t, decimal.MustParse("-12.34"), got.Transaction.Amount)
	require.Equal(t, "coffee", got.Transaction.Description)
}

func TestTransactionUsecase_PatchTransactionByDTO_OK(t *testing.T) {
	t.Parallel()

	s := newDependencies(t)
	defer finishDependencies(s)

	var wg sync.WaitGroup
	wg.Add(1)

	s.transactionCacheRepo.ClearForPrefixesMock.Set(func(ctx context.Context, prefixes ...string) error {
		defer wg.Done()
		return nil
	})

	id := uuid.New()
	accID := uuid.New()

	now := time.Now().Truncate(time.Microsecond)

	tx := &entity.Transaction{
		ID:          id,
		AccountID:   accID,
		IsIncome:    true,
		Amount:      decimal.MustParse("10.00"),
		OccurredOn:  civil.Date{Year: 2025, Month: 12, Day: 1},
		CategoryID:  2,
		Description: "old",
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	newAmount := decimal.MustParse("55.559")
	newDesc := "   new desc   "

	var call int32
	s.transactionRepo.FindOneByIDMock.Set(func(ctx context.Context, gotID uuid.UUID, qp *uctypes.QueryGetOneParams) (*entity.Transaction, error) {
		require.Equal(t, id, gotID)

		n := atomic.AddInt32(&call, 1)
		if n == 1 {
			require.Nil(t, qp)
			return tx, nil
		}

		require.NotNil(t, qp)
		require.True(t, qp.ForUpdate)
		return tx, nil
	})

	s.transactionRepo.UpdateMock.Set(func(ctx context.Context, got *entity.Transaction) error {
		require.Equal(t, decimal.MustParse("55.55"), got.Amount)
		require.Equal(t, "new desc", got.Description)
		return nil
	})

	err := s.uc.PatchTransactionByDTO(testCtx(), id, usecase.PatchTransactionDataInput{
		Version:     now.UnixMicro(),
		Amount:      &newAmount,
		Description: &newDesc,
	}, true)
	require.NoError(t, err)

	wg.Wait()
}

func TestTransactionUsecase_DeleteTransactionByID_OK(t *testing.T) {
	t.Parallel()

	s := newDependencies(t)
	defer finishDependencies(s)

	var wg sync.WaitGroup
	wg.Add(1)

	s.transactionCacheRepo.ClearForPrefixesMock.Set(func(ctx context.Context, prefixes ...string) error {
		defer wg.Done()
		return nil
	})

	id := uuid.New()
	accID := uuid.New()

	tx := &entity.Transaction{
		ID:         id,
		AccountID:  accID,
		IsIncome:   true,
		Amount:     decimal.MustParse("10.00"),
		OccurredOn: civil.Date{Year: 2025, Month: 12, Day: 1},
		CategoryID: 2,
	}

	var call int32
	s.transactionRepo.FindOneByIDMock.Set(func(ctx context.Context, gotID uuid.UUID, qp *uctypes.QueryGetOneParams) (*entity.Transaction, error) {
		require.Equal(t, id, gotID)

		n := atomic.AddInt32(&call, 1)
		if n == 1 {
			require.Nil(t, qp)
			return tx, nil
		}

		require.NotNil(t, qp)
		require.True(t, qp.ForUpdate)
		return tx, nil
	})

	s.transactionRepo.UpdateMock.Set(func(ctx context.Context, got *entity.Transaction) error {
		require.NotNil(t, got.DeletedAt)
		return nil
	})

	err := s.uc.DeleteTransactionByID(testCtx(), id)
	require.NoError(t, err)

	wg.Wait()
}

func TestTransactionUsecase_ImportTransactionsFromCSV_OK_empty(t *testing.T) {
	t.Parallel()

	s := newDependencies(t)
	defer finishDependencies(s)

	var wg sync.WaitGroup
	wg.Add(1)

	s.transactionCacheRepo.ClearForPrefixesMock.Set(func(ctx context.Context, prefixes ...string) error {
		defer wg.Done()
		return nil
	})

	accID := uuid.New()

	s.transactionCSVRepo.ItemsFromCSVMock.Set(func(ctx context.Context, data []byte, gotAcc uuid.UUID) ([]*entity.Transaction, error) {
		require.Equal(t, accID, gotAcc)
		require.Equal(t, []byte("csv-data"), data)
		return nil, nil
	})

	err := s.uc.ImportTransactionsFromCSV(testCtx(), []byte("csv-data"), accID)
	require.NoError(t, err)

	wg.Wait()
}
