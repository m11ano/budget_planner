package providing

import "github.com/m11ano/budget_planner/backend/auth/pkg/backoff"

func NewBackoff() *backoff.Controller {
	return backoff.NewController()
}
