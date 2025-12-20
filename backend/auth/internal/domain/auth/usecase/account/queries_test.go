package account

import (
	"context"
	"testing"

	"github.com/google/uuid"
	appErrors "github.com/m11ano/budget_planner/backend/auth/internal/app/errors"
	"github.com/m11ano/budget_planner/backend/auth/internal/common/uctypes"
	"github.com/m11ano/budget_planner/backend/auth/internal/domain/auth/entity"
	"github.com/m11ano/budget_planner/backend/auth/internal/domain/auth/usecase"
	"github.com/stretchr/testify/require"
)

func TestAccountUsecase_FindOneByEmail_Table(t *testing.T) {
	t.Parallel()

	table := []tableTest{
		{
			name: "Success",
			run: func(t *testing.T, s *dependencies) {
				email := "user@example.com"
				item := newAccount(uuid.New(), email)

				qp := &uctypes.QueryGetOneParams{ForUpdate: false}

				s.accountRepoMock.FindOneByEmailMock.Set(func(ctx context.Context, gotEmail string, gotQP *uctypes.QueryGetOneParams) (*entity.Account, error) {
					require.Equal(t, email, gotEmail)
					require.Equal(t, qp, gotQP)
					return item, nil
				})

				got, err := s.uc.FindOneByEmail(s.ctx, email, qp)
				require.NoError(t, err)
				require.NotNil(t, got)
				require.NotNil(t, got.Account)
				require.Equal(t, item.ID, got.Account.ID)
				require.Equal(t, item.Email, got.Account.Email)
			},
		},
		{
			name: "Error_repo",
			run: func(t *testing.T, s *dependencies) {
				email := "user@example.com"

				s.accountRepoMock.FindOneByEmailMock.Set(func(ctx context.Context, gotEmail string, gotQP *uctypes.QueryGetOneParams) (*entity.Account, error) {
					require.Equal(t, email, gotEmail)
					return nil, appErrors.ErrInternal
				})

				_, err := s.uc.FindOneByEmail(s.ctx, email, nil)
				require.Error(t, err)
				require.ErrorIs(t, err, appErrors.ErrInternal)
			},
		},
	}

	for _, tt := range table {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			s := newDependencies(t)
			defer finishDependencies(t, s)
			tt.run(t, s)
		})
	}
}

func TestAccountUsecase_FindOneByID_Table(t *testing.T) {
	t.Parallel()

	table := []tableTest{
		{
			name: "Success",
			run: func(t *testing.T, s *dependencies) {
				id := uuid.New()
				item := newAccount(id, "user@example.com")

				qp := &uctypes.QueryGetOneParams{ForUpdate: true}

				s.accountRepoMock.FindOneByIDMock.Set(func(ctx context.Context, gotID uuid.UUID, gotQP *uctypes.QueryGetOneParams) (*entity.Account, error) {
					require.Equal(t, id, gotID)
					require.Equal(t, qp, gotQP)
					return item, nil
				})

				got, err := s.uc.FindOneByID(s.ctx, id, qp)
				require.NoError(t, err)
				require.NotNil(t, got)
				require.NotNil(t, got.Account)
				require.Equal(t, item.ID, got.Account.ID)
				require.Equal(t, item.Email, got.Account.Email)
			},
		},
		{
			name: "Error_repo",
			run: func(t *testing.T, s *dependencies) {
				id := uuid.New()

				s.accountRepoMock.FindOneByIDMock.Set(func(ctx context.Context, gotID uuid.UUID, gotQP *uctypes.QueryGetOneParams) (*entity.Account, error) {
					require.Equal(t, id, gotID)
					return nil, appErrors.ErrNotFound
				})

				_, err := s.uc.FindOneByID(s.ctx, id, nil)
				require.Error(t, err)
				require.ErrorIs(t, err, appErrors.ErrNotFound)
			},
		},
	}

	for _, tt := range table {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			s := newDependencies(t)
			defer finishDependencies(t, s)
			tt.run(t, s)
		})
	}
}

func TestAccountUsecase_FindList_Table(t *testing.T) {
	t.Parallel()

	table := []tableTest{
		{
			name: "Success_empty",
			run: func(t *testing.T, s *dependencies) {
				opts := &usecase.AccountListOptions{}
				qp := &uctypes.QueryGetListParams{}

				s.accountRepoMock.FindListMock.Set(func(ctx context.Context, gotOpts *usecase.AccountListOptions, gotQP *uctypes.QueryGetListParams) ([]*entity.Account, error) {
					require.Equal(t, opts, gotOpts)
					require.Equal(t, qp, gotQP)
					return []*entity.Account{}, nil
				})

				got, err := s.uc.FindList(s.ctx, opts, qp)
				require.NoError(t, err)
				require.Len(t, got, 0)
			},
		},
		{
			name: "Success_two_items",
			run: func(t *testing.T, s *dependencies) {
				opts := &usecase.AccountListOptions{}
				qp := &uctypes.QueryGetListParams{}

				a1 := newAccount(uuid.New(), "a1@example.com")
				a2 := newAccount(uuid.New(), "a2@example.com")

				s.accountRepoMock.FindListMock.Set(func(ctx context.Context, gotOpts *usecase.AccountListOptions, gotQP *uctypes.QueryGetListParams) ([]*entity.Account, error) {
					require.Equal(t, opts, gotOpts)
					require.Equal(t, qp, gotQP)
					return []*entity.Account{a1, a2}, nil
				})

				got, err := s.uc.FindList(s.ctx, opts, qp)
				require.NoError(t, err)
				require.Len(t, got, 2)
				require.Equal(t, a1.ID, got[0].Account.ID)
				require.Equal(t, a2.ID, got[1].Account.ID)
			},
		},
		{
			name: "Error_repo",
			run: func(t *testing.T, s *dependencies) {
				opts := &usecase.AccountListOptions{}

				s.accountRepoMock.FindListMock.Set(func(ctx context.Context, gotOpts *usecase.AccountListOptions, gotQP *uctypes.QueryGetListParams) ([]*entity.Account, error) {
					return nil, appErrors.ErrInternal
				})

				_, err := s.uc.FindList(s.ctx, opts, nil)
				require.Error(t, err)
				require.ErrorIs(t, err, appErrors.ErrInternal)
			},
		},
	}

	for _, tt := range table {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			s := newDependencies(t)
			defer finishDependencies(t, s)
			tt.run(t, s)
		})
	}
}

