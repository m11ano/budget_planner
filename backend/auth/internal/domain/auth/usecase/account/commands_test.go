package account

import (
	"context"
	"net"
	"testing"

	"github.com/gojuno/minimock/v3"
	"github.com/google/uuid"
	appErrors "github.com/m11ano/budget_planner/backend/auth/internal/app/errors"
	"github.com/m11ano/budget_planner/backend/auth/internal/common/uctypes"
	"github.com/m11ano/budget_planner/backend/auth/internal/domain/auth/entity"
	"github.com/m11ano/budget_planner/backend/auth/internal/domain/auth/usecase"
	"github.com/stretchr/testify/require"
)

func mustAccount(t *testing.T, email, password string) *entity.Account {
	t.Helper()
	a, err := entity.NewAccount(email, password, true, true)
	require.NoError(t, err)
	return a
}

func TestAccountUsecase_CreateAccountByDTO_Table(t *testing.T) {
	t.Parallel()

	ip := net.ParseIP("127.0.0.1")

	table := []tableTest{
		{
			name: "Success",
			run: func(t *testing.T, s *dependencies) {
				in := usecase.CreateAccountDataInput{
					Email:             "user1@example.com",
					Password:          "pass123456",
					SkipPasswordCheck: true,
					IsConfirmed:       true,
					ProfileName:       "Ivan",
					ProfileSurname:    "Petrov",
				}

				s.accountRepoMock.FindOneByEmailMock.Set(func(
					ctx context.Context,
					email string,
					qp *uctypes.QueryGetOneParams,
				) (*entity.Account, error) {
					require.Equal(t, in.Email, email)
					require.Nil(t, qp)
					return nil, appErrors.ErrNotFound
				})

				s.accountRepoMock.CreateMock.Set(func(ctx context.Context, acc *entity.Account) error {
					require.Equal(t, in.Email, acc.Email)
					return nil
				})

				got, err := s.uc.CreateAccountByDTO(s.ctx, in, &ip)
				require.NoError(t, err)
				require.NotNil(t, got)
				require.NotNil(t, got.Account)
				require.Equal(t, in.Email, got.Account.Email)
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

func TestAccountUsecase_PatchAccountByDTO_Table(t *testing.T) {
	t.Parallel()

	table := []tableTest{
		{
			name: "Success_password_changed_revokes_sessions",
			run: func(t *testing.T, s *dependencies) {
				id := uuid.New()

				acc := mustAccount(t, "user2@example.com", "oldpass123")
				acc.ID = id

				newPass := "newpass123456"
				in := usecase.PatchAccountDataInput{
					Version:           acc.Version(),
					Password:          &newPass,
					SkipPasswordCheck: true,
				}

				s.accountRepoMock.FindOneByIDMock.Set(func(
					ctx context.Context,
					gotID uuid.UUID,
					qp *uctypes.QueryGetOneParams,
				) (*entity.Account, error) {
					require.Equal(t, id, gotID)
					require.NotNil(t, qp)
					require.True(t, qp.ForUpdate)
					return acc, nil
				})

				s.accountRepoMock.UpdateMock.Set(func(ctx context.Context, got *entity.Account) error {
					require.Equal(t, id, got.ID)
					return nil
				})

				s.sessionUCMock.RevokeSessionsByAccountIDMock.
					Expect(minimock.AnyContext, id).
					Return(nil)

				err := s.uc.PatchAccountByDTO(s.ctx, id, in, false)
				require.NoError(t, err)
			},
		},
		{
			name: "Success_email_changed_same_account_ok",
			run: func(t *testing.T, s *dependencies) {
				id := uuid.New()

				acc := mustAccount(t, "user5@example.com", "oldpass123")
				acc.ID = id

				newEmail := "user5_new@example.com"
				in := usecase.PatchAccountDataInput{
					Version: acc.Version(),
					Email:   &newEmail,
				}

				s.accountRepoMock.FindOneByIDMock.Set(func(
					ctx context.Context,
					gotID uuid.UUID,
					qp *uctypes.QueryGetOneParams,
				) (*entity.Account, error) {
					require.Equal(t, id, gotID)
					require.NotNil(t, qp)
					require.True(t, qp.ForUpdate)
					return acc, nil
				})

				s.accountRepoMock.FindOneByEmailMock.Set(func(
					ctx context.Context,
					email string,
					qp *uctypes.QueryGetOneParams,
				) (*entity.Account, error) {
					require.Equal(t, newEmail, email)
					require.Nil(t, qp)
					return &entity.Account{ID: id, Email: newEmail}, nil
				})

				s.accountRepoMock.UpdateMock.Set(func(ctx context.Context, got *entity.Account) error {
					require.Equal(t, id, got.ID)
					require.Equal(t, newEmail, got.Email)
					return nil
				})

				err := s.uc.PatchAccountByDTO(s.ctx, id, in, false)
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

func TestAccountUsecase_UpdateAccount_Table(t *testing.T) {
	t.Parallel()

	table := []tableTest{
		{
			name: "Success",
			run: func(t *testing.T, s *dependencies) {
				item := mustAccount(t, "u6@example.com", "pass123456")

				s.accountRepoMock.UpdateMock.Set(func(ctx context.Context, got *entity.Account) error {
					require.Equal(t, item, got)
					return nil
				})

				err := s.uc.UpdateAccount(s.ctx, item)
				require.NoError(t, err)
			},
		},
		{
			name: "Error_repo",
			run: func(t *testing.T, s *dependencies) {
				item := mustAccount(t, "u7@example.com", "pass123456")

				s.accountRepoMock.UpdateMock.Set(func(ctx context.Context, got *entity.Account) error {
					return appErrors.ErrInternal
				})

				err := s.uc.UpdateAccount(s.ctx, item)
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
