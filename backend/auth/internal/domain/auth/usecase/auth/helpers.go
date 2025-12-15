package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/m11ano/budget_planner/backend/auth/internal/domain/auth/entity"
	"github.com/m11ano/budget_planner/backend/auth/pkg/auth"
	authpb "github.com/m11ano/budget_planner/backend/auth/pkg/proto_pb/claims"
)

func (uc *UsecaseImpl) makeRefreshClaims(session *entity.Session) *entity.SessionRefreshClaims {
	return &entity.SessionRefreshClaims{
		SessionID:      session.ID,
		RefreshKey:     session.RefreshToken,
		RefreshVersion: session.RefreshVersion,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(session.RefreshTokenIssuedAt),
			ExpiresAt: jwt.NewNumericDate(session.RefreshTokenExpiresAt),
		},
	}
}

func (uc *UsecaseImpl) makeAccessClaims(
	sessionID uuid.UUID,
	accountID uuid.UUID,
	meta map[string]string,
	issuedAt time.Time,
	duration time.Duration,
) *auth.SessionAccessClaims {
	return &auth.SessionAccessClaims{
		AccessClaims: authpb.AccessClaims{
			AccountId: accountID.String(),
			SessionId: sessionID.String(),
			Meta:      meta,
		},
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(issuedAt),
			ExpiresAt: jwt.NewNumericDate(issuedAt.Add(duration)),
		},
	}
}
