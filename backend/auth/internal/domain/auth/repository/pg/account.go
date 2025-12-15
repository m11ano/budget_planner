package pg

import (
	"net"
	"time"

	"github.com/google/uuid"
	"github.com/m11ano/budget_planner/backend/auth/internal/domain/auth/entity"
	"github.com/m11ano/budget_planner/backend/auth/pkg/dbhelper"
)

const (
	AccountTable = "account"
)

var AccountTableFields = []string{}

func init() {
	AccountTableFields = dbhelper.ExtractDBFields(&AccountDBModel{})
}

type AccountDBModel struct {
	ID             uuid.UUID  `db:"id"`
	Email          string     `db:"email"`
	PasswordHash   string     `db:"password_hash"`
	IsConfirmed    bool       `db:"is_confirmed"`
	IsBlocked      bool       `db:"is_blocked"`
	LastLoginAt    *time.Time `db:"last_login_at"`
	LastRequestAt  *time.Time `db:"last_request_at"`
	LastRequestIP  *net.IP    `db:"last_request_ip"`
	ProfileName    string     `db:"profile_name"`
	ProfileSurname string     `db:"profile_surname"`

	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt time.Time  `db:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at"`
}

func (db *AccountDBModel) ToEntity() *entity.Account {
	return &entity.Account{
		ID:             db.ID,
		Email:          db.Email,
		PasswordHash:   db.PasswordHash,
		IsConfirmed:    db.IsConfirmed,
		IsBlocked:      db.IsBlocked,
		LastLoginAt:    db.LastLoginAt,
		LastRequestAt:  db.LastRequestAt,
		LastRequestIP:  db.LastRequestIP,
		ProfileName:    db.ProfileName,
		ProfileSurname: db.ProfileSurname,

		CreatedAt: db.CreatedAt,
		UpdatedAt: db.UpdatedAt,
		DeletedAt: db.DeletedAt,
	}
}

func MapAccountEntityToDBModel(entity *entity.Account) *AccountDBModel {
	return &AccountDBModel{
		ID:             entity.ID,
		Email:          entity.Email,
		PasswordHash:   entity.PasswordHash,
		IsConfirmed:    entity.IsConfirmed,
		IsBlocked:      entity.IsBlocked,
		LastLoginAt:    entity.LastLoginAt,
		LastRequestAt:  entity.LastRequestAt,
		LastRequestIP:  entity.LastRequestIP,
		ProfileName:    entity.ProfileName,
		ProfileSurname: entity.ProfileSurname,

		CreatedAt: entity.CreatedAt,
		UpdatedAt: entity.UpdatedAt,
		DeletedAt: entity.DeletedAt,
	}
}
