package session

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

func TestSessionUsecase_FindOneByID_Table(t *testing.T) {
	t.Parallel()

	table := []tableTest{
		{
			name: "Success",
			run: func(t *testing.T, s *dependencies) {
				id := uuid.New()
				qp := &uctypes.QueryGetOneParams{}
				want := newSession(id, uuid.New())

				s.repoMock.FindOneByIDMock.Set(func(
					ctx context.Context,
					gotID uuid.UUID,
					gotQP *uctypes.QueryGetOneParams,
				) (*entity.Session, error) {
					require.Equal(t, id, gotID)
					require.Equal(t, qp, gotQP)
					return want, nil
				})

				got, err := s.uc.FindOneByID(s.ctx, id, qp)
				require.NoError(t, err)
				require.Equal(t, want, got)
			},
		},
		{
			name: "Error_repo",
			run: func(t *testing.T, s *dependencies) {
				id := uuid.New()
				qp := &uctypes.QueryGetOneParams{}

				s.repoMock.FindOneByIDMock.Set(func(
					ctx context.Context,
					gotID uuid.UUID,
					gotQP *uctypes.QueryGetOneParams,
				) (*entity.Session, error) {
					require.Equal(t, id, gotID)
					require.Equal(t, qp, gotQP)
					return nil, appErrors.ErrInternal
				})

				got, err := s.uc.FindOneByID(s.ctx, id, qp)
				require.Nil(t, got)
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

func TestSessionUsecase_FindList_Table(t *testing.T) {
	t.Parallel()

	table := []tableTest{
		{
			name: "Success",
			run: func(t *testing.T, s *dependencies) {
				accountID := uuid.New()
				listOpts := &usecase.SessionListOptions{FilterAccountID: &accountID}
				qp := &uctypes.QueryGetListParams{}

				want := []*entity.Session{
					newSession(uuid.New(), accountID),
					newSession(uuid.New(), accountID),
				}

				s.repoMock.FindListMock.Set(func(
					ctx context.Context,
					gotOpts *usecase.SessionListOptions,
					gotQP *uctypes.QueryGetListParams,
				) ([]*entity.Session, error) {
					require.Equal(t, listOpts, gotOpts)
					require.Equal(t, qp, gotQP)
					return want, nil
				})

				got, err := s.uc.FindList(s.ctx, listOpts, qp)
				require.NoError(t, err)
				require.Equal(t, want, got)
			},
		},
		{
			name: "Success_empty",
			run: func(t *testing.T, s *dependencies) {
				listOpts := &usecase.SessionListOptions{}
				qp := &uctypes.QueryGetListParams{}

				s.repoMock.FindListMock.Set(func(
					ctx context.Context,
					gotOpts *usecase.SessionListOptions,
					gotQP *uctypes.QueryGetListParams,
				) ([]*entity.Session, error) {
					require.Equal(t, listOpts, gotOpts)
					require.Equal(t, qp, gotQP)
					return []*entity.Session{}, nil
				})

				got, err := s.uc.FindList(s.ctx, listOpts, qp)
				require.NoError(t, err)
				require.Len(t, got, 0)
			},
		},
		{
			name: "Error_repo",
			run: func(t *testing.T, s *dependencies) {
				listOpts := &usecase.SessionListOptions{}
				qp := &uctypes.QueryGetListParams{}

				s.repoMock.FindListMock.Set(func(
					ctx context.Context,
					gotOpts *usecase.SessionListOptions,
					gotQP *uctypes.QueryGetListParams,
				) ([]*entity.Session, error) {
					require.Equal(t, listOpts, gotOpts)
					require.Equal(t, qp, gotQP)
					return nil, appErrors.ErrInternal
				})

				got, err := s.uc.FindList(s.ctx, listOpts, qp)
				require.Nil(t, got)
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
