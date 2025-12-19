package pg

import (
	"time"

	"cloud.google.com/go/civil"
	"github.com/google/uuid"
	"github.com/govalues/decimal"
	"github.com/m11ano/budget_planner/backend/ledger/internal/domain/budget/entity"
	"github.com/m11ano/budget_planner/backend/ledger/pkg/dbhelper"
)

const (
	BudgetTable = "budget"
)

var BudgetTableFields = []string{}

func init() {
	BudgetTableFields = dbhelper.ExtractDBFields(&BudgetDBModel{})
}

type BudgetDBModel struct {
	ID         uuid.UUID       `db:"id"`
	AccountID  uuid.UUID       `db:"account_id"`
	Period     civil.Date      `db:"period"`
	CategoryID uint64          `db:"category_id"`
	Amount     decimal.Decimal `db:"amount"`

	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt time.Time  `db:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at"`
}

func (db *BudgetDBModel) ToEntity() *entity.Budget {
	return &entity.Budget{
		ID:         db.ID,
		AccountID:  db.AccountID,
		Amount:     db.Amount,
		CategoryID: db.CategoryID,
		Period:     db.Period,

		CreatedAt: db.CreatedAt,
		UpdatedAt: db.UpdatedAt,
		DeletedAt: db.DeletedAt,
	}
}

func MapBudgetEntityToDBModel(entity *entity.Budget) *BudgetDBModel {
	return &BudgetDBModel{
		ID:         entity.ID,
		AccountID:  entity.AccountID,
		Amount:     entity.Amount,
		CategoryID: entity.CategoryID,
		Period:     entity.Period,

		CreatedAt: entity.CreatedAt,
		UpdatedAt: entity.UpdatedAt,
		DeletedAt: entity.DeletedAt,
	}
}
