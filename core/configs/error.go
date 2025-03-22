package configs

import (
	"errors"
	"github.com/golang-migrate/migrate/v4"
)

const (
	errorFailedOpenConnection       = "open connection with DB error: %v"
	errorFailedPingConnection       = "ping connection with DB error: %v"
	errorFailedCloseDatabase        = "failed to close database: %v"
	errorDatabaseHasOpenConnections = "database still has %d open connections"
)

var (
	errorInvalidPath = errors.New("invalid path")
)

func isNonNilAndNotExpectedMigrationError(err error) bool {
	return err != nil && isNoMigration(err) && isMigrationNoChange(err)
}

func isMigrationNoChange(err error) bool {
	return err != nil && errors.Is(err, migrate.ErrNoChange)
}

func isNoMigration(err error) bool {
	return err != nil && errors.Is(err, migrate.ErrNilVersion)
}
