package entity

import (
	"cloud.google.com/go/civil"
	"github.com/google/uuid"
	"github.com/govalues/decimal"
)

type TransactionReportItem struct {
	AccountID  uuid.UUID
	Sum        decimal.Decimal
	Period     civil.Date
	CategoryID uint64
}
