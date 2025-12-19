package usecase

type Facade struct {
	Category    CategoryUsecase
	Transaction TransactionUsecase
	Budget      BudgetUsecase
}

func NewFacade(
	category CategoryUsecase,
	transactionUC TransactionUsecase,
	budgetUC BudgetUsecase,
) *Facade {
	return &Facade{
		Category:    category,
		Transaction: transactionUC,
		Budget:      budgetUC,
	}
}
