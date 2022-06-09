package db

import "time"

// Expense represents an individual expense.
type Expense struct {
	ID               string    `db:"expense_id"`        // Unique identifier.
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
