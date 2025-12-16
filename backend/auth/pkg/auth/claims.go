package auth

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	authpb "github.com/m11ano/budget_planner/backend/auth/pkg/proto_pb/access_claims"
)

type SessionAccessClaims struct {
	authpb.AccessClaims
	jwt.RegisteredClaims
}

func (c *SessionAccessClaims) GetSessionID() (uuid.UUID, error) {
	return uuid.Parse(c.SessionId)
}

func (c *SessionAccessClaims) GetAccountID() (uuid.UUID, error) {
	return uuid.Parse(c.AccountId)
}
