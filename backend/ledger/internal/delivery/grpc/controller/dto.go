package controller

import (
	budgetEntity "github.com/m11ano/budget_planner/backend/ledger/internal/domain/budget/entity"
	budgetUC "github.com/m11ano/budget_planner/backend/ledger/internal/domain/budget/usecase"
	desc "github.com/m11ano/budget_planner/backend/ledger/pkg/proto_pb/ledger_service"
	"github.com/samber/lo"
)

func TransactionToProto(itemDTO *budgetUC.TransactionDTO) *desc.Transaction {
	return &desc.Transaction{
		Id:        itemDTO.Transaction.ID.String(),
		AccountId: itemDTO.Transaction.AccountID.String(),
		IsIncome:  itemDTO.Transaction.IsIncome,
		Amount:    itemDTO.Transaction.Amount.String(),
		OccurredOn: &desc.Date{
			Year:  int32(itemDTO.Transaction.OccurredOn.Year),
			Month: int32(itemDTO.Transaction.OccurredOn.Month),
			Day:   int32(itemDTO.Transaction.OccurredOn.Day),
		},
		CategoryId:  int64(itemDTO.Transaction.CategoryID),
		Description: itemDTO.Transaction.Description,
		CreatedAt:   toProtoTimestamp(&itemDTO.Transaction.CreatedAt),
		UpdatedAt:   toProtoTimestamp(&itemDTO.Transaction.UpdatedAt),
	}
}

func BudgetToProto(itemDTO *budgetUC.BudgetDTO) *desc.Budget {
	return &desc.Budget{
		Id:        itemDTO.Budget.ID.String(),
		AccountId: itemDTO.Budget.AccountID.String(),
		Amount:    itemDTO.Budget.Amount.String(),
		Period: &desc.DateMonth{
			Year:  int32(itemDTO.Budget.Period.Year),
			Month: int32(itemDTO.Budget.Period.Month),
		},
		CategoryId: int64(itemDTO.Budget.CategoryID),
		CreatedAt:  toProtoTimestamp(&itemDTO.Budget.CreatedAt),
		UpdatedAt:  toProtoTimestamp(&itemDTO.Budget.UpdatedAt),
	}
}

func ReportToProto(itemDTO *budgetEntity.ReportItem) *desc.PeriodReport {
	periodReport := &desc.PeriodReport{
		PeriodStart: &desc.Date{
			Year:  int32(itemDTO.DateFrom.Year),
			Month: int32(itemDTO.DateFrom.Month),
			Day:   int32(itemDTO.DateFrom.Day),
		},
		PeriodEnd: &desc.Date{
			Year:  int32(itemDTO.DateTo.Year),
			Month: int32(itemDTO.DateTo.Month),
			Day:   int32(itemDTO.DateTo.Day),
		},
		Items: make([]*desc.ReportItem, 0, len(itemDTO.Items)),
	}

	for _, item := range itemDTO.Items {
		repItem := &desc.ReportItem{
			CategoryId: int64(item.CategoryID),
		}

		if item.Sum != nil {
			repItem.Sum = lo.ToPtr(item.Sum.String())
		}

		spentBudget, err := item.SpentBudget()
		if err == nil && spentBudget != nil {
			repItem.SpentBudget = lo.ToPtr(spentBudget.String())
		}

		if item.BudgetAmount != nil {
			repItem.ItemBudget = lo.ToPtr(item.BudgetAmount.String())
		}

		periodReport.Items = append(periodReport.Items, repItem)
	}

	return periodReport
}
