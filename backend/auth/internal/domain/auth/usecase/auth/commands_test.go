package auth

import (
	"context"
	"net"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/m11ano/budget_planner/backend/auth/internal/common/uctypes"
	"github.com/m11ano/budget_planner/backend/auth/internal/domain/auth/entity"
	"github.com/m11ano/budget_planner/backend/auth/internal/domain/auth/usecase"
	"github.com/m11ano/budget_planner/backend/auth/pkg/auth"
	"github.com/stretchr/testify/require"
)

func TestAuthUsecase_LoginByEmail_Table(t *testing.T) {
	t.Parallel()

	table := []tableTest{
		{
			name: "Success",
			run: func(t *testing.T, s *dependencies) {
				email := "User1@Example.com"
				password := "pass123456"
				ip := net.ParseIP("127.0.0.1")

				dto := mustAccountDTO(t, "user1@example.com", password, true, false)

				s.accountUCMock.FindOneByEmailMock.Set(func(
					ctx context.Context,
					gotEmail string,
					qp *uctypes.QueryGetOneParams,
				) (*usecase.AccountDTO, error) {
					require.Equal(t, "user1@example.com", gotEmail)
					require.NotNil(t, qp)
					require.True(t, qp.ForUpdate)
					return dto, nil
				})

				s.sessionUCMock.CreateMock.Set(func(ctx context.Context, sess *entity.Session) error {
					require.Equal(t, dto.Account.ID, sess.AccountID)
					return nil
				})

				s.accountUCMock.UpdateAccountMock.Set(func(ctx context.Context, acc *entity.Account) error {
					require.Equal(t, dto.Account.ID, acc.ID)
					return nil
				})

				got, err := s.uc.LoginByEmail(s.ctx, email, password, ip)
				require.NoError(t, err)
				require.NotNil(t, got)
				require.NotNil(t, got.Session)
				require.NotNil(t, got.AccountDTO)
				require.Equal(t, dto.Account.ID, got.AccountDTO.Account.ID)
				require.NotNil(t, got.RefreshClaims)
				require.NotNil(t, got.AccessClaims)
			},
		},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			s := newDependencies(t)
			defer finishDependencies(t, s)
			tt.run(t, s)
		})
	}
}

func TestAuthUsecase_IssueAccessJWT_Table(t *testing.T) {
	t.Parallel()

	table := []tableTest{
		{
			name: "Success",
			run: func(t *testing.T, s *dependencies) {
				now := time.Now()
				claims := &auth.SessionAccessClaims{
					RegisteredClaims: jwt.RegisteredClaims{
						IssuedAt:  jwt.NewNumericDate(now),
						NotBefore: jwt.NewNumericDate(now),
						ExpiresAt: jwt.NewNumericDate(now.Add(10 * time.Minute)),
					},
				}

				token, err := s.uc.IssueAccessJWT(claims)
				require.NoError(t, err)
				require.NotEmpty(t, token)
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

func TestAuthUsecase_IssueRefreshJWT_And_ParseRefreshToken_Table(t *testing.T) {
	t.Parallel()

	table := []tableTest{
		{
			name: "Success_roundtrip",
			run: func(t *testing.T, s *dependencies) {
				now := time.Now()
				claims := &entity.SessionRefreshClaims{
					RegisteredClaims: jwt.RegisteredClaims{
						IssuedAt:  jwt.NewNumericDate(now),
						NotBefore: jwt.NewNumericDate(now),
						ExpiresAt: jwt.NewNumericDate(now.Add(10 * time.Minute)),
					},
				}
				claims.SessionID = uuid.New()
				claims.RefreshKey = uuid.New()
				claims.RefreshVersion = 1

				token, err := s.uc.IssueRefreshJWT(claims)
				require.NoError(t, err)
				require.NotEmpty(t, token)

				got, err := s.uc.ParseRefreshToken(token, true)
				require.NoError(t, err)
				require.NotNil(t, got)
				require.Equal(t, claims.SessionID, got.SessionID)
				require.Equal(t, claims.RefreshKey, got.RefreshKey)
				require.Equal(t, claims.RefreshVersion, got.RefreshVersion)
			},
		},
		{
			name: "Negative_invalid_token",
			run: func(t *testing.T, s *dependencies) {
				_, err := s.uc.ParseRefreshToken("not-a-token", true)
				require.Error(t, err)
				require.ErrorIs(t, err, usecase.ErrInvalidToken)
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

func TestAuthUsecase_GenerateNewClaims_Table(t *testing.T) {
	t.Parallel()

	table := []tableTest{
		{
			name: "Success",
			run: func(t *testing.T, s *dependencies) {
				ip := net.ParseIP("127.0.0.1")
				now := time.Now()

				accDTO := mustAccountDTO(t, "user1@example.com", "pass123456", true, false)
				accID := accDTO.Account.ID

				session := entity.NewSession(
					accID,
					ip,
					now.Add(-time.Minute),
					time.Duration(s.refreshLifetimeHr)*time.Hour,
				)

				refresh := &entity.SessionRefreshClaims{
					SessionID:        session.ID,
					RefreshKey:       session.RefreshToken,
					RefreshVersion:   session.RefreshVersion,
					RegisteredClaims: jwt.RegisteredClaims{},
				}

				s.sessionUCMock.FindOneByIDMock.Set(func(
					ctx context.Context,
					id uuid.UUID,
					qp *uctypes.QueryGetOneParams,
				) (*entity.Session, error) {
					require.Equal(t, refresh.SessionID, id)
					require.NotNil(t, qp)
					require.True(t, qp.ForUpdate)
					return session, nil
				})

				s.accountUCMock.FindOneByIDMock.Set(func(
					ctx context.Context,
					id uuid.UUID,
					qp *uctypes.QueryGetOneParams,
				) (*usecase.AccountDTO, error) {
					require.Equal(t, accID, id)
					require.NotNil(t, qp)
					require.True(t, qp.ForUpdate)
					return accDTO, nil
				})

				s.sessionUCMock.UpdateMock.Set(func(ctx context.Context, got *entity.Session) error {
					require.Equal(t, session.ID, got.ID)
					return nil
				})

				s.accountUCMock.UpdateAccountMock.Set(func(ctx context.Context, got *entity.Account) error {
					require.Equal(t, accID, got.ID)
					return nil
				})

				out, err := s.uc.GenerateNewClaims(s.ctx, refresh, ip)
				require.NoError(t, err)
				require.NotNil(t, out)
				require.NotNil(t, out.Session)
				require.NotNil(t, out.RefreshClaims)
				require.NotNil(t, out.AccessClaims)
				require.NotNil(t, out.AccountDTO)
				require.Equal(t, accID, out.AccountDTO.Account.ID)
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
