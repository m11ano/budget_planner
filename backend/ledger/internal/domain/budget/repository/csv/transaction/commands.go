package transaction

import (
	"bytes"
	"context"
	"encoding/csv"
	"fmt"

	appErrors "github.com/m11ano/budget_planner/backend/ledger/internal/app/errors"
	"github.com/m11ano/budget_planner/backend/ledger/internal/domain/budget/usecase"
)

func (r *Repository) ItemsToCSV(ctx context.Context, items []*usecase.TransactionDTO) ([]byte, error) {
	const op = "ItemsToCSV"

	buf := &bytes.Buffer{}
	writer := csv.NewWriter(buf)

	buf.WriteString("\ufeff")

	headers := []string{
		"ID",
		"AccountID",
		"IsIncome",
		"Amount",
		"OccurredOn",
		"Description",
		"CategoryID",
		"CategoryTitle",
	}

	if err := writer.Write(headers); err != nil {
		return nil, appErrors.Chainf(err, "%s.%s", r.pkg, op)
	}

	for _, row := range items {
		row := []string{
			row.Transaction.ID.String(),
			row.Transaction.AccountID.String(),
			fmt.Sprintf("%t", row.Transaction.IsIncome),
			row.Transaction.Amount.String(),
			row.Transaction.OccurredOn.String(),
			row.Transaction.Description,
			fmt.Sprintf("%d", row.Transaction.CategoryID),
			row.Category.Title,
		}

		if err := writer.Write(row); err != nil {
			return nil, appErrors.Chainf(err, "%s.%s", r.pkg, op)
		}
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		return nil, appErrors.Chainf(err, "%s.%s", r.pkg, op)
	}

	return buf.Bytes(), nil
}
