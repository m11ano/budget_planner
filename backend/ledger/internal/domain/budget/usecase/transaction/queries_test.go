package transaction

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

func TestTransactionUsecase_FindOneByID_OK(t *testing.T) {
	t.Parallel()

	s := newDependencies(t)
	defer finishDependencies(s)

	id := uuid.New()
	accID := uuid.New()
	catID := uint64(10)

	tx := &entity.Transaction{
		ID:         id,
		AccountID:  accID,
		IsIncome:   true,
		Amount:     decimal.MustParse("123.45"),
		OccurredOn: civil.Date{Year: 2025, Month: 12, Day: 20},
		CategoryID: catID,
		CreatedAt:  time.Now().Truncate(time.Microsecond),
		UpdatedAt:  time.Now().Truncate(time.Microsecond),
	}

	cat := &entity.Category{ID: catID}

	s.transactionRepo.FindOneByIDMock.Set(func(ctx context.Context, gotID uuid.UUID, qp *uctypes.QueryGetOneParams) (*entity.Transaction, error) {
		require.Equal(t, id, gotID)
		require.Nil(t, qp)
		return tx, nil
	})

	s.categoryRepo.FindListMock.Set(func(ctx context.Context, opt *usecase.CategoryListOptions, qp *uctypes.QueryGetListParams) ([]*entity.Category, error) {
		require.NotNil(t, opt)
		require.NotNil(t, opt.FilterIDs)
		require.Len(t, *opt.FilterIDs, 1)
		require.Equal(t, catID, (*opt.FilterIDs)[0])
		require.Nil(t, qp)
		return []*entity.Category{cat}, nil
	})

	got, err := s.uc.FindOneByID(testCtx(), id, nil)
	require.NoError(t, err)
	require.NotNil(t, got)
	require.Equal(t, tx, got.Transaction)
	require.Equal(t, cat, got.Category)
}

func TestTransactionUsecase_FindList_OK(t *testing.T) {
	t.Parallel()

	s := newDependencies(t)
	defer finishDependencies(s)

	accID := uuid.New()
	cat1 := uint64(2)
	cat2 := uint64(7)

	listOptions := &usecase.TransactionListOptions{FilterAccountID: &accID}
	qp := &uctypes.QueryGetListParams{Limit: 50}

	items := []*entity.Transaction{
		{
			ID:         uuid.New(),
			AccountID:  accID,
			IsIncome:   true,
			Amount:     decimal.MustParse("10.00"),
			OccurredOn: civil.Date{Year: 2025, Month: 12, Day: 1},
			CategoryID: cat1,
			CreatedAt:  time.Now().Truncate(time.Microsecond),
			UpdatedAt:  time.Now().Truncate(time.Microsecond),
		},
		{
			ID:         uuid.New(),
			AccountID:  accID,
			IsIncome:   false,
			Amount:     decimal.MustParse("-20.00"),
			OccurredOn: civil.Date{Year: 2025, Month: 12, Day: 2},
			CategoryID: cat2,
			CreatedAt:  time.Now().Truncate(time.Microsecond),
			UpdatedAt:  time.Now().Truncate(time.Microsecond),
		},
	}

	s.transactionRepo.FindListMock.Set(func(ctx context.Context, gotOpt *usecase.TransactionListOptions, gotQP *uctypes.QueryGetListParams) ([]*entity.Transaction, error) {
		require.Equal(t, listOptions, gotOpt)
		require.Equal(t, qp, gotQP)
		return items, nil
	})

	s.categoryRepo.FindListMock.Set(func(ctx context.Context, opt *usecase.CategoryListOptions, _ *uctypes.QueryGetListParams) ([]*entity.Category, error) {
		require.NotNil(t, opt)
		require.NotNil(t, opt.FilterIDs)
		ids := *opt.FilterIDs
		require.Len(t, ids, 2)
		require.ElementsMatch(t, []uint64{cat1, cat2}, ids)

		return []*entity.Category{
			{ID: cat1},
			{ID: cat2},
		}, nil
	})

	got, err := s.uc.FindList(testCtx(), listOptions, qp)
	require.NoError(t, err)
	require.Len(t, got, 2)
	require.Equal(t, items[0], got[0].Transaction)
	require.Equal(t, cat1, got[0].Category.ID)
	require.Equal(t, items[1], got[1].Transaction)
	require.Equal(t, cat2, got[1].Category.ID)
}

