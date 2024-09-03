package db

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
)

// Migrate applies all pending SQL migrations from the migrations directory to the database.
// It reads each SQL file and executes its contents against the provided database connection.
func Migrate(db *sql.DB) error {
	// Find all migration files in the specified directory.
	migrationFiles, err := filepath.Glob("./db/migrations/*.sql")
	if err != nil {
		return fmt.Errorf("failed to list migration files: %w", err)
	}

	// Apply each migration file to the database.
	for _, file := range migrationFiles {
		if err := applyMigration(db, file); err != nil {
			return fmt.Errorf("error applying migration %s: %w", file, err)
		}
	}

	return nil
}

// applyMigration reads the SQL file and executes its contents against the database.
// It is called by the Migrate function for each migration file.
func applyMigration(db *sql.DB, filePath string) error {
	// Read the content of the SQL migration file.
	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read migration file %s: %w", filePath, err)
	}

	// Execute the SQL commands in the migration file.
	_, err = db.Exec(string(data))
	if err != nil {
		return fmt.Errorf("failed to execute migration file %s: %w", filePath, err)
	}

	return nil
}
