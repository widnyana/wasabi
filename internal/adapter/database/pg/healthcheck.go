package pg

import (
	"context"
	"database/sql"
)

// HealthChecker is a struct that holds the database client and provides a method to check the health of the database connection.
type HealthChecker struct {
	client *sql.DB
}

// NewHealthChecker creates a new HealthChecker instance with the provided database client.
// It returns a HealthChecker that can be used to check the health of the database connection.
func NewHealthChecker(client *sql.DB) HealthChecker {
	return HealthChecker{client}
}

// CheckHealth checks the health of the database connection by pinging the database.
// It returns an error if the ping fails, indicating a problem with the database connection.
func (healthChecker HealthChecker) CheckHealth(ctx context.Context) error {
	return healthChecker.client.PingContext(ctx)
}
