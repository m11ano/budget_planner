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
	TransactionTable = "transaction"
)

var TransactionTableFields = []string{}

func init() {
	TransactionTableFields = dbhelper.ExtractDBFields(&TransactionDBModel{})
}

type TransactionDBModel struct {
	ID          uuid.UUID       `db:"id"`
	AccountID   uuid.UUID       `db:"account_id"`
	IsIncome    bool            `db:"is_income"`
	Amount      decimal.Decimal `db:"amount"`
	OccurredOn  civil.Date      `db:"occurred_on"`
	CategoryID  uint64          `db:"category_id"`
	Description string          `db:"description"`

	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt time.Time  `db:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at"`
}

func (db *TransactionDBModel) ToEntity() *entity.Transaction {
	return &entity.Transaction{
		ID:          db.ID,
		AccountID:   db.AccountID,
		IsIncome:    db.IsIncome,
		Amount:      db.Amount,
		OccurredOn:  db.OccurredOn,
		CategoryID:  db.CategoryID,
		Description: db.Description,

		CreatedAt: db.CreatedAt,
		UpdatedAt: db.UpdatedAt,
		DeletedAt: db.DeletedAt,
	}
}

func MapTransactionEntityToDBModel(entity *entity.Transaction) *TransactionDBModel {
	return &TransactionDBModel{
		ID:          entity.ID,
		AccountID:   entity.AccountID,
		IsIncome:    entity.IsIncome,
		Amount:      entity.Amount,
		OccurredOn:  entity.OccurredOn,
		CategoryID:  entity.CategoryID,
		Description: entity.Description,

		CreatedAt: entity.CreatedAt,
		UpdatedAt: entity.UpdatedAt,
		DeletedAt: entity.DeletedAt,
	}
}
