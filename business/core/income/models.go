package income

import (
	"time"
	"unsafe"

	"github.com/gloompi/ultimate-service/business/core/income/db"
)

// Income represents an individual income.
type Income struct {
	ID               string    `json:"id"`                // Unique identifier.
	Name             string    `json:"name"`              // Display name of the income.
	Category         string    `json:"category"`          // Category of the income.
	Currency         string    `json:"currency"`          // Currency of the income.
	Amount           int       `json:"amount"`            // Amount of money of the income.
	Reoccurrence     int       `json:"reoccurrence"`      // Execute transaction each day, week, month.
	Duration         int       `json:"duration"`          // Range of time transaction needs to be happening.
	ReoccurrenceType string    `json:"reoccurrence_type"` // Type of reoccurrence (Monthly, Daily, Once).
	DurationType     string    `json:"duration_type"`     // Type of duration (Months, Days).
	UserID           string    `json:"user_id"`           // ID of the user who created the income.
	DateCreated      time.Time `json:"date_created"`      // When the income was added.
	DateUpdated      time.Time `json:"date_updated"`      // When the income record was last modified.
}

// NewIncome is what we require from clients when adding a Income.
type NewIncome struct {
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

// UpdateProduct defines what information may be provided to modify an
// existing Product. All fields are optional so clients can send just the
// fields they want changed. It uses pointer fields so we can differentiate
// between a field that was not provided and a field that was provided as
// explicitly blank. Normally we do not want to use pointers to basic types but
// we make exceptions around marshalling/unmarshalling.
type UpdateIncome struct {
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

func toIncome(dbInc db.Income) Income {
	iu := (*Income)(unsafe.Pointer(&dbInc))
	return *iu
}

func toIncomeSlice(dbIncs []db.Income) []Income {
	incs := make([]Income, len(dbIncs))
	for i, dbInc := range dbIncs {
		incs[i] = toIncome(dbInc)
	}
	return incs
}
