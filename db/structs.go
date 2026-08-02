package db

import "time"

// Question represents a row in the "questions" table.
//
// Columns:
//
//	id            – INTEGER PRIMARY KEY
//	name          – TEXT NOT NULL
//	tags          – TEXT (JSON array)
//	votings_id    – TEXT (JSON array of voting IDs)
//	profile_committee_id  – INTEGER
//	responsible_committee_id – INTEGER
//	other_committees      – TEXT (JSON array of committee IDs)
//	authors         – TEXT (JSON array of author IDs)
type Question struct {
	Id                     int64     `db:"id" json:"id"`
	Name                   string    `db:"name" json:"name"`
	Tags                   string    `db:"tags" json:"tags"`
	VotingsId              string    `db:"votings_id" json:"votings_id"`
	ProfileCommitteeId     int64     `db:"profile_committee_id" json:"profile_committee_id"`
	ResponsibleCommitteeId int64     `db:"responsible_committee_id" json:"responsible_committee_id"`
	OtherCommittees        string    `db:"other_committees" json:"other_committees"`
	Authors                string    `db:"authors" json:"authors"`
	CreatedAt              time.Time `db:"created_at" json:"created_at"`
	UpdatedAt              time.Time `db:"updated_at" json:"updated_at"`
}

// Department represents a row in the "departments" table.
//
// Columns:
//
//	id    – INTEGER PRIMARY KEY
//	head  – INTEGER NOT NULL (FK to deputies.id)
//	name  – TEXT
//	size  – INTEGER
//	members – TEXT (JSON array of deputy IDs)
type Department struct {
	Id      int64  `db:"id" json:"id"`
	HeadId  int64  `db:"head" json:"head_id"`
	Name    string `db:"name" json:"name"`
	Size    int    `db:"size" json:"size"`
	Members string `db:"members" json:"members"` // JSON array of deputy IDs
}

// Deputy represents a row in the "deputies" table.
//
// Columns:
//
//	id        – INTEGER PRIMARY KEY AUTOINCREMENT (internal row id)
//	api_id    – INTEGER UNIQUE (Duma API deputy ID, e.g. 99100142)
//	full_name – TEXT
//	faction   – INTEGER (FK to factions.id)
//	department – INTEGER (FK to departments.id)
type Deputy struct {
	Id         int64  `db:"id" json:"id"`
	ApiId      int64  `db:"api_id" json:"api_id"`
	FullName   string `db:"full_name" json:"full_name"`
	FactionId  int64  `db:"faction" json:"faction_id"`
	Department int64  `db:"department" json:"department_id"`
}

// Faction represents a row in the "factions" table.
//
// Columns:
//
//	id             – INTEGER PRIMARY KEY
//	name           – TEXT
//	head           – INTEGER (FK to deputies.id)
//	department     – INTEGER (FK to departments.id)
//	target_questions – TEXT (JSON array)
//	target_tags      – TEXT (JSON map)
type Faction struct {
	Id              int64  `db:"id" json:"id"`
	Name            string `db:"name" json:"name"`
	HeadId          int64  `db:"head" json:"head_id"`
	DepartmentId    int64  `db:"department" json:"department_id"`
	TargetQuestions string `db:"target_questions" json:"target_questions"` // JSON array
	TargetTags      string `db:"target_tags" json:"target_tags"`           // JSON map
}

// Voting represents a row in the "votings" table.
//
// Columns:
//
//	id         – INTEGER PRIMARY KEY
//	name       – TEXT
//	date       – DATETIME
//	question_id – INTEGER (FK to questions.id)
//	factions   – TEXT (JSON array of faction IDs)
//	result     – TEXT (JSON map)
type Voting struct {
	Id         int64     `db:"id" json:"id"`
	Name       string    `db:"name" json:"name"`
	Date       time.Time `db:"date" json:"date"`
	QuestionId int64     `db:"question_id" json:"question_id"`
	Factions   string    `db:"factions" json:"factions"` // JSON array of faction IDs
	Result     string    `db:"result" json:"result"`     // JSON map
}

// SyncStatus represents a row in the "sync_status" table.
//
// Columns:
//
//	id                     – INTEGER PRIMARY KEY
//	last_successful_update – DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
type SyncStatus struct {
	Id                   int64     `db:"id" json:"id"`
	LastSuccessfulUpdate time.Time `db:"last_successful_update" json:"last_successful_update"`
}
