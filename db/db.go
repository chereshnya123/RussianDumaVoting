package db

import (
	"database/sql"
	"embed"
	"fmt"
	"io/fs"
	"log"
	"path/filepath"
	"sort"
	"time"

	_ "github.com/lib/pq"
)

type Database struct {
	db *sql.DB
}

//go:embed migrations/*.sql
var migrationFiles embed.FS

func NewDatabase(dbPath string) (*Database, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	database := &Database{db: db}

	if err := database.ApplyMigrations(); err != nil {
		return nil, fmt.Errorf("failed to create tables: %w", err)
	}

	return database, nil
}

func (d *Database) Close() error {
	return d.db.Close()
}

func (d *Database) ApplyMigrations() error {
	entries, err := fs.Glob(migrationFiles, "migrations/*.sql")
	if err != nil {
		return fmt.Errorf("failed to read migrations: %w", err)
	}

	sort.Strings(entries)

	for _, entry := range entries {
		filename := filepath.Base(entry)
		log.Printf("Applying migration: %s", filename)

		content, err := migrationFiles.ReadFile(entry)
		if err != nil {
			return fmt.Errorf("failed to read %s: %w", filename, err)
		}

		_, err = d.db.Exec(string(content))
		if err != nil {
			return fmt.Errorf("failed to execute %s: %w", filename, err)
		}
	}

	log.Println("All migrations applied successfully")
	return nil
}

func (d *Database) GetLastUpdateTime() (time.Time, error) {
	var lastUpdate time.Time

	err := d.db.QueryRow("SELECT MAX(last_successful_update) FROM sync_status").Scan(&lastUpdate)

	if err != nil {
		return time.Time{}, err
	}

	return lastUpdate, nil
}
