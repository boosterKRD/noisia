package db

import (
	"context"
)

// TestConninfo defines connection string to test database.
const TestConninfo = "host=10.233.98.28 port=1521 user=noisia password=123 database=noisia_fixtures"

// NewTestDB creates connection for test database.
func NewTestDB() (DB, error) {
	return NewPostgresDB(context.Background(), TestConninfo)
}
