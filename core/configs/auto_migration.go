package configs

import (
	"database/sql"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	pg "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/winartodev/apollo/core"
	"github.com/winartodev/apollo/core/helpers"
	"log"
)

const (
	migrationsDir = "core/database/migration"
)

type SchemaMigration struct {
	version int
	dirty   bool
}

type AutoMigration struct {
	migrate *migrate.Migrate
}

func NewAutoMigration(databaseName string, db *sql.DB) (*AutoMigration, error) {
	sourceURL, err := generateSourceURL()
	if err != nil {
		return nil, fmt.Errorf("failed to generate source URL: %w", err)
	}

	x, err := pg.WithInstance(db, &pg.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to create database driver instance: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(*sourceURL, databaseName, x)
	if err != nil {
		return nil, fmt.Errorf("failed to get schema version after migration: %w", err)
	}

	log.Printf("AutoMigration initialized successfully for database: %s", databaseName)
	return &AutoMigration{
		migrate: m,
	}, nil
}

func (am *AutoMigration) Start() error {
	err := am.migrate.Up()
	if err == nil {
		schema, err := am.getSchemaMigration()
		if err != nil {
			return fmt.Errorf("failed to get schema version after migration error: %w", err)
		}

		log.Printf("Migration successful. Current schema version: %d", schema.version)
		return nil
	}

	if isErrorNoChange(err) {
		log.Printf("No migrations to run. Database is already at the latest version.")
		return nil
	}

	// Handle if migrate failed
	schema, err := am.getSchemaMigration()
	if err != nil {
		return fmt.Errorf("failed to get schema version after migration error: %w", err)
	}

	log.Printf("Migration failed. Current schema version: %d, dirty: %v", schema.version, schema.dirty)

	err = am.Fix(schema.version)
	if err != nil {
		return fmt.Errorf("failed to fix schema version: %w", err)
	}

	err = am.Rollback()
	if err != nil {
		return fmt.Errorf("failed to rollback migration: %w", err)
	}

	log.Printf("Migration recovery completed successfully.")

	return nil
}

func (am *AutoMigration) Fix(version int) error {
	log.Printf("Fixing schema version to: %d", version)

	err := am.migrate.Force(version)
	if err != nil {
		return fmt.Errorf("failed to force schema version: %w", err)
	}

	schema, err := am.getSchemaMigration()
	if err != nil {
		return fmt.Errorf("failed to get schema version after fix: %w", err)
	}

	log.Printf("Schema fix successful. Current version: %d, dirty: %v", schema.version, schema.dirty)
	return nil
}

func (am *AutoMigration) Rollback() error {
	log.Printf("Rolling back schema by 1 step.")

	err := am.migrate.Steps(-1)
	if err != nil {
		return fmt.Errorf("failed to rollback schema: %w", err)
	}

	schema, err := am.getSchemaMigration()
	if err != nil {
		return fmt.Errorf("failed to get schema version after rollback: %w", err)
	}

	log.Printf("Rollback successful. Current version: %d, dirty: %v", schema.version, schema.dirty)

	return nil
}

func (am *AutoMigration) getSchemaMigration() (*SchemaMigration, error) {
	version, dirty, err := am.migrate.Version()
	if isNonNilAndNotExpectedMigrationError(err) {
		return nil, fmt.Errorf("failed to retrieve schema version: %w", err)
	}

	return &SchemaMigration{
		version: int(version),
		dirty:   dirty,
	}, nil
}

func generateSourceURL() (*string, error) {
	filePath, err := helpers.GetCompletePath(migrationsDir)
	if err != nil {
		return nil, errorInvalidPath
	}

	var completePath string
	if helpers.CurrentOS(core.OSWindows) {
		completePath = fmt.Sprintf("file:%s", filePath)
	} else {
		completePath = fmt.Sprintf("file://%s", filePath)
	}

	return &completePath, nil
}
