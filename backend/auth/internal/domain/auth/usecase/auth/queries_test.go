package auth

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

func TestAuthUsecase_IsSessionConfirmed_Table(t *testing.T) {
	t.Parallel()

	table := []tableTest{
		{
			name: "Success_true",
			run: func(t *testing.T, s *dependencies) {
				sessionID := uuid.New()
				accDTO := mustAccountDTO(t, "user1@example.com", "pass123456", true, false)

				s.sessionUCMock.FindOneByIDMock.Set(func(ctx context.Context, id uuid.UUID, qp *uctypes.QueryGetOneParams) (*entity.Session, error) {
					require.Equal(t, sessionID, id)
					require.Nil(t, qp)
					return &entity.Session{
						ID:        sessionID,
						AccountID: accDTO.Account.ID,
					}, nil
				})

				s.accountUCMock.FindOneByIDMock.Set(func(ctx context.Context, id uuid.UUID, qp *uctypes.QueryGetOneParams) (*usecase.AccountDTO, error) {
					require.Equal(t, accDTO.Account.ID, id)
					require.Nil(t, qp)
					return accDTO, nil
				})

				ok, err := s.uc.IsSessionConfirmed(s.ctx, sessionID)
				require.NoError(t, err)
				require.True(t, ok)
			},
		},
		{
			name: "Negative_session_not_found_returns_false_nil",
			run: func(t *testing.T, s *dependencies) {
				sessionID := uuid.New()

				s.sessionUCMock.FindOneByIDMock.Set(func(ctx context.Context, id uuid.UUID, qp *uctypes.QueryGetOneParams) (*entity.Session, error) {
					return nil, appErrors.ErrNotFound
				})

				ok, err := s.uc.IsSessionConfirmed(s.ctx, sessionID)
				require.NoError(t, err)
				require.False(t, ok)
			},
		},
		{
			name: "Negative_account_blocked_returns_false_nil",
			run: func(t *testing.T, s *dependencies) {
				sessionID := uuid.New()
				accDTO := mustAccountDTO(t, "user2@example.com", "pass123456", true, true)

				s.sessionUCMock.FindOneByIDMock.Set(func(ctx context.Context, id uuid.UUID, qp *uctypes.QueryGetOneParams) (*entity.Session, error) {
					return &entity.Session{
						ID:        sessionID,
						AccountID: accDTO.Account.ID,
					}, nil
				})

				s.accountUCMock.FindOneByIDMock.Set(func(ctx context.Context, id uuid.UUID, qp *uctypes.QueryGetOneParams) (*usecase.AccountDTO, error) {
					return accDTO, nil
				})

				ok, err := s.uc.IsSessionConfirmed(s.ctx, sessionID)
				require.NoError(t, err)
				require.False(t, ok)
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
