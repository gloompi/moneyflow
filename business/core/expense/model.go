package expense

import (
	"time"
	"unsafe"

	"github.com/gloompi/ultimate-service/business/core/expense/db"
)

// Expense represents an individual expense.
type Expense struct {
	ID               string    `db:"id"`                // Unique identifier.
	Name             string    `db:"name"`              // Display name of the expense.
	Category         string    `db:"category"`          // Category of the expense.
	Currency         string    `db:"currency"`          // Currency of the expense.
	Amount           int       `db:"amount"`            // Amount of money of the expense.
	Reoccurrence     int       `db:"reoccurrence"`      // Execute transaction each day, week, month.
	Duration         int       `db:"duration"`          // Range of time transaction needs to be happening.
	ReoccurrenceType string    `db:"reoccurrence_type"` // Type of reoccurrence (Monthly, Daily, Once).
	DurationType     string    `db:"duration_type"`     // Type of duration (Months, Days).
	UserID           string    `db:"user_id"`           // ID of the user who created the expense.
	DateCreated      time.Time `db:"date_created"`      // When the expense was added.
	DateUpdated      time.Time `db:"date_updated"`      // When the expense record was last modified.
}

// NewExpense is what we require from clients when adding a Expense.
type NewExpense struct {
	Name             string `json:"name" validate:"required"`
	Category         string `json:"category" validate:"required"`
	Currency         string `json:"currency" validate:"required"`
	Amount           int    `json:"amount" validate:"omitempty,gte=1"`
	Reoccurrence     int    `json:"reoccurrence" validate:"omitempty,gte=1"`
	Duration         int    `json:"duration" validate:"omitempty,gte=1"`
	ReoccurrenceType string `json:"reoccurrence_type"`
	DurationType     string `json:"duration_type"`
	UserID           string `json:"user_id" validate:"required"`
}

// UpdateExpense defines what information may be provided to modify an
// existing Expense. All fields are optional so clients can send just the
// fields they want changed. It uses pointer fields so we can differentiate
// between a field that was not provided and a field that was provided as
// explicitly blank. Normally we do not want to use pointers to basic types but
// we make exceptions around marshalling/unmarshalling.
type UpdateExpense struct {
	Name             *string `json:"name"`
	Category         *string `json:"category" validate:"required"`
	Currency         *string `json:"currency" validate:"required"`
	Amount           *int    `json:"amount" validate:"required"`
	Reoccurrence     *int    `json:"reoccurrence"`
	Duration         *int    `json:"duration"`
	ReoccurrenceType *string `json:"reoccurrence_type"`
	DurationType     *string `json:"duration_type"`
}

// =============================================================================

func toExpense(dbExp db.Expense) Expense {
	iu := (*Expense)(unsafe.Pointer(&dbExp))
	return *iu
}

func toExpenseSlice(dbExps []db.Expense) []Expense {
	exps := make([]Expense, len(dbExps))
	for i, dbExp := range dbExps {
		exps[i] = toExpense(dbExp)
	}
	return exps
}
