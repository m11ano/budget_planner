package pg

import (
	"net"
	"time"

	"github.com/google/uuid"
	"github.com/m11ano/budget_planner/backend/auth/internal/domain/auth/entity"
	"github.com/m11ano/budget_planner/backend/auth/pkg/dbhelper"
)

const (
	SessionTable = "auth_session"
)

var SessionTableFields = []string{}

func init() {
	SessionTableFields = dbhelper.ExtractDBFields(&SessionDBModel{})
}

type SessionDBModel struct {
	ID                    uuid.UUID `db:"id"`
	AccountID             uuid.UUID `db:"account_id"`
	RefreshToken          uuid.UUID `db:"refresh_token"`
	RefreshVersion        uint64    `db:"refresh_version"`
	RefreshTokenIssuedAt  time.Time `db:"refresh_token_issued_at"`
	RefreshTokenExpiresAt time.Time `db:"refresh_token_expires_at"`
	RefreshTokenRequestIP net.IP    `db:"refresh_token_request_ip"`
	CreateRequestIP       net.IP    `db:"create_request_ip"`

	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt time.Time  `db:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at"`
}

func (db *SessionDBModel) ToEntity() *entity.Session {
	return &entity.Session{
		ID:                    db.ID,
		AccountID:             db.AccountID,
		RefreshToken:          db.RefreshToken,
		RefreshVersion:        db.RefreshVersion,
		RefreshTokenIssuedAt:  db.RefreshTokenIssuedAt,
		RefreshTokenExpiresAt: db.RefreshTokenExpiresAt,
		RefreshTokenRequestIP: db.RefreshTokenRequestIP,
		CreateRequestIP:       db.CreateRequestIP,

		CreatedAt: db.CreatedAt,
		UpdatedAt: db.UpdatedAt,
		DeletedAt: db.DeletedAt,
	}
}

func MapSessionEntityToDBModel(entity *entity.Session) *SessionDBModel {
	return &SessionDBModel{
		ID:                    entity.ID,
		AccountID:             entity.AccountID,
		RefreshToken:          entity.RefreshToken,
		RefreshVersion:        entity.RefreshVersion,
		RefreshTokenIssuedAt:  entity.RefreshTokenIssuedAt,
		RefreshTokenExpiresAt: entity.RefreshTokenExpiresAt,
		RefreshTokenRequestIP: entity.RefreshTokenRequestIP,
		CreateRequestIP:       entity.CreateRequestIP,

		CreatedAt: entity.CreatedAt,
		UpdatedAt: entity.UpdatedAt,
		DeletedAt: entity.DeletedAt,
	}
}
