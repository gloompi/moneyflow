// Package expense provides a core business API. Right now these
package expense

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/gloompi/ultimate-service/business/core/expense/db"
	"github.com/gloompi/ultimate-service/business/sys/database"
	"github.com/gloompi/ultimate-service/business/sys/validate"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

// Set of error variables for CRUD operations.
var (
	ErrNotFound  = errors.New("expense not found")
	ErrInvalidID = errors.New("ID is not in its proper form")
)

// Core manages the set of APIs for expense access.
type Core struct {
	store db.Store
}

// NewCore constructs a core for expense api access.
func NewCore(log *zap.SugaredLogger, sqlxDB *sqlx.DB) Core {
	return Core{
		store: db.NewStore(log, sqlxDB),
	}
}

// Create adds an Expense to the database. It returns the created Expense with
// fields like ID and DateCreated populated.
func (c Core) Create(ctx context.Context, ne NewExpense, now time.Time) (Expense, error) {
	if err := validate.Check(ne); err != nil {
		return Expense{}, fmt.Errorf("validating data: %w", err)
	}

	dbExp := db.Expense{
		ID:               validate.GenerateID(),
		Name:             ne.Name,
		Category:         ne.Category,
		Currency:         ne.Currency,
		Amount:           ne.Amount,
		Reoccurrence:     ne.Reoccurrence,
		Duration:         ne.Duration,
		ReoccurrenceType: ne.ReoccurrenceType,
		DurationType:     ne.DurationType,
		UserID:           ne.UserID,
		DateCreated:      now,
		DateUpdated:      now,
	}

	if err := c.store.Create(ctx, dbExp); err != nil {
		return Expense{}, fmt.Errorf("create: %w", err)
	}

	return toExpense(dbExp), nil
}

// Update modifies data about a Expense. It will error if the specified ID is
// invalid or does not reference an existing Expense.
func (c Core) Update(ctx context.Context, expenseID string, ue UpdateExpense, now time.Time) error {
	if err := validate.CheckID(expenseID); err != nil {
		return ErrInvalidID
	}

	if err := validate.Check(ue); err != nil {
		return fmt.Errorf("validating data: %w", err)
	}

	dbExp, err := c.store.QueryByID(ctx, expenseID)
	if err != nil {
		if errors.Is(err, database.ErrDBNotFound) {
			return ErrNotFound
		}
		return fmt.Errorf("updating expense expenseID[%s]: %w", expenseID, err)
	}

	if ue.Name != nil {
		dbExp.Name = *ue.Name
	}
	if ue.Category != nil {
		dbExp.Category = *ue.Category
	}
	if ue.Currency != nil {
		dbExp.Currency = *ue.Currency
	}
	if ue.Amount != nil {
		dbExp.Amount = *ue.Amount
	}
	if ue.Reoccurrence != nil {
		dbExp.Reoccurrence = *ue.Reoccurrence
	}
	if ue.Duration != nil {
		dbExp.Duration = *ue.Duration
	}
	if ue.ReoccurrenceType != nil {
		dbExp.ReoccurrenceType = *ue.ReoccurrenceType
	}
	if ue.DurationType != nil {
		dbExp.DurationType = *ue.DurationType
	}
	dbExp.DateUpdated = now

	if err := c.store.Update(ctx, dbExp); err != nil {
		return fmt.Errorf("update: %w", err)
	}

	return nil
}

// Delete removes the expense identified by a given ID.
func (c Core) Delete(ctx context.Context, expenseID string) error {
	if err := validate.CheckID(expenseID); err != nil {
		return ErrInvalidID
	}

	if err := c.store.Delete(ctx, expenseID); err != nil {
		return fmt.Errorf("delete: %w", err)
	}

	return nil
}

// Query gets all Expense from the database.
func (c Core) Query(ctx context.Context, pageNumber int, rowsPerPage int) ([]Expense, error) {
	dbExps, err := c.store.Query(ctx, pageNumber, rowsPerPage)
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}

	return toExpenseSlice(dbExps), nil
}

// QueryByID finds the expense identified by a given ID.
func (c Core) QueryByID(ctx context.Context, expenseID string) (Expense, error) {
	if err := validate.CheckID(expenseID); err != nil {
		return Expense{}, ErrInvalidID
	}

	dbExp, err := c.store.QueryByID(ctx, expenseID)
	if err != nil {
		if errors.Is(err, database.ErrDBNotFound) {
			return Expense{}, ErrNotFound
		}
		return Expense{}, fmt.Errorf("query: %w", err)
	}

	return toExpense(dbExp), nil
}

// QueryByUserID finds the expenses identified by a given User ID.
func (c Core) QueryByUserID(ctx context.Context, userID string) ([]Expense, error) {
	if err := validate.CheckID(userID); err != nil {
		return nil, ErrInvalidID
	}

	dbInc, err := c.store.QueryByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}

	return toExpenseSlice(dbInc), nil
}
