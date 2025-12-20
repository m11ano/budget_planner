package redis

import (
	"cloud.google.com/go/civil"
	"github.com/google/uuid"
	"github.com/govalues/decimal"
	"github.com/m11ano/budget_planner/backend/ledger/internal/domain/budget/entity"
)

type ReportItemModelItem struct {
	Sum          *decimal.Decimal `json:"sum"`
	Period       civil.Date       `json:"period"`
	CategoryID   uint64           `json:"categoryID"`
	BudgetID     *uuid.UUID       `json:"budgetID"`
	BudgetAmount *decimal.Decimal `json:"budgetAmount"`
}

type ReportItemModel struct {
	AccountID uuid.UUID              `json:"accountID"`
	DateFrom  civil.Date             `json:"dateFrom"`
	DateTo    civil.Date             `json:"dateTo"`
	Items     []*ReportItemModelItem `json:"items"`
}

func (db *ReportItemModel) ToEntity() *entity.ReportItem {
	if db == nil {
		return nil
	}

	items := make([]*entity.AccountTransactionReportItem, 0, len(db.Items))
	for _, it := range db.Items {
		if it == nil {
			continue
		}

		items = append(items, &entity.AccountTransactionReportItem{
			Sum:          it.Sum,
			Period:       it.Period,
			CategoryID:   it.CategoryID,
			BudgetID:     it.BudgetID,
			BudgetAmount: it.BudgetAmount,
		})
	}

	return &entity.ReportItem{
		AccountID: db.AccountID,
		DateFrom:  db.DateFrom,
		DateTo:    db.DateTo,
		Items:     items,
	}
}

func MapTransactionEntityToDBModel(e *entity.ReportItem) *ReportItemModel {
	if e == nil {
		return nil
	}

	items := make([]*ReportItemModelItem, 0, len(e.Items))
	for _, it := range e.Items {
		if it == nil {
			continue
		}

		items = append(items, &ReportItemModelItem{
			Sum:          it.Sum,
			Period:       it.Period,
			CategoryID:   it.CategoryID,
			BudgetID:     it.BudgetID,
			BudgetAmount: it.BudgetAmount,
		})
	}

	return &ReportItemModel{
		AccountID: e.AccountID,
		DateFrom:  e.DateFrom,
		DateTo:    e.DateTo,
		Items:     items,
	}
}
