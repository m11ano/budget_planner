package redis

import (
	"time"

	"cloud.google.com/go/civil"
	"github.com/google/uuid"
	"github.com/govalues/decimal"
	"github.com/m11ano/budget_planner/backend/ledger/internal/domain/budget/entity"
)

type BudgetPagedList struct {
	Items []*BudgetModel `json:"items"`
	Total uint64         `json:"total"`
}

type BudgetModel struct {
	ID         uuid.UUID       `json:"id"`
	AccountID  uuid.UUID       `json:"accountID"`
	Period     civil.Date      `json:"period"`
	CategoryID uint64          `json:"categoryID"`
	Amount     decimal.Decimal `json:"amount"`

	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
	DeletedAt *time.Time `json:"deletedAt"`
}

func (db *BudgetModel) ToEntity() *entity.Budget {
	if db == nil {
		return nil
	}

	return &entity.Budget{
		ID:         db.ID,
		AccountID:  db.AccountID,
		Period:     db.Period,
		CategoryID: db.CategoryID,
		Amount:     db.Amount,
		CreatedAt:  db.CreatedAt,
		UpdatedAt:  db.UpdatedAt,
		DeletedAt:  db.DeletedAt,
	}
}

func MapBudgetEntityToDBModel(e *entity.Budget) *BudgetModel {
	if e == nil {
		return nil
	}

	return &BudgetModel{
		ID:         e.ID,
		AccountID:  e.AccountID,
		Period:     e.Period,
		CategoryID: e.CategoryID,
		Amount:     e.Amount,
		CreatedAt:  e.CreatedAt,
		UpdatedAt:  e.UpdatedAt,
		DeletedAt:  e.DeletedAt,
	}
}
