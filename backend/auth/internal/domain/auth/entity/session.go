package entity

import (
	"net"
	"time"

	"github.com/google/uuid"
)

type Session struct {
	ID                    uuid.UUID
	AccountID             uuid.UUID
	RefreshToken          uuid.UUID
	RefreshVersion        uint64
	RefreshTokenIssuedAt  time.Time
	RefreshTokenExpiresAt time.Time
	RefreshTokenRequestIP net.IP
	CreateRequestIP       net.IP

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

func (item *Session) Version() int64 {
	return item.UpdatedAt.UnixMicro()
}

func (item *Session) GenerateNewRefresh(timeIssuedAt time.Time, duration time.Duration, ip net.IP) {
	item.RefreshTokenIssuedAt = timeIssuedAt.Truncate(time.Microsecond)
	item.RefreshTokenExpiresAt = item.RefreshTokenIssuedAt.Add(duration)
	item.RefreshTokenRequestIP = ip
	item.RefreshToken = uuid.New()
	item.RefreshVersion++
}

func (item *Session) IsAlive(timeForCheck time.Time) bool {
	return item.RefreshTokenExpiresAt.Before(timeForCheck)
}

func NewSession(accountID uuid.UUID, ip net.IP, timeIssuedAt time.Time, refreshDuration time.Duration) *Session {
	timeNow := time.Now().Truncate(time.Microsecond)

	item := &Session{
		ID:              uuid.New(),
		AccountID:       accountID,
		CreateRequestIP: ip,

		CreatedAt: timeNow,
		UpdatedAt: timeNow,
	}

	item.GenerateNewRefresh(timeIssuedAt, refreshDuration, ip)

	return item
}
