package transaction

import (
	"context"
	"log/slog"

	"cloud.google.com/go/civil"
	"github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/google/uuid"
	"github.com/govalues/decimal"
	appErrors "github.com/m11ano/budget_planner/backend/ledger/internal/app/errors"
	"github.com/m11ano/budget_planner/backend/ledger/internal/common/uctypes"
	"github.com/m11ano/budget_planner/backend/ledger/internal/domain/budget/entity"
	"github.com/m11ano/budget_planner/backend/ledger/internal/domain/budget/repository/pg"
	"github.com/m11ano/budget_planner/backend/ledger/internal/domain/budget/usecase"
	"github.com/m11ano/budget_planner/backend/ledger/internal/infra/loghandler"
)

func (r *Repository) buildWhereForList(listOptions *usecase.TransactionListOptions, withDeleted bool) (where squirrel.And) {
	defer func() {
		if !withDeleted {
			where = append(where, squirrel.Expr("deleted_at IS NULL"))
		}
	}()

	if listOptions == nil {
		return where
	}

	if listOptions.FilterAccountID != nil {
		where = append(where, squirrel.Eq{"account_id": *listOptions.FilterAccountID})
	}

	return where
}

func (r *Repository) buildSortForList(listOptions *usecase.TransactionListOptions) []string {
	if listOptions == nil || len(listOptions.Sort) == 0 {
		return []string{"created_at DESC"}
	}

	sort := make([]string, 0, len(listOptions.Sort))

	for _, sortOption := range listOptions.Sort {
		switch sortOption.Field {
		case usecase.TransactionListOptionsSortFieldOccurredOn:
			if sortOption.IsDesc {
				sort = append(sort, "occurred_on DESC")
			} else {
				sort = append(sort, "occurred_on ASC")
			}
		case usecase.TransactionListOptionsSortFieldCreatedAt:
			if sortOption.IsDesc {
				sort = append(sort, "created_at DESC")
			} else {
				sort = append(sort, "created_at ASC")
			}
		}
	}

	return sort
}

func (r *Repository) FindList(
	ctx context.Context,
	listOptions *usecase.TransactionListOptions,
	queryParams *uctypes.QueryGetListParams,
) ([]*entity.Transaction, error) {
	const op = "FindList"

	withDeleted := queryParams != nil && queryParams.WithDeleted

	where := r.buildWhereForList(listOptions, withDeleted)

	fields := pg.TransactionTableFields

	q := r.qb.Select(fields...).From(pg.TransactionTable).Where(where)

	sort := r.buildSortForList(listOptions)
	if len(sort) > 0 {
		q = q.OrderBy(sort...)
	}

	if queryParams != nil {
		if queryParams.ForUpdateSkipLocked {
			q = q.Suffix("FOR UPDATE SKIP LOCKED")
		} else if queryParams.ForUpdate {
			q = q.Suffix("FOR UPDATE")
		} else if queryParams.ForShare {
			q = q.Suffix("FOR SHARE")
		}

		if queryParams.Limit > 0 {
			q = q.Limit(queryParams.Limit)
		}

		if queryParams.Offset > 0 {
			q = q.Offset(queryParams.Offset)
		}
	}

	query, args, err := q.ToSql()
	if err != nil {
		r.logger.ErrorContext(loghandler.WithSource(ctx), "building query", slog.Any("error", err))
		return nil, appErrors.Chainf(appErrors.ErrInternal.WithWrap(err), "%s.%s", r.pkg, op)
	}

	rows, err := r.pgClient.GetConn(ctx).Query(ctx, query, args...)
	if err != nil {
		convErr, ok := appErrors.ConvertPgxToAppErr(err)
		if !ok {
			r.logger.ErrorContext(loghandler.WithSource(ctx), "query row error", slog.Any("error", err))
		}
		return nil, appErrors.Chainf(convErr, "%s.%s", r.pkg, op)
	}

	defer rows.Close()

	dbData := []*pg.TransactionDBModel{}

	if err := pgxscan.ScanAll(&dbData, rows); err != nil {
		convErr, ok := appErrors.ConvertPgxToAppErr(err)
		if !ok {
			r.logger.ErrorContext(loghandler.WithSource(ctx), "scan row error", slog.Any("error", err))
		}
		return nil, appErrors.Chainf(convErr, "%s.%s", r.pkg, op)
	}

	result := make([]*entity.Transaction, 0, len(dbData))
	for _, dbItem := range dbData {
		result = append(result, dbItem.ToEntity())
	}

	return result, nil
}