func TestAccountUsecase_FindPagedList_Table(t *testing.T) {
	t.Parallel()

	table := []tableTest{
		{
			name: "Success",
			run: func(t *testing.T, s *dependencies) {
				opts := &usecase.AccountListOptions{}
				qp := &uctypes.QueryGetListParams{}

				a1 := newAccount(uuid.New(), "a1@example.com")
				a2 := newAccount(uuid.New(), "a2@example.com")

				s.accountRepoMock.FindPagedListMock.Set(func(ctx context.Context, gotOpts *usecase.AccountListOptions, gotQP *uctypes.QueryGetListParams) ([]*entity.Account, uint64, error) {
					require.Equal(t, opts, gotOpts)
					require.Equal(t, qp, gotQP)
					return []*entity.Account{a1, a2}, 2, nil
				})

				got, total, err := s.uc.FindPagedList(s.ctx, opts, qp)
				require.NoError(t, err)
				require.Equal(t, uint64(2), total)
				require.Len(t, got, 2)
				require.Equal(t, a1.ID, got[0].Account.ID)
				require.Equal(t, a2.ID, got[1].Account.ID)
			},
		},
		{
			name: "Error_repo",
			run: func(t *testing.T, s *dependencies) {
				opts := &usecase.AccountListOptions{}

				s.accountRepoMock.FindPagedListMock.Set(func(ctx context.Context, gotOpts *usecase.AccountListOptions, gotQP *uctypes.QueryGetListParams) ([]*entity.Account, uint64, error) {
					return nil, 0, appErrors.ErrInternal
				})

				_, _, err := s.uc.FindPagedList(s.ctx, opts, nil)
				require.Error(t, err)
				require.ErrorIs(t, err, appErrors.ErrInternal)
			},
		},
	}

	for _, tt := range table {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			s := newDependencies(t)
			defer finishDependencies(t, s)
			tt.run(t, s)
		})
	}
}

func TestAccountUsecase_FindListInMap_Table(t *testing.T) {
	t.Parallel()

	table := []tableTest{
		{
			name: "Success_empty",
			run: func(t *testing.T, s *dependencies) {
				opts := &usecase.AccountListOptions{}

				s.accountRepoMock.FindListMock.Set(func(ctx context.Context, gotOpts *usecase.AccountListOptions, gotQP *uctypes.QueryGetListParams) ([]*entity.Account, error) {
					require.Equal(t, opts, gotOpts)
					return []*entity.Account{}, nil
				})

				got, err := s.uc.FindListInMap(s.ctx, opts, nil)
				require.NoError(t, err)
				require.Len(t, got, 0)
			},
		},
		{
			name: "Success_two_items",
			run: func(t *testing.T, s *dependencies) {
				opts := &usecase.AccountListOptions{}

				a1 := newAccount(uuid.New(), "a1@example.com")
				a2 := newAccount(uuid.New(), "a2@example.com")

				s.accountRepoMock.FindListMock.Set(func(ctx context.Context, gotOpts *usecase.AccountListOptions, gotQP *uctypes.QueryGetListParams) ([]*entity.Account, error) {
					require.Equal(t, opts, gotOpts)
					return []*entity.Account{a1, a2}, nil
				})

				got, err := s.uc.FindListInMap(s.ctx, opts, nil)
				require.NoError(t, err)
				require.Len(t, got, 2)

				require.NotNil(t, got[a1.ID])
				require.Equal(t, a1.ID, got[a1.ID].Account.ID)

				require.NotNil(t, got[a2.ID])
				require.Equal(t, a2.ID, got[a2.ID].Account.ID)
			},
		},
		{
			name: "Error_repo",
			run: func(t *testing.T, s *dependencies) {
				opts := &usecase.AccountListOptions{}

				s.accountRepoMock.FindListMock.Set(func(ctx context.Context, gotOpts *usecase.AccountListOptions, gotQP *uctypes.QueryGetListParams) ([]*entity.Account, error) {
					return nil, appErrors.ErrInternal
				})

				_, err := s.uc.FindListInMap(s.ctx, opts, nil)
				require.Error(t, err)
				require.ErrorIs(t, err, appErrors.ErrInternal)
			},
		},
	}

	for _, tt := range table {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			s := newDependencies(t)
			defer finishDependencies(t, s)
			tt.run(t, s)
		})
	}
}
