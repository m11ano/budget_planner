package auth

import (
	"context"
	"time"

	"github.com/m11ano/budget_planner/backend/auth/internal/domain/auth/entity"
)

func (uc *UsecaseImpl) delete(ctx context.Context, session *entity.Session) error {
	nowTime := time.Now()
	session.DeletedAt = &nowTime
	return uc.repo.Update(ctx, session)
}
