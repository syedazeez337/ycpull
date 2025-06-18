package main

import (
	"database/sql"
	"fmt"
	"strings"

	_ "github.com/mattn/go-sqlite3" // SQLite driver
)

// InitDB initializes the SQLite database and creates the startups table if it doesn't exist.
func InitDB(dbPath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database at %s: %w", dbPath, err)
	}

	createTableSQL := `
	CREATE TABLE IF NOT EXISTS startups (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT,
		slug TEXT UNIQUE,
		description TEXT,
		batch TEXT,
		logo TEXT,
		website TEXT,
		tags TEXT,
		location TEXT
	);`

	_, err = db.Exec(createTableSQL)
	if err != nil {
		return nil, fmt.Errorf("failed to create startups table: %w", err)
	}

	return db, nil
}

// StoreStartups inserts a slice of Startup objects into the database.
// It converts the Tags slice into a comma-separated string for storage.
func StoreStartups(db *sql.DB, startups []Startup) error {
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	stmt, err := tx.Prepare(`
		INSERT INTO startups (name, slug, description, batch, logo, website, tags, location)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(slug) DO NOTHING;`) // Avoid duplicates based on slug
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to prepare insert statement: %w", err)
	}
	defer stmt.Close()

	for _, s := range startups {
		tagsStr := strings.Join(s.Tags, ",")
		_, err := stmt.Exec(s.Name, s.Slug, s.Description, s.Batch, s.Logo, s.Website, tagsStr, s.Location)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to insert startup %s (slug: %s): %w", s.Name, s.Slug, err)
		}
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
