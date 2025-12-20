package transaction

import (
	"bytes"
	"context"
	"encoding/csv"
	"fmt"
	"strconv"
	"strings"

	"cloud.google.com/go/civil"
	"github.com/google/uuid"
	"github.com/govalues/decimal"
	appErrors "github.com/m11ano/budget_planner/backend/ledger/internal/app/errors"
	"github.com/m11ano/budget_planner/backend/ledger/internal/domain/budget/entity"
	"github.com/m11ano/budget_planner/backend/ledger/internal/domain/budget/usecase"
)

func (r *Repository) ItemsToCSV(ctx context.Context, items []*usecase.TransactionDTO) ([]byte, error) {
	const op = "ItemsToCSV"

	buf := &bytes.Buffer{}
	writer := csv.NewWriter(buf)

	buf.Write([]byte{0xEF, 0xBB, 0xBF})

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

func (r *Repository) ItemsFromCSV(ctx context.Context, data []byte, accountID uuid.UUID) ([]*entity.Transaction, error) {
	const op = "ItemsFromCSV"

	reader := csv.NewReader(bytes.NewReader(data))

	records, err := reader.ReadAll()
	if err != nil {
		return nil, appErrors.Chainf(err, "%s.%s", r.pkg, op)
	}

	out := make([]*entity.Transaction, 0, len(records))

	for idx, row := range records {
		if len(row) < 5 {
			return nil, appErrors.Chainf(
				appErrors.ErrBadRequest.WithHints("invalid number of columns", fmt.Sprintf("line %d", idx+1)),
				"%s.%s",
				r.pkg,
				op,
			)
		}

		isIncome := false
		if strings.ToLower(row[0]) == "true" {
			isIncome = true
		}

		amount, err := decimal.Parse(row[1])
		if err != nil {
			return nil, appErrors.Chainf(
				appErrors.ErrBadRequest.WithParent(err).WithHints("invalid amount", fmt.Sprintf("line %d", idx)),
				"%s.%s",
				r.pkg,
				op,
			)
		}

		occurredOn, err := civil.ParseDate(row[2])
		if err != nil {
			return nil, appErrors.Chainf(
				appErrors.ErrBadRequest.WithParent(err).WithHints("invalid occurred_on", fmt.Sprintf("line %d", idx+1)),
				"%s.%s",
				r.pkg,
				op,
			)
		}

		categoryID, err := strconv.ParseUint(row[4], 10, 64)
		if err != nil {
			return nil, appErrors.Chainf(
				appErrors.ErrBadRequest.WithParent(err).WithHints("invalid category_id", fmt.Sprintf("line %d", idx+1)),
				"%s.%s",
				r.pkg,
				op,
			)
		}

		transaction, err := entity.NewTransaction(
			accountID,
			isIncome,
			amount,
			occurredOn,
			categoryID,
		)
		if err != nil {
			appErr, ok := appErrors.ExtractError(err)
			if ok {
				hints := appErr.Hints()
				hints = append(hints, fmt.Sprintf("line %d", idx+1))

				return nil, appErrors.Chainf(
					appErr.WithHints(hints...),
					"%s.%s: line %d",
					r.pkg,
					op,
					idx+1,
				)
			}
			return nil, appErrors.Chainf(err, "%s.%s: line %d", r.pkg, op, idx+1)
		}

		err = transaction.SetDescription(row[3])
		if err != nil {
			appErr, ok := appErrors.ExtractError(err)
			if ok {
				hints := appErr.Hints()
				hints = append(hints, fmt.Sprintf("line %d", idx+1))

				return nil, appErrors.Chainf(
					appErr.WithHints(hints...),
					"%s.%s: line %d",
					r.pkg,
					op,
					idx+1,
				)
			}
			return nil, appErrors.Chainf(err, "%s.%s: line %d", r.pkg, op, idx+1)
		}

		out = append(out, transaction)
	}

	return out, nil
}
