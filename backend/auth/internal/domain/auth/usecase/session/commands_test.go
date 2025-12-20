package session

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	appErrors "github.com/m11ano/budget_planner/backend/auth/internal/app/errors"
	"github.com/m11ano/budget_planner/backend/auth/internal/common/uctypes"
	"github.com/m11ano/budget_planner/backend/auth/internal/domain/auth/entity"
	"github.com/m11ano/budget_planner/backend/auth/internal/domain/auth/usecase"
	"github.com/stretchr/testify/require"
)

func newSession(id, accountID uuid.UUID) *entity.Session {
	return &entity.Session{
		ID:        id,
		AccountID: accountID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func TestSessionUsecase_Create_Table(t *testing.T) {
	t.Parallel()

	table := []tableTest{
		{
			name: "Success",
			run: func(t *testing.T, s *dependencies) {
				item := newSession(uuid.New(), uuid.New())

				s.repoMock.CreateMock.Set(func(ctx context.Context, got *entity.Session) error {
					require.Equal(t, item, got)
					return nil
				})

				err := s.uc.Create(s.ctx, item)
				require.NoError(t, err)
			},
		},
		{
			name: "Error_repo",
			run: func(t *testing.T, s *dependencies) {
				item := newSession(uuid.New(), uuid.New())

				s.repoMock.CreateMock.Set(func(ctx context.Context, got *entity.Session) error {
					return appErrors.ErrInternal
				})

				err := s.uc.Create(s.ctx, item)
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

func TestSessionUsecase_Update_Table(t *testing.T) {
	t.Parallel()

	table := []tableTest{
		{
			name: "Success",
			run: func(t *testing.T, s *dependencies) {
				item := newSession(uuid.New(), uuid.New())

				s.repoMock.UpdateMock.Set(func(ctx context.Context, got *entity.Session) error {
					require.Equal(t, item, got)
					return nil
				})

				err := s.uc.Update(s.ctx, item)
				require.NoError(t, err)
			},
		},
		{
			name: "Error_repo",
			run: func(t *testing.T, s *dependencies) {
				item := newSession(uuid.New(), uuid.New())

				s.repoMock.UpdateMock.Set(func(ctx context.Context, got *entity.Session) error {
					return appErrors.ErrInternal
				})

				err := s.uc.Update(s.ctx, item)
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

func TestSessionUsecase_RevokeSessionsByAccountID_Table(t *testing.T) {
	t.Parallel()

	table := []tableTest{
		{
			name: "Success_empty_list",
			run: func(t *testing.T, s *dependencies) {
				accountID := uuid.New()

				s.repoMock.FindListMock.Set(func(
					ctx context.Context,
					opts *usecase.SessionListOptions,
					qp *uctypes.QueryGetListParams,
				) ([]*entity.Session, error) {
					require.NotNil(t, opts)
					require.NotNil(t, opts.FilterAccountID)
					require.Equal(t, accountID, *opts.FilterAccountID)

					require.NotNil(t, qp)
					require.True(t, qp.ForUpdate)

					return []*entity.Session{}, nil
				})

				err := s.uc.RevokeSessionsByAccountID(s.ctx, accountID)
				require.NoError(t, err)
			},
		},
		{
			name: "Success_multiple_sessions",
			run: func(t *testing.T, s *dependencies) {
				accountID := uuid.New()
				s1 := newSession(uuid.New(), accountID)
				s2 := newSession(uuid.New(), accountID)

				s.repoMock.FindListMock.Set(func(
					ctx context.Context,
					opts *usecase.SessionListOptions,
					qp *uctypes.QueryGetListParams,
				) ([]*entity.Session, error) {
					require.Equal(t, accountID, *opts.FilterAccountID)
					require.True(t, qp.ForUpdate)
					return []*entity.Session{s1, s2}, nil
				})

				updated := map[uuid.UUID]bool{}
				s.repoMock.UpdateMock.Set(func(ctx context.Context, got *entity.Session) error {
					updated[got.ID] = true
					return nil
				})

				err := s.uc.RevokeSessionsByAccountID(s.ctx, accountID)
				require.NoError(t, err)
				require.True(t, updated[s1.ID])
				require.True(t, updated[s2.ID])
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

func TestSessionUsecase_RevokeSessionByID_Table(t *testing.T) {
	t.Parallel()

	table := []tableTest{
		{
			name: "Success",
			run: func(t *testing.T, s *dependencies) {
				id := uuid.New()
				item := newSession(id, uuid.New())

				s.repoMock.FindOneByIDMock.Set(func(
					ctx context.Context,
					gotID uuid.UUID,
					qp *uctypes.QueryGetOneParams,
				) (*entity.Session, error) {
					require.Equal(t, id, gotID)
					require.NotNil(t, qp)
					require.True(t, qp.ForUpdate)
					return item, nil
				})

				s.repoMock.UpdateMock.Set(func(ctx context.Context, got *entity.Session) error {
					require.Equal(t, id, got.ID)
					return nil
				})

				err := s.uc.RevokeSessionByID(s.ctx, id)
				require.NoError(t, err)
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
