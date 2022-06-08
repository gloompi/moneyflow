// Package income provides a core business API. Right now these
package income

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/gloompi/ultimate-service/business/core/income/db"
	"github.com/gloompi/ultimate-service/business/sys/database"
	"github.com/gloompi/ultimate-service/business/sys/validate"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

// Set of error variables for CRUD operations.
var (
	ErrNotFound  = errors.New("income not found")
	ErrInvalidID = errors.New("ID is not in its proper form")
)

// Core manages the set of APIs for product access.
type Core struct {
	store db.Store
}

// NewCore constructs a core for product api access.
func NewCore(log *zap.SugaredLogger, sqlxDB *sqlx.DB) Core {
	return Core{
		store: db.NewStore(log, sqlxDB),
	}
}

// Create adds an Income to the database. It returns the created Income with
// fields like ID and DateCreated populated.
func (c Core) Create(ctx context.Context, ni NewIncome, now time.Time) (Income, error) {
	if err := validate.Check(ni); err != nil {
		return Income{}, fmt.Errorf("validating data: %w", err)
	}

	dbInc := db.Income{
		ID:               validate.GenerateID(),
		Name:             ni.Name,
		Category:         ni.Category,
		Currency:         ni.Currency,
		Amount:           ni.Amount,
		Reoccurrence:     ni.Reoccurrence,
		Duration:         ni.Duration,
		ReoccurrenceType: ni.ReoccurrenceType,
		DurationType:     ni.DurationType,
		UserID:           ni.UserID,
		DateCreated:      now,
		DateUpdated:      now,
	}

	if err := c.store.Create(ctx, dbInc); err != nil {
		return Income{}, fmt.Errorf("create: %w", err)
	}

	return toIncome(dbInc), nil
}

// Update modifies data about a Income. It will error if the specified ID is
// invalid or does not reference an existing Income.
func (c Core) Update(ctx context.Context, incomeID string, ui UpdateIncome, now time.Time) error {
	if err := validate.CheckID(incomeID); err != nil {
		return ErrInvalidID
	}

	if err := validate.Check(ui); err != nil {
		return fmt.Errorf("validating data: %w", err)
	}

	dbInc, err := c.store.QueryByID(ctx, incomeID)
	if err != nil {
		if errors.Is(err, database.ErrDBNotFound) {
			return ErrNotFound
		}
		return fmt.Errorf("updating income incomeID[%s]: %w", incomeID, err)
	}

	if ui.Name != nil {
		dbInc.Name = *ui.Name
	}
	if ui.Category != nil {
		dbInc.Category = *ui.Category
	}
	if ui.Currency != nil {
		dbInc.Currency = *ui.Currency
	}
	if ui.Amount != nil {
		dbInc.Amount = *ui.Amount
	}
	if ui.Reoccurrence != nil {
		dbInc.Reoccurrence = *ui.Reoccurrence
	}
	if ui.Duration != nil {
		dbInc.Duration = *ui.Duration
	}
	if ui.ReoccurrenceType != nil {
		dbInc.ReoccurrenceType = *ui.ReoccurrenceType
	}
	if ui.DurationType != nil {
		dbInc.DurationType = *ui.DurationType
	}
	dbInc.DateUpdated = now

	if err := c.store.Update(ctx, dbInc); err != nil {
		return fmt.Errorf("update: %w", err)
	}

	return nil
}

// Delete removes the income identified by a given ID.
func (c Core) Delete(ctx context.Context, incomeID string) error {
	if err := validate.CheckID(incomeID); err != nil {
		return ErrInvalidID
	}

	if err := c.store.Delete(ctx, incomeID); err != nil {
		return fmt.Errorf("delete: %w", err)
	}

	return nil
}

// Query gets all Incomes from the database.
func (c Core) Query(ctx context.Context, pageNumber int, rowsPerPage int) ([]Income, error) {
	dbIncs, err := c.store.Query(ctx, pageNumber, rowsPerPage)
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}

	return toIncomeSlice(dbIncs), nil
}

// // QueryByID finds the income identified by a given ID.
func (c Core) QueryByID(ctx context.Context, incomeID string) (Income, error) {
	if err := validate.CheckID(incomeID); err != nil {
		return Income{}, ErrInvalidID
	}

	dbInc, err := c.store.QueryByID(ctx, incomeID)
	if err != nil {
		if errors.Is(err, database.ErrDBNotFound) {
			return Income{}, ErrNotFound
		}
		return Income{}, fmt.Errorf("query: %w", err)
	}

	return toIncome(dbInc), nil
}

// QueryByUserID finds the products identified by a given User ID.
func (c Core) QueryByUserID(ctx context.Context, userID string) ([]Income, error) {
	if err := validate.CheckID(userID); err != nil {
		return nil, ErrInvalidID
	}

	dbInc, err := c.store.QueryByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}

	return toIncomeSlice(dbInc), nil
}
