package db

import "time"

// Income represents an individual income.
type Income struct {
	ID               string    `db:"income_id"`         // Unique identifier.
	Name             string    `db:"name"`              // Display name of the income.
	Category         string    `db:"category"`          // Category of the income.
	Currency         string    `db:"currency"`          // Currency of the income.
	Amount           int       `db:"amount"`            // Amount of money of the income.
	Reoccurrence     int       `db:"reoccurrence"`      // Execute transaction each day, week, month.
	Duration         int       `db:"duration"`          // Range of time transaction needs to be happening.
	ReoccurrenceType string    `db:"reoccurrence_type"` // Type of reoccurrence (Monthly, Daily, Once).
	DurationType     string    `db:"duration_type"`     // Type of duration (Months, Days).
	UserID           string    `db:"user_id"`           // ID of the user who created the income.
	DateCreated      time.Time `db:"date_created"`      // When the income was added.
	DateUpdated      time.Time `db:"date_updated"`      // When the income record was last modified.
}

// Income represents an income structure by user.
type IncomeByUser struct {
	Income
	Total int `db:"total"` // Total income.
}
