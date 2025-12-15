package entity

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type SessionRefreshClaims struct {
	SessionID      uuid.UUID `json:"sid,omitempty"`
	RefreshKey     uuid.UUID `json:"rk,omitempty"`
	RefreshVersion uint64    `json:"rv,omitempty"`
	jwt.RegisteredClaims
}