func TestTransactionUsecase_FindPagedList_OK(t *testing.T) {
	t.Parallel()

	s := newDependencies(t)
	defer finishDependencies(s)

	accID := uuid.New()
	catID := uint64(7)

	listOptions := &usecase.TransactionListOptions{FilterAccountID: &accID}
	qp := &uctypes.QueryGetListParams{Limit: 10, Offset: 20}

	items := []*entity.Transaction{
		{
			ID:         uuid.New(),
			AccountID:  accID,
			IsIncome:   true,
			Amount:     decimal.MustParse("99.99"),
			OccurredOn: civil.Date{Year: 2025, Month: 11, Day: 30},
			CategoryID: catID,
			CreatedAt:  time.Now().Truncate(time.Microsecond),
			UpdatedAt:  time.Now().Truncate(time.Microsecond),
		},
	}
	total := uint64(42)

	s.transactionRepo.FindPagedListMock.Set(func(ctx context.Context, gotOpt *usecase.TransactionListOptions, gotQP *uctypes.QueryGetListParams) ([]*entity.Transaction, uint64, error) {
		require.Equal(t, listOptions, gotOpt)
		require.Equal(t, qp, gotQP)
		return items, total, nil
	})

	s.categoryRepo.FindListMock.Set(func(ctx context.Context, opt *usecase.CategoryListOptions, _ *uctypes.QueryGetListParams) ([]*entity.Category, error) {
		require.NotNil(t, opt)
		require.NotNil(t, opt.FilterIDs)
		require.ElementsMatch(t, []uint64{catID}, *opt.FilterIDs)
		return []*entity.Category{{ID: catID}}, nil
	})

	got, gotTotal, err := s.uc.FindPagedList(testCtx(), listOptions, qp)
	require.NoError(t, err)
	require.Equal(t, total, gotTotal)
	require.Len(t, got, 1)
	require.Equal(t, items[0], got[0].Transaction)
	require.Equal(t, catID, got[0].Category.ID)
}

func TestTransactionUsecase_FindListInMap_OK(t *testing.T) {
	t.Parallel()

	s := newDependencies(t)
	defer finishDependencies(s)

	accID := uuid.New()
	cat1 := uint64(2)
	cat2 := uint64(3)

	listOptions := &usecase.TransactionListOptions{FilterAccountID: &accID}
	qp := &uctypes.QueryGetListParams{Limit: 100}

	items := []*entity.Transaction{
		{
			ID:         uuid.New(),
			AccountID:  accID,
			IsIncome:   true,
			Amount:     decimal.MustParse("10.00"),
			OccurredOn: civil.Date{Year: 2025, Month: 12, Day: 10},
			CategoryID: cat1,
			CreatedAt:  time.Now().Truncate(time.Microsecond),
			UpdatedAt:  time.Now().Truncate(time.Microsecond),
		},
		{
			ID:         uuid.New(),
			AccountID:  accID,
			IsIncome:   false,
			Amount:     decimal.MustParse("-5.00"),
			OccurredOn: civil.Date{Year: 2025, Month: 12, Day: 11},
			CategoryID: cat2,
			CreatedAt:  time.Now().Truncate(time.Microsecond),
			UpdatedAt:  time.Now().Truncate(time.Microsecond),
		},
	}

	s.transactionRepo.FindListMock.Set(func(ctx context.Context, gotOpt *usecase.TransactionListOptions, gotQP *uctypes.QueryGetListParams) ([]*entity.Transaction, error) {
		require.Equal(t, listOptions, gotOpt)
		require.Equal(t, qp, gotQP)
		return items, nil
	})

	s.categoryRepo.FindListMock.Set(func(ctx context.Context, opt *usecase.CategoryListOptions, _ *uctypes.QueryGetListParams) ([]*entity.Category, error) {
		require.NotNil(t, opt)
		require.NotNil(t, opt.FilterIDs)
		require.ElementsMatch(t, []uint64{cat1, cat2}, *opt.FilterIDs)
		return []*entity.Category{{ID: cat1}, {ID: cat2}}, nil
	})

	got, err := s.uc.FindListInMap(testCtx(), listOptions, qp)
	require.NoError(t, err)
	require.Len(t, got, 2)
	require.NotNil(t, got[items[0].ID])
	require.NotNil(t, got[items[1].ID])
	require.Equal(t, items[0], got[items[0].ID].Transaction)
	require.Equal(t, items[1], got[items[1].ID].Transaction)
}

