package db

import (
	"database/sql"
	"embed"
	"errors"
	"fmt"
	"io/fs"
	"log"
	"path/filepath"
	"sort"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// Database provides access to the SQLite database.
type Database struct {
	db *sql.DB
}

//go:embed migrations/*.sql
var migrationFiles embed.FS

// NewDatabase opens the SQLite database and applies migrations.
func NewDatabase(dbPath string) (*Database, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	database := &Database{db: db}

	if err := database.ApplyMigrations(); err != nil {
		return nil, fmt.Errorf("failed to apply migrations: %w", err)
	}

	return database, nil
}

// Close closes the database connection.
func (d *Database) Close() error {
	return d.db.Close()
}

// ApplyMigrations applies all SQL migration files in the embedded migrations directory.
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

		if _, err := d.db.Exec(string(content)); err != nil {
			return fmt.Errorf("failed to execute %s: %w", filename, err)
		}
	}

	log.Println("All migrations applied successfully")
	return nil
}

// GetLastUpdateTime returns the last successful update time.
// If sync_status is empty, returns Unix epoch (time zero).
func (d *Database) GetLastUpdateTime() (time.Time, error) {
	var lastUpdate time.Time
	var rowsCount int
	_ = d.db.QueryRow("SELECT COUNT(last_successful_update) FROM sync_status").Scan(&rowsCount)
	if rowsCount == 0 {
		return time.Unix(0, 0), nil
	}
	err := d.db.QueryRow("SELECT MAX(last_successful_update) FROM sync_status").Scan(&lastUpdate)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return time.Time{}, fmt.Errorf("failed to get last update time: %w", err)
	}

	if errors.Is(err, sql.ErrNoRows) {
		return time.Unix(0, 0), nil
	}
	return lastUpdate, nil
}

func (d *Database) GetLatestVoting() (Voting, error) {
	var v Voting
	query := `SELECT id, name, date, question_id, factions, result
	          FROM votings ORDER BY date DESC LIMIT 1`
	if err := d.db.QueryRow(query).Scan(
		&v.Id, &v.Name, &v.Date, &v.QuestionId, &v.Factions, &v.Result,
	); err != nil {
		return Voting{}, fmt.Errorf("get latest voting: %w", err)
	}
	return v, nil
}

func (d *Database) GetQuestionByID(questionId int64) (Question, error) {
	var q Question
	query := `SELECT id, name, tags, votings_id, profile_committee_id,
	                  responsible_committee_id, other_committees, authors
	          FROM questions WHERE id = ?`
	if err := d.db.QueryRow(query, questionId).Scan(
		&q.Id, &q.Name, &q.Tags, &q.VotingsId,
		&q.ProfileCommitteeId, &q.ResponsibleCommitteeId,
		&q.OtherCommittees, &q.Authors,
	); err != nil {
		return Question{}, fmt.Errorf("get question by id %d: %w", questionId, err)
	}
	return q, nil
}

// SaveDeputy inserts a deputy into the database. SQLite auto-generates the internal
// rowid (ID); the Deputy.APIID field is stored in the api_id column.
func (d *Database) SaveDeputy(deputy *Deputy) error {
	query := `INSERT INTO deputies (api_id, full_name, faction, department)
	          VALUES (?, ?, ?, ?)`
	result, err := d.db.Exec(query, deputy.ApiId, deputy.FullName, deputy.FactionId, deputy.Department)
	if err != nil {
		return fmt.Errorf("insert deputy api_id=%d: %w", deputy.ApiId, err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("get insert ID: %w", err)
	}
	deputy.Id = id

	return nil
}

// SaveDeputyUpsert inserts a deputy or updates it if the api_id already exists.
func (d *Database) SaveDeputyUpsert(deputy *Deputy) error {
	query := `INSERT INTO deputies (api_id, full_name, faction, department)
	          VALUES (?, ?, ?, ?)
	          ON CONFLICT(api_id) DO UPDATE SET
	            full_name = excluded.full_name,
	            faction   = excluded.faction,
	            department = excluded.department`
	if _, err := d.db.Exec(query, deputy.ApiId, deputy.FullName, deputy.FactionId, deputy.Department); err != nil {
		return fmt.Errorf("upsert deputy api_id=%d: %w", deputy.ApiId, err)
	}

	return nil
}

func (d *Database) isFactionExists(factionApiId int64) bool {
	query := `SELECT COUNT(*) FROM factions where api_id == ?`
	var count int64
	err := d.db.QueryRow(query, factionApiId).Scan(&count)
	if err != nil && err != sql.ErrNoRows {
		panic(fmt.Sprintf("can not check if faction exists. Faction api_id = %d, err = %s", factionApiId, err.Error()))
	}

	return count != 0
}

// SaveFaction inserts a faction into the database.
func (d *Database) SaveFaction(faction *Faction) error {
	if d.isFactionExists(faction.ApiId) {
		return nil
	}

	query := `INSERT INTO factions (api_id, name, head)
						VALUES (?, ?, ?)`
	result, err := d.db.Exec(query, faction.ApiId, faction.Name, faction.HeadId)
	if err != nil {
		return fmt.Errorf("insert faction api_id=%d: %w", faction.ApiId, err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("get faction insert ID: %w", err)
	}
	faction.Id = id

	return nil
}