func (r *Repository) FindPagedList(
	ctx context.Context,
	listOptions *usecase.TransactionListOptions,
	queryParams *uctypes.QueryGetListParams,
) ([]*entity.Transaction, uint64, error) {
	const op = "FindPagedList"

	withDeleted := queryParams != nil && queryParams.WithDeleted

	where := r.buildWhereForList(listOptions, withDeleted)

	fields := pg.TransactionTableFields

	q := r.qb.Select(fields...).From(pg.TransactionTable).Where(where)

	sort := r.buildSortForList(listOptions)
	if len(sort) > 0 {
		q = q.OrderBy(sort...)
	}

	if queryParams != nil {
		if queryParams.ForUpdateSkipLocked {
			q = q.Suffix("FOR UPDATE SKIP LOCKED")
		} else if queryParams.ForUpdate {
			q = q.Suffix("FOR UPDATE")
		} else if queryParams.ForShare {
			q = q.Suffix("FOR SHARE")
		}

		if queryParams.Limit > 0 {
			q = q.Limit(queryParams.Limit)
		}

		if queryParams.Offset > 0 {
			q = q.Offset(queryParams.Offset)
		}
	}

	query, args, err := q.ToSql()
	if err != nil {
		r.logger.ErrorContext(loghandler.WithSource(ctx), "building query", slog.Any("error", err))
		return nil, 0, appErrors.Chainf(appErrors.ErrInternal.WithWrap(err), "%s.%s", r.pkg, op)
	}

	totalQ := r.qb.Select("COUNT(*) as total").From(pg.TransactionTable).Where(where)
	totalQuery, totalArgs, err := totalQ.ToSql()
	if err != nil {
		r.logger.ErrorContext(loghandler.WithSource(ctx), "building query for total", slog.Any("error", err))
		return nil, 0, appErrors.Chainf(appErrors.ErrInternal.WithWrap(err), "%s.%s", r.pkg, op)
	}

	var total uint64
	var result []*entity.Transaction

	err = r.pgClient.Do(ctx, func(ctx context.Context) error {
		rows, err := r.pgClient.GetConn(ctx).Query(ctx, query, args...)
		if err != nil {
			convErr, ok := appErrors.ConvertPgxToAppErr(err)
			if !ok {
				r.logger.ErrorContext(loghandler.WithSource(ctx), "query row error", slog.Any("error", err))
			}
			return appErrors.Chainf(convErr, "%s.%s", r.pkg, op)
		}

		defer rows.Close()

		dbData := []*pg.TransactionDBModel{}

		if err := pgxscan.ScanAll(&dbData, rows); err != nil {
			convErr, ok := appErrors.ConvertPgxToAppErr(err)
			if !ok {
				r.logger.ErrorContext(loghandler.WithSource(ctx), "scan row error", slog.Any("error", err))
			}
			return appErrors.Chainf(convErr, "%s.%s", r.pkg, op)
		}

		result = make([]*entity.Transaction, 0, len(dbData))
		for _, dbItem := range dbData {
			result = append(result, dbItem.ToEntity())
		}

		row := r.pgClient.GetConn(ctx).QueryRow(ctx, totalQuery, totalArgs...)
		if err := row.Scan(&total); err != nil {
			convErr, ok := appErrors.ConvertPgxToAppErr(err)
			if !ok {
				r.logger.ErrorContext(loghandler.WithSource(ctx), "scan total error", slog.Any("error", err))
			}
			return appErrors.Chainf(convErr, "%s.%s", r.pkg, op)
		}

		return nil
	})
	if err != nil {
		return nil, 0, err
	}

	return result, total, nil
}

