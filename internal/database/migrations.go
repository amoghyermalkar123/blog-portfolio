// internal/database/migrations.go
package database

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

type Migration struct {
	ID      int
	Name    string
	Content string
}

func (db *Database) RunMigrations() error {
	// Create migrations table if it doesn't exist
	_, err := db.Exec(`
        CREATE TABLE IF NOT EXISTS schema_migrations (
            version INTEGER PRIMARY KEY,
            name TEXT NOT NULL,
            applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
        )
    `)
	if err != nil {
		return fmt.Errorf("failed to create migrations table: %v", err)
	}

	// Get all migration files
	migrations, err := loadMigrationFiles()
	if err != nil {
		return fmt.Errorf("failed to load migration files: %v", err)
	}

	// Begin transaction
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %v", err)
	}
	defer tx.Rollback()

	// Run each migration
	for _, migration := range migrations {
		applied, err := isMigrationApplied(tx, migration.ID)
		if err != nil {
			return fmt.Errorf("failed to check migration status: %v", err)
		}

		if !applied {
			db.logger.Info("Applying migration:", migration.Name)

			// Execute migration
			if _, err := tx.Exec(migration.Content); err != nil {
				return fmt.Errorf("failed to apply migration %s: %v", migration.Name, err)
			}

			// Record migration
			if err := recordMigration(tx, migration); err != nil {
				return fmt.Errorf("failed to record migration %s: %v", migration.Name, err)
			}
		}
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit migrations: %v", err)
	}

	return nil
}

func loadMigrationFiles() ([]Migration, error) {
	migrationsDir := "./migrations"
	entries, err := os.ReadDir(migrationsDir)
	if err != nil {
		return nil, err
	}

	var migrations []Migration
	for _, entry := range entries {
		if !entry.IsDir() && filepath.Ext(entry.Name()) == ".up.sql" {
			content, err := os.ReadFile(filepath.Join(migrationsDir, entry.Name()))
			if err != nil {
				return nil, err
			}

			// Parse migration ID from filename
			var id int
			_, err = fmt.Sscanf(entry.Name(), "%d_", &id)
			if err != nil {
				return nil, fmt.Errorf("invalid migration filename format: %s", entry.Name())
			}

			migrations = append(migrations, Migration{
				ID:      id,
				Name:    entry.Name(),
				Content: string(content),
			})
		}
	}

	return migrations, nil
}

func isMigrationApplied(tx *sql.Tx, version int) (bool, error) {
	var exists bool
	err := tx.QueryRow("SELECT EXISTS(SELECT 1 FROM schema_migrations WHERE version = ?)", version).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func recordMigration(tx *sql.Tx, migration Migration) error {
	_, err := tx.Exec("INSERT INTO schema_migrations (version, name) VALUES (?, ?)",
		migration.ID, migration.Name)
	return err
}

func (db *Database) RollbackMigration() error {
	// Get the last applied migration
	var version int
	var name string
	err := db.QueryRow(`
        SELECT version, name 
        FROM schema_migrations 
        ORDER BY version DESC 
        LIMIT 1
    `).Scan(&version, &name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("no migrations to rollback")
		}
		return fmt.Errorf("failed to get last migration: %v", err)
	}

	// Read the down migration file
	downFile := fmt.Sprintf("./migrations/%d_%s.down.sql", version, name)
	content, err := os.ReadFile(downFile)
	if err != nil {
		return fmt.Errorf("failed to read down migration: %v", err)
	}

	// Begin transaction
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %v", err)
	}
	defer tx.Rollback()

	// Execute down migration
	if _, err := tx.Exec(string(content)); err != nil {
		return fmt.Errorf("failed to execute down migration: %v", err)
	}

	// Remove migration record
	if _, err := tx.Exec("DELETE FROM schema_migrations WHERE version = ?", version); err != nil {
		return fmt.Errorf("failed to remove migration record: %v", err)
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit rollback: %v", err)
	}

	return nil
}
