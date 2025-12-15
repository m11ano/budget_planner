package usecase

type Facade struct {
	Account AccountUsecase
	Session SessionUsecase
	Auth    AuthUsecase
}

func NewFacade(
	accountUC AccountUsecase,
	sessionUC SessionUsecase,
	authUC AuthUsecase,
) *Facade {
	return &Facade{
		Account: accountUC,
		Session: sessionUC,
		Auth:    authUC,
	}
}