func TestTransactionUsecase_FindPagedListAsCSV_OK(t *testing.T) {
	t.Parallel()

	s := newDependencies(t)
	defer finishDependencies(s)

	accID := uuid.New()
	catID := uint64(3)

	listOptions := &usecase.TransactionListOptions{FilterAccountID: &accID}
	qp := &uctypes.QueryGetListParams{Limit: 10, Offset: 0}

	items := []*entity.Transaction{
		{
			ID:         uuid.New(),
			AccountID:  accID,
			IsIncome:   true,
			Amount:     decimal.MustParse("1.11"),
			OccurredOn: civil.Date{Year: 2025, Month: 12, Day: 1},
			CategoryID: catID,
			CreatedAt:  time.Now().Truncate(time.Microsecond),
			UpdatedAt:  time.Now().Truncate(time.Microsecond),
		},
	}
	total := uint64(1)

	s.transactionRepo.FindPagedListMock.Set(func(ctx context.Context, gotOpt *usecase.TransactionListOptions, gotQP *uctypes.QueryGetListParams) ([]*entity.Transaction, uint64, error) {
		require.Equal(t, listOptions, gotOpt)
		require.Equal(t, qp, gotQP)
		return items, total, nil
	})

	s.categoryRepo.FindListMock.Set(func(ctx context.Context, opt *usecase.CategoryListOptions, _ *uctypes.QueryGetListParams) ([]*entity.Category, error) {
		require.NotNil(t, opt)
		require.NotNil(t, opt.FilterIDs)
		require.ElementsMatch(t, []uint64{catID}, *opt.FilterIDs)
		return []*entity.Category{{ID: catID}}, nil
	})

	wantCSV := []byte("csv,data\n")
	s.transactionCSVRepo.ItemsToCSVMock.Set(func(ctx context.Context, gotItems []*usecase.TransactionDTO) ([]byte, error) {
		require.Len(t, gotItems, 1)
		require.Equal(t, items[0], gotItems[0].Transaction)
		require.NotNil(t, gotItems[0].Category)
		require.Equal(t, catID, gotItems[0].Category.ID)
		return wantCSV, nil
	})

	data, gotTotal, err := s.uc.FindPagedListAsCSV(testCtx(), listOptions, qp)
	require.NoError(t, err)
	require.Equal(t, total, gotTotal)
	require.Equal(t, wantCSV, data)
}

func TestTransactionUsecase_CountReportItems_CacheHit_OK(t *testing.T) {
	t.Parallel()

	s := newDependencies(t)
	defer finishDependencies(s)

	accID := uuid.New()

	filter := usecase.CountReportItemsQueryFilter{
		AccountID: accID,
	}

	want := []*entity.ReportItem{
		{
			AccountID: accID,
			DateFrom:  civil.Date{Year: 2025, Month: 12, Day: 1},
			DateTo:    civil.Date{Year: 2025, Month: 12, Day: 31},
			Items:     []*entity.AccountTransactionReportItem{},
		},
	}

	s.transactionCacheRepo.GetReportsMock.Set(func(ctx context.Context, key string) ([]*entity.ReportItem, error) {
		require.NotEmpty(t, key)
		return want, nil
	})

	items, hit, err := s.uc.CountReportItems(testCtx(), filter)
	require.NoError(t, err)
	require.True(t, hit)
	require.Equal(t, want, items)
}
