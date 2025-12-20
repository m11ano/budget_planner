package budget

import (
	"context"
	"errors"
	"testing"
	"time"

	"cloud.google.com/go/civil"
	"github.com/google/uuid"
	"github.com/govalues/decimal"
	appErrors "github.com/m11ano/budget_planner/backend/ledger/internal/app/errors"
	"github.com/m11ano/budget_planner/backend/ledger/internal/common/uctypes"
	"github.com/m11ano/budget_planner/backend/ledger/internal/domain/budget/entity"
	"github.com/m11ano/budget_planner/backend/ledger/internal/domain/budget/usecase"
	"github.com/stretchr/testify/require"
)

func TestBudgetUsecase_FindOneByID_Table(t *testing.T) {
	t.Parallel()

	type want struct {
		dto     *usecase.BudgetDTO
		hit     bool
		wantErr error
	}

	id := uuid.New()
	accID := uuid.New()

	item := &entity.Budget{
		ID:         id,
		AccountID:  accID,
		Period:     civil.Date{Year: 2025, Month: 12, Day: 1},
		CategoryID: 10,
		Amount:     decimal.MustParse("123.45"),
	}

	tests := []struct {
		name string

		queryParams *uctypes.QueryGetOneParams
		cacheHas    bool
		repoErr     error

		want want
	}{
		{
			name:        "Cache_hit",
			queryParams: nil,
			cacheHas:    true,
			want: want{
				dto: &usecase.BudgetDTO{Budget: item},
				hit: true,
			},
		},
		{
			name:        "Cache_miss_repo_ok",
			queryParams: nil,
			cacheHas:    false,
			want: want{
				dto: &usecase.BudgetDTO{Budget: item},
				hit: false,
			},
		},
		{
			name:        "SkipCache_true_repo_ok",
			queryParams: &uctypes.QueryGetOneParams{SkipCache: true},
			cacheHas:    false,
			want: want{
				dto: &usecase.BudgetDTO{Budget: item},
				hit: false,
			},
		},
		{
			name:        "Repo_error_simple_negative",
			queryParams: &uctypes.QueryGetOneParams{SkipCache: true},
			repoErr:     appErrors.ErrNotFound,
			want: want{
				wantErr: appErrors.ErrNotFound,
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			s := newDependencies(t)
			defer finishDependencies(s)

			key := buildKeyForFindOneByID(id)

			if tt.queryParams == nil || !tt.queryParams.SkipCache {
				if tt.cacheHas {
					s.budgetCacheRepo.GetBudgetMock.Set(func(ctx context.Context, gotKey string) (*entity.Budget, error) {
						require.Equal(t, key, gotKey)
						return item, nil
					})
				} else {
					s.budgetCacheRepo.GetBudgetMock.Set(func(ctx context.Context, gotKey string) (*entity.Budget, error) {
						require.Equal(t, key, gotKey)
						return nil, appErrors.ErrNotFound
					})
				}
			}

			needRepo := (tt.queryParams != nil && tt.queryParams.SkipCache) || (!tt.cacheHas)
			if needRepo {
				s.budgetRepo.FindOneByIDMock.Set(func(ctx context.Context, gotID uuid.UUID, qp *uctypes.QueryGetOneParams) (*entity.Budget, error) {
					require.Equal(t, id, gotID)
					if tt.queryParams == nil {
						require.Nil(t, qp)
					} else {
						require.NotNil(t, qp)
						require.Equal(t, tt.queryParams.SkipCache, qp.SkipCache)
					}

					if tt.repoErr != nil {
						return nil, tt.repoErr
					}
					return item, nil
				})

				if tt.repoErr == nil {
					s.budgetCacheRepo.SaveBudgetMock.Set(func(ctx context.Context, gotKey string, gotItem *entity.Budget, _ *time.Duration) error {
						require.Equal(t, key, gotKey)
						require.Equal(t, item, gotItem)
						return nil
					})
				}
			}

			got, hit, err := s.uc.FindOneByID(testCtx(), id, tt.queryParams)

			if tt.want.wantErr != nil {
				require.Error(t, err)
				require.True(t, errors.Is(err, tt.want.wantErr))
				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.want.hit, hit)
			require.NotNil(t, got)
			require.NotNil(t, got.Budget)
			require.Equal(t, item, got.Budget)
		})
	}
}

