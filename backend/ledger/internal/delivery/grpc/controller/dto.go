package controller

import (
	budgetUC "github.com/m11ano/budget_planner/backend/ledger/internal/domain/budget/usecase"
	desc "github.com/m11ano/budget_planner/backend/ledger/pkg/proto_pb/ledger_service"
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