func (r *Repository) FindOneByID(
	ctx context.Context,
	id uuid.UUID,
	queryParams *uctypes.QueryGetOneParams,
) (*entity.Transaction, error) {
	const op = "FindOneByID"

	withDeleted := queryParams != nil && queryParams.WithDeleted

	where := squirrel.And{
		squirrel.Eq{"id": id},
	}

	if !withDeleted {
		where = append(where, squirrel.Expr("deleted_at IS NULL"))
	}

	q := r.qb.Select(pg.TransactionTableFields...).From(pg.TransactionTable).Where(where)

	if queryParams != nil {
		if queryParams.ForUpdateSkipLocked {
			q = q.Suffix("FOR UPDATE SKIP LOCKED")
		} else if queryParams.ForUpdate {
			q = q.Suffix("FOR UPDATE")
		} else if queryParams.ForShare {
			q = q.Suffix("FOR SHARE")
		}
	}

	query, args, err := q.ToSql()
	if err != nil {
		r.logger.ErrorContext(loghandler.WithSource(ctx), "building query", slog.Any("error", err))
		return nil, appErrors.Chainf(appErrors.ErrInternal.WithWrap(err), "%s.%s", r.pkg, op)
	}

	rows, err := r.pgClient.GetConn(ctx).Query(ctx, query, args...)
	if err != nil {
		convErr, ok := appErrors.ConvertPgxToAppErr(err)
		if !ok {
			r.logger.ErrorContext(loghandler.WithSource(ctx), "query row error", slog.Any("error", err))
		}
		return nil, appErrors.Chainf(convErr, "%s.%s", r.pkg, op)
	}

	defer rows.Close()

	dbData := &pg.TransactionDBModel{}

	if err := pgxscan.ScanOne(dbData, rows); err != nil {
		convErr, ok := appErrors.ConvertPgxToAppErr(err)
		if !ok {
			r.logger.ErrorContext(loghandler.WithSource(ctx), "scan row error", slog.Any("error", err))
		}
		return nil, appErrors.Chainf(convErr, "%s.%s", r.pkg, op)
	}

	return dbData.ToEntity(), nil
}

func (r *Repository) CountReportItems(
	ctx context.Context,
	filter *usecase.CountReportItemsFilter,
) ([]*entity.TransactionReportItem, error) {
	const op = "CountReportItems"

	type reportRow struct {
		AccountID  uuid.UUID       `db:"account_id"`
		Sum        decimal.Decimal `db:"sum"`
		Period     civil.Date      `db:"period"`
		CategoryID uint64          `db:"category_id"`
	}

	where := squirrel.And{
		squirrel.Expr("deleted_at IS NULL"),
	}

	if filter != nil {
		if filter.AccountID != nil {
			where = append(where, squirrel.Eq{"account_id": *filter.AccountID})
		}
		if filter.PeriodFrom != nil {
			where = append(where, squirrel.GtOrEq{"occurred_on": *filter.PeriodFrom})
		}
		if filter.PeriodTo != nil {
			where = append(where, squirrel.LtOrEq{"occurred_on": *filter.PeriodTo})
		}
		if filter.CategoryID != nil {
			where = append(where, squirrel.Eq{"category_id": *filter.CategoryID})
		}
		if len(filter.ExcludeIDs) > 0 {
			where = append(where, squirrel.NotEq{"id": filter.ExcludeIDs})
		}
	}

	periodExpr := "date_trunc('month', occurred_on)::date AS period"
	sumExpr := "SUM(amount) AS sum"

	selectCols := []string{
		"account_id",
		periodExpr,
		sumExpr,
	}

	groupBy := []string{
		"account_id",
		"period",
	}

	if filter != nil && filter.GroupByCategory {
		selectCols = append(selectCols, "category_id")
		groupBy = append(groupBy, "category_id")
	}

	q := r.qb.
		Select(selectCols...).
		From(pg.TransactionTable).
		Where(where).
		GroupBy(groupBy...).
		OrderBy("period ASC")

	query, args, err := q.ToSql()
	if err != nil {
		r.logger.ErrorContext(loghandler.WithSource(ctx), "building query", slog.Any("error", err))
		return nil, appErrors.Chainf(appErrors.ErrInternal.WithWrap(err), "%s.%s", r.pkg, op)
	}

	rows, err := r.pgClient.GetConn(ctx).Query(ctx, query, args...)
	if err != nil {
		convErr, ok := appErrors.ConvertPgxToAppErr(err)
		if !ok {
			r.logger.ErrorContext(loghandler.WithSource(ctx), "query error", slog.Any("error", err))
		}
		return nil, appErrors.Chainf(convErr, "%s.%s", r.pkg, op)
	}
	defer rows.Close()

	var dbData []*reportRow
	if err := pgxscan.ScanAll(&dbData, rows); err != nil {
		convErr, ok := appErrors.ConvertPgxToAppErr(err)
		if !ok {
			r.logger.ErrorContext(loghandler.WithSource(ctx), "scan error", slog.Any("error", err))
		}
		return nil, appErrors.Chainf(convErr, "%s.%s", r.pkg, op)
	}

	result := make([]*entity.TransactionReportItem, 0, len(dbData))
	for _, it := range dbData {
		result = append(result, &entity.TransactionReportItem{
			AccountID:  it.AccountID,
			Sum:        it.Sum,
			Period:     it.Period,
			CategoryID: it.CategoryID,
		})
	}

	return result, nil
}
