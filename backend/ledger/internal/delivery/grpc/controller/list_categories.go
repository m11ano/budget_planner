package controller

import (
	"context"

	appErrors "github.com/m11ano/budget_planner/backend/ledger/internal/app/errors"
	"github.com/m11ano/budget_planner/backend/ledger/internal/common/uctypes"
	budgetUC "github.com/m11ano/budget_planner/backend/ledger/internal/domain/budget/usecase"
	desc "github.com/m11ano/budget_planner/backend/ledger/pkg/proto_pb/ledger_service"
)

func (c *controller) ListCategories(ctx context.Context, req *desc.ListCategoriesRequest) (*desc.ListCategoriesResponse, error) {
	const op = "ListCategories"

	items, err := c.budgetFacade.Category.FindList(
		ctx,
		&budgetUC.CategoryListOptions{
			Sort: []uctypes.SortOption[budgetUC.CategoryListOptionsSortField]{
				{
					Field:  budgetUC.CategoryListOptionsSortFieldID,
					IsDesc: false,
				},
			},
		},
		nil,
	)
	if err != nil {
		return nil, appErrors.Chainf(err, "%s.%s", c.pkg, op)
	}

	out := &desc.ListCategoriesResponse{
		Items: make([]*desc.Category, 0, len(items)),
	}

	for _, item := range items {
		out.Items = append(out.Items, &desc.Category{
			Id:    int64(item.Category.ID),
			Title: item.Category.Title,
		})
	}

	return out, nil
}
