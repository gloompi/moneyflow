// Package db contains product related CRUD functionality.
package db

import (
	"context"
	"fmt"

	"github.com/gloompi/ultimate-service/business/sys/database"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

// Store manages the set of APIs for user access.
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

// Create adds an Income to the database. It returns the created Income with
// fields like ID and DateCreated populated.
func (s Store) Create(ctx context.Context, inc Income) error {
	const q = `
	INSERT INTO incomes
		(income_id, user_id, name, category, currency, amount, reocurrence, duration, reocurrence_type, duration_type, date_created, date_updated)
	VALUES
		(:income_id, :user_id, :name, :category, :currency, amount, reocurrence, duration, reocurrence_type, duration_type, :date_created, :date_updated)`

	if err := database.NamedExecContext(ctx, s.log, s.db, q, inc); err != nil {
		return fmt.Errorf("inserting income: %w", err)
	}

	return nil
}

// Update modifies data about a Income. It will error if the specified ID is
// invalid or does not reference an existing Income.
func (s Store) Update(ctx context.Context, inc Income) error {
	const q = `
	UPDATE
		incomes
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
		income_id = :income_id`

	if err := database.NamedExecContext(ctx, s.log, s.db, q, inc); err != nil {
		return fmt.Errorf("updating income incomeID[%s]: %w", inc.ID, err)
	}

	return nil
}

// Delete removes the income identified by a given ID.
func (s Store) Delete(ctx context.Context, incomeID string) error {
	data := struct {
		IncomeID string `db:"income_id"`
	}{
		IncomeID: incomeID,
	}

	const q = `
	DELETE FROM
		incomes
	WHERE
		income_id = :income_id`

	if err := database.NamedExecContext(ctx, s.log, s.db, q, data); err != nil {
		return fmt.Errorf("deleting income incomeID[%s]: %w", incomeID, err)
	}

	return nil
}

// Query gets all Incomes from the database.
func (s Store) Query(ctx context.Context, pageNumber int, rowsPerPage int) ([]Income, error) {
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
		incomes
	ORDER BY
		date_created
	OFFSET :offset ROWS FETCH NEXT :rows_per_page ROWS ONLY`

	var incs []Income
	if err := database.NamedQuerySlice(ctx, s.log, s.db, q, data, &incs); err != nil {
		return nil, fmt.Errorf("selecting incomes: %w", err)
	}

	return incs, nil
}

// QueryByID finds the income identified by a given ID.
func (s Store) QueryByID(ctx context.Context, incomeID string) (Income, error) {
	data := struct {
		IncomeID string `db:"income_id"`
	}{
		IncomeID: incomeID,
	}

	const q = `
	SELECT
		*
	FROM
		incomes
	WHERE 
		income_id = :income_id`

	var inc Income
	if err := database.NamedQueryStruct(ctx, s.log, s.db, q, data, &inc); err != nil {
		return Income{}, fmt.Errorf("selecting income incomeID[%q]: %w", incomeID, err)
	}

	return inc, nil
}

// QueryByUserID finds the income identified by a given User ID.
func (s Store) QueryByUserID(ctx context.Context, userID string) ([]Income, error) {
	data := struct {
		UserID string `db:"user_id"`
	}{
		UserID: userID,
	}

	const q = `
	SELECT
		*
	FROM
		incomes
	WHERE
		user_id = :user_id`

	var incs []Income
	if err := database.NamedQuerySlice(ctx, s.log, s.db, q, data, &incs); err != nil {
		return nil, fmt.Errorf("selecting incomes userID[%s]: %w", userID, err)
	}

	return incs, nil
}