func TestBudgetUsecase_FindList_Table(t *testing.T) {
	t.Parallel()

	accID := uuid.New()
	listOptions := &usecase.BudgetListOptions{
		FilterAccountID: &accID,
	}

	items := []*entity.Budget{
		{ID: uuid.New(), AccountID: accID, CategoryID: 1, Period: civil.Date{Year: 2025, Month: 12, Day: 1}, Amount: decimal.MustParse("10.00")},
		{ID: uuid.New(), AccountID: accID, CategoryID: 2, Period: civil.Date{Year: 2025, Month: 12, Day: 1}, Amount: decimal.MustParse("20.00")},
	}

	tests := []struct {
		name string

		qp        *uctypes.QueryGetListParams
		cacheHas  bool
		repoCalls bool
		wantHit   bool
	}{
		{name: "Cache_hit", qp: &uctypes.QueryGetListParams{Limit: 50}, cacheHas: true, repoCalls: false, wantHit: true},
		{name: "Cache_miss_repo_ok", qp: &uctypes.QueryGetListParams{Limit: 50}, cacheHas: false, repoCalls: true, wantHit: false},
		{name: "SkipCache_true_repo_ok", qp: &uctypes.QueryGetListParams{Limit: 50, SkipCache: true}, cacheHas: false, repoCalls: true, wantHit: false},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			s := newDependencies(t)
			defer finishDependencies(s)

			key := buildKeyForFindList(listOptions, tt.qp)

			// cache
			if !tt.qp.SkipCache {
				if tt.cacheHas {
					s.budgetCacheRepo.GetBudgetsListMock.Set(func(ctx context.Context, gotKey string) ([]*entity.Budget, error) {
						require.Equal(t, key, gotKey)
						return items, nil
					})
				} else {
					s.budgetCacheRepo.GetBudgetsListMock.Set(func(ctx context.Context, gotKey string) ([]*entity.Budget, error) {
						require.Equal(t, key, gotKey)
						return nil, appErrors.ErrNotFound
					})
				}
			}

			if tt.repoCalls {
				s.budgetRepo.FindListMock.Set(func(ctx context.Context, gotOpt *usecase.BudgetListOptions, gotQP *uctypes.QueryGetListParams) ([]*entity.Budget, error) {
					require.Equal(t, listOptions, gotOpt)
					require.Equal(t, tt.qp, gotQP)
					return items, nil
				})

				s.budgetCacheRepo.SaveBudgetsListMock.Set(func(ctx context.Context, gotKey string, gotItems []*entity.Budget, _ *time.Duration) error {
					require.Equal(t, key, gotKey)
					require.Equal(t, items, gotItems)
					return nil
				})
			}

			got, hit, err := s.uc.FindList(testCtx(), listOptions, tt.qp)
			require.NoError(t, err)
			require.Equal(t, tt.wantHit, hit)
			require.Len(t, got, len(items))
			require.Equal(t, items[0], got[0].Budget)
			require.Equal(t, items[1], got[1].Budget)
		})
	}
}

func TestBudgetUsecase_FindPagedList_Table(t *testing.T) {
	t.Parallel()

	accID := uuid.New()
	listOptions := &usecase.BudgetListOptions{
		FilterAccountID: &accID,
	}

	items := []*entity.Budget{
		{ID: uuid.New(), AccountID: accID, CategoryID: 1, Period: civil.Date{Year: 2025, Month: 11, Day: 1}, Amount: decimal.MustParse("10.00")},
	}
	total := uint64(42)

	tests := []struct {
		name string

		qp       *uctypes.QueryGetListParams
		cacheHas bool
		wantHit  bool
	}{
		{name: "Cache_hit", qp: &uctypes.QueryGetListParams{Limit: 10, Offset: 2}, cacheHas: true, wantHit: true},
		{name: "Cache_miss_repo_ok", qp: &uctypes.QueryGetListParams{Limit: 10, Offset: 2}, cacheHas: false, wantHit: false},
		{name: "SkipCache_true_repo_ok", qp: &uctypes.QueryGetListParams{Limit: 10, Offset: 2, SkipCache: true}, cacheHas: false, wantHit: false},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			s := newDependencies(t)
			defer finishDependencies(s)

			key := buildKeyForFindPagedList(listOptions, tt.qp)

			if !tt.qp.SkipCache {
				if tt.cacheHas {
					s.budgetCacheRepo.GetBudgetsPagedListMock.Set(func(ctx context.Context, gotKey string) ([]*entity.Budget, uint64, error) {
						require.Equal(t, key, gotKey)
						return items, total, nil
					})
				} else {
					s.budgetCacheRepo.GetBudgetsPagedListMock.Set(func(ctx context.Context, gotKey string) ([]*entity.Budget, uint64, error) {
						require.Equal(t, key, gotKey)
						return nil, 0, appErrors.ErrNotFound
					})
				}
			}

			needRepo := tt.qp.SkipCache || !tt.cacheHas
			if needRepo {
				s.budgetRepo.FindPagedListMock.Set(func(ctx context.Context, gotOpt *usecase.BudgetListOptions, gotQP *uctypes.QueryGetListParams) ([]*entity.Budget, uint64, error) {
					require.Equal(t, listOptions, gotOpt)
					require.Equal(t, tt.qp, gotQP)
					return items, total, nil
				})

				s.budgetCacheRepo.SaveBudgetsPagedListMock.Set(func(ctx context.Context, gotKey string, gotItems []*entity.Budget, gotTotal uint64, _ *time.Duration) error {
					require.Equal(t, key, gotKey)
					require.Equal(t, items, gotItems)
					require.Equal(t, total, gotTotal)
					return nil
				})
			}

			got, gotTotal, hit, err := s.uc.FindPagedList(testCtx(), listOptions, tt.qp)
			require.NoError(t, err)
			require.Equal(t, tt.wantHit, hit)
			require.Equal(t, total, gotTotal)
			require.Len(t, got, len(items))
			require.Equal(t, items[0], got[0].Budget)
		})
	}
}
