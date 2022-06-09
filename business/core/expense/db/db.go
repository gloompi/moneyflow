// Package db contains product related CRUD functionality.
package db

import (
	"context"
	"fmt"

	"github.com/gloompi/ultimate-service/business/sys/database"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

// Store manages the set of APIs for expense access.
type Store struct {
	log          *zap.SugaredLogger
	tr           database.Transactor
	db           sqlx.ExtContext
	isWithinTran bool
}

// NewStore constructs a data for api access.
func NewStore(log *zap.SugaredLogger, db *sqlx.DB) Store {
	return Store{
		log: log,
		tr:  db,
		db:  db,
	}
}

// WithinTran runs passed function and do commit/rollback at the end.
func (s Store) WithinTran(ctx context.Context, fn func(sqlx.ExtContext) error) error {
	if s.isWithinTran {
		fn(s.db)
	}
	return database.WithinTran(ctx, s.log, s.tr, fn)
}

// Tran return new Store with transaction in it.
func (s Store) Tran(tx sqlx.ExtContext) Store {
	return Store{
		log:          s.log,
		tr:           s.tr,
		db:           tx,
		isWithinTran: true,
	}
}

// Create adds an Expense to the database. It returns the created Expense with
// fields like ID and DateCreated populated.
func (s Store) Create(ctx context.Context, exp Expense) error {
	const q = `
	INSERT INTO expenses
		(expense_id, user_id, name, category, currency, amount, reocurrence, duration, reocurrence_type, duration_type, date_created, date_updated)
	VALUES
		(:expense_id, :user_id, :name, :category, :currency, :amount, :reocurrence, :duration, :reocurrence_type, :duration_type, :date_created, :date_updated)`

	if err := database.NamedExecContext(ctx, s.log, s.db, q, exp); err != nil {
		return fmt.Errorf("inserting expense: %w", err)
	}

	return nil
}

// Update modifies data about a Expense. It will error if the specified ID is
// invalid or does not reference an existing Expense.
func (s Store) Update(ctx context.Context, exp Expense) error {
	const q = `
	UPDATE
		expenses
	SET
		"name" = :name,
		"category" = :category,
		"currency" = :currency,
		"amount" = :amount,
		"reocurrence" = :reocurrence,
		"duration" = :duration,
		"reocurrence_type" = :reocurrence_type,
		"duration_type" = :duration_type,
		"date_updated" = :date_updated
	WHERE
		expense_id = :expense_id`

	if err := database.NamedExecContext(ctx, s.log, s.db, q, exp); err != nil {
		return fmt.Errorf("updating expense expenseID[%s]: %w", exp.ID, err)
	}

	return nil
}

// Delete removes the expense identified by a given ID.
func (s Store) Delete(ctx context.Context, expenseID string) error {
	data := struct {
		ExpenseID string `db:"expense_id"`
	}{
		ExpenseID: expenseID,
	}

	const q = `
	DELETE FROM
		expenses
	WHERE
		expense_id = :expense_id`

	if err := database.NamedExecContext(ctx, s.log, s.db, q, data); err != nil {
		return fmt.Errorf("deleting expense expenseID[%s]: %w", expenseID, err)
	}

	return nil
}

// Query gets all Expenses from the database.
func (s Store) Query(ctx context.Context, pageNumber int, rowsPerPage int) ([]Expense, error) {
	data := struct {
		Offset      int `db:"offset"`
		RowsPerPage int `db:"rows_per_page"`
	}{
		Offset:      (pageNumber - 1) * rowsPerPage,
		RowsPerPage: rowsPerPage,
	}

	const q = `
	SELECT
		*
	FROM
		expenses
	ORDER BY
		date_created
	OFFSET :offset ROWS FETCH NEXT :rows_per_page ROWS ONLY`

	var exps []Expense
	if err := database.NamedQuerySlice(ctx, s.log, s.db, q, data, &exps); err != nil {
		return nil, fmt.Errorf("selecting expenses: %w", err)
	}

	return exps, nil
}

// QueryByID finds the expense identified by a given ID.
func (s Store) QueryByID(ctx context.Context, expenseID string) (Expense, error) {
	data := struct {
		ExpenseID string `db:"expense_id"`
	}{
		ExpenseID: expenseID,
	}

	const q = `
	SELECT
		*
	FROM
		expenses
	WHERE
		expense_id = :expense_id`

	var exp Expense
	if err := database.NamedQueryStruct(ctx, s.log, s.db, q, data, &exp); err != nil {
		return Expense{}, fmt.Errorf("selecting expense expenseID[%q]: %w", expenseID, err)
	}

	return exp, nil
}

// QueryByUserID finds the expense identified by a given User ID.
func (s Store) QueryByUserID(ctx context.Context, userID string) ([]Expense, error) {
	data := struct {
		UserID string `db:"user_id"`
	}{
		UserID: userID,
	}

	const q = `
	SELECT
		*
	FROM
		expenses
	WHERE
		user_id = :user_id`

	var exps []Expense
	if err := database.NamedQuerySlice(ctx, s.log, s.db, q, data, &exps); err != nil {
		return nil, fmt.Errorf("selecting expenses userID[%s]: %w", userID, err)
	}

	return exps, nil
}
