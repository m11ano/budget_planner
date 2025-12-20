package ledger

import (
	"fmt"
	"time"

	"cloud.google.com/go/civil"
	desc "github.com/m11ano/budget_planner/backend/gateway/pkg/proto_pb/ledger_service"
)

type TransactionOutput struct {
	ID          string     `json:"id"`
	AccountID   string     `json:"accountID"`
	IsIncome    bool       `json:"isIncome"`
	Amount      string     `json:"amount"`
	OccurredOn  civil.Date `json:"occurredOn" swaggertype:"string" example:"2025-12-20"`
	CategoryID  uint64     `json:"categoryID"`
	Description string     `json:"description"`
	CreatedAt   *time.Time `json:"createdAt"`
	UpdatedAt   *time.Time `json:"updatedAt"`
}

func NewTransactionOutput(transaction *desc.Transaction) *TransactionOutput {
	occurredOn, err := civil.ParseDate(fmt.Sprintf(
		"%04d-%02d-%02d",
		transaction.OccurredOn.Year,
		transaction.OccurredOn.Month,
		transaction.OccurredOn.Day,
	))
	if err != nil {
		panic(err)
	}

	return &TransactionOutput{
		ID:          transaction.Id,
		AccountID:   transaction.AccountId,
		IsIncome:    transaction.IsIncome,
		Amount:      transaction.Amount,
		OccurredOn:  occurredOn,
		CategoryID:  uint64(transaction.CategoryId),
		Description: transaction.Description,
		CreatedAt:   fromProtoTimestamp(transaction.CreatedAt),
		UpdatedAt:   fromProtoTimestamp(transaction.UpdatedAt),
	}
}

type BudgetOutputPeriod struct {
	Month int32 `json:"month"`
	Year  int32 `json:"year"`
}

type BudgetOutput struct {
	ID         string             `json:"id"`
	AccountID  string             `json:"accountID"`
	Amount     string             `json:"amount"`
	Period     BudgetOutputPeriod `json:"period"`
	CategoryID uint64             `json:"categoryID"`
	CreatedAt  *time.Time         `json:"createdAt"`
	UpdatedAt  *time.Time         `json:"updatedAt"`
}

func NewBudgetOutput(budget *desc.Budget) *BudgetOutput {
	return &BudgetOutput{
		ID:        budget.Id,
		AccountID: budget.AccountId,
		Amount:    budget.Amount,
		Period: BudgetOutputPeriod{
			Month: budget.Period.Month,
			Year:  budget.Period.Year,
		},
		CategoryID: uint64(budget.CategoryId),
		CreatedAt:  fromProtoTimestamp(budget.CreatedAt),
		UpdatedAt:  fromProtoTimestamp(budget.UpdatedAt),
	}
}

type CategoryOutput struct {
	ID    uint64 `json:"id"`
	Title string `json:"title"`
}

func NewCategoryOutput(category *desc.Category) *CategoryOutput {
	return &CategoryOutput{
		ID:    uint64(category.Id),
		Title: category.Title,
	}
}

type ReportOutputItem struct {
	CategoryID  uint64  `json:"categoryID"`
	Sum         *string `json:"sum"`
	SpentBudget *string `json:"spentBudget"`
	ItemBudget  *string `json:"itemBudget"`
}

type ReportOutput struct {
	PeriodStart civil.Date          `json:"periodStart" swaggertype:"string" example:"2025-01-01"`
	PeriodEnd   civil.Date          `json:"periodEnd" swaggertype:"string" example:"2025-12-20"`
	Items       []*ReportOutputItem `json:"items"`
}

func NewReportOutput(report *desc.PeriodReport) *ReportOutput {
	periodStart, err := civil.ParseDate(fmt.Sprintf(
		"%04d-%02d-%02d",
		report.PeriodStart.Year,
		report.PeriodStart.Month,
		report.PeriodStart.Day,
	))
	if err != nil {
		panic(err)
	}

	periodEnd, err := civil.ParseDate(fmt.Sprintf(
		"%04d-%02d-%02d",
		report.PeriodEnd.Year,
		report.PeriodEnd.Month,
		report.PeriodEnd.Day,
	))
	if err != nil {
		panic(err)
	}

	result := &ReportOutput{
		PeriodStart: periodStart,
		PeriodEnd:   periodEnd,
		Items:       make([]*ReportOutputItem, 0, len(report.Items)),
	}

	for _, item := range report.Items {
		result.Items = append(result.Items, &ReportOutputItem{
			CategoryID:  uint64(item.CategoryId),
			Sum:         item.Sum,
			SpentBudget: item.SpentBudget,
			ItemBudget:  item.ItemBudget,
		})
	}

	return result
}
