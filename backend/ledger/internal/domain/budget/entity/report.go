package entity

import (
	"cloud.google.com/go/civil"
	"github.com/google/uuid"
)

type ReportItem struct {
	AccountID uuid.UUID
	DateFrom  civil.Date
	DateTo    civil.Date
	Items     []*TransactionReportItem
}
