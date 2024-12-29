// internal/database/database.go
package database

import (
	"blog-portfolio/internal/logger"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

type Database struct {
	*sql.DB
	logger *logger.Logger
}

func New(logger *logger.Logger) (*Database, error) {
	// Ensure database directory exists
	dbDir := "./data"
	if err := os.MkdirAll(dbDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create database directory: %v", err)
	}

	// Connect to SQLite database
	dbPath := filepath.Join(dbDir, "blog.db")
	db, err := sql.Open("sqlite3", dbPath+"?_foreign_keys=on")
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %v", err)
	}

	// Test the connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}

	// Set connection pool settings
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)

	return &Database{
		DB:     db,
		logger: logger,
	}, nil
}

// Close closes the database connection
func (db *Database) Close() error {
	return db.DB.Close()
}
