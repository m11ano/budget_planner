package usecase

import (
	"context"
	"net"

	"github.com/google/uuid"
	appErrors "github.com/m11ano/budget_planner/backend/auth/internal/app/errors"
	"github.com/m11ano/budget_planner/backend/auth/internal/domain/auth/entity"
	"github.com/m11ano/budget_planner/backend/auth/pkg/auth"
)

type AuthSessionDTO struct {
	Session       *entity.Session
	RefreshClaims *entity.SessionRefreshClaims
	AccessClaims  *auth.SessionAccessClaims
	AccountDTO    *AccountDTO
}

var ErrPasswordIncorrect = appErrors.ErrUnauthorized.Extend("password incorrect")

var ErrAccountNotConfirmed = appErrors.ErrUnauthorized.Extend("account not confirmed").WithTextCode("ACCOUNT_NOT_CONFIRMED")

var ErrAccountBlocked = appErrors.ErrUnauthorized.Extend("account blocked").WithTextCode("ACCOUNT_BLOCKED")

var ErrAccessDenied = appErrors.ErrUnauthorized.Extend("access denied")

var ErrInvalidToken = appErrors.ErrUnauthorized.Extend("invalid token")

var ErrExpiredToken = appErrors.ErrUnauthorized.Extend("expired token").WithTextCode("EXPIRED_TOKEN")

//go:generate minimock -i github.com/m11ano/budget_planner/backend/auth/internal/domain/auth/usecase.AuthUsecase -o mocks/auth_usecase.go
type AuthUsecase interface {
	LoginByEmail(
		ctx context.Context,
		email string,
		password string,
		ip net.IP,
	) (res *AuthSessionDTO, resErr error)

	IsSessionConfirmed(
		ctx context.Context,
		sessionID uuid.UUID,
	) (res bool, resErr error)

	IssueAccessJWT(access *auth.SessionAccessClaims) (resToken string, resErr error)

	IssueRefreshJWT(refresh *entity.SessionRefreshClaims) (resToken string, resErr error)

	ParseRefreshToken(token string, validate bool) (res *entity.SessionRefreshClaims, resErr error)

	GenerateNewClaims(ctx context.Context, refresh *entity.SessionRefreshClaims, ip net.IP) (res *AuthSessionDTO, resErr error)
}
