package db

import "time"

// Deputy represents a row in the "deputies" table.
//
// Columns:
//
//	id        – INTEGER PRIMARY KEY AUTOINCREMENT (internal row id)
//	api_id    – INTEGER UNIQUE (Duma API deputy ID, e.g. 99100142)
//	full_name – TEXT
//	faction   – INTEGER (FK to factions.id)
type Deputy struct {
	Id        int64  `db:"id" json:"id"`
	ApiId     int64  `db:"api_id" json:"api_id"`
	FullName  string `db:"full_name" json:"full_name"`
	FactionId int64  `db:"faction" json:"faction_id"`
}

// Faction represents a row in the "factions" table.
//
// Columns:
//
//	id       – INTEGER PRIMARY KEY
//	api_id   – INTEGER UNIQUE
//	name     – TEXT
//	head     – INTEGER (FK to deputies.id)
type Faction struct {
	Id     int64  `db:"id" json:"id"`
	ApiId  int64  `db:"api_id" json:"api_id"`
	Name   string `db:"name" json:"name"`
	HeadId int64  `db:"head" json:"head_id"`
}

// LawDraft represents a row in the "law_drafts" table.
//
// Columns:
//
//	id      – INTEGER PRIMARY KEY
//	name    – TEXT NOT NULL
//	number  – TEXT NOT NULL
type LawDraft struct {
	Id     int64  `db:"id" json:"id"`
	Name   string `db:"name" json:"name"`
	Number string `db:"number" json:"number"`
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

// VoteStageFactionBulletin represents a row in the "vote_stage_votes" table.
//
// Columns:
//
//	id              – INTEGER PRIMARY KEY
//	for_count       – INTEGER
//	against_count   – INTEGER
//	abstained_count – INTEGER
//	no_vote_count   – INTEGER
//	faction_id      – INTEGER (FK to factions.id)
//	vote_stage_id   – INTEGER (FK to vote_stages.id)
//
// Non-persisted fields used during conversion from API response:
type VoteStageFactionBulletin struct {
	Id             int64  `db:"id" json:"id"`
	ForCount       int64  `db:"for_count" json:"for_count"`
	AgainstCount   int64  `db:"against_count" json:"against_count"`
	AbstainedCount int64  `db:"abstained_count" json:"abstained_count"`
	NoVoteCount    int64  `db:"no_vote_count" json:"no_vote_count"`
	FactionId      int64  `db:"faction_id" json:"faction_id"`
	VoteStageId    int64  `db:"vote_stage_id" json:"vote_stage_id"`
	FactionCode    string `db:"-" json:"-"` // API faction code (used to look up faction_id)
	FactionName    string `db:"-" json:"-"` // Faction full name
	FactionAbbr    string `db:"-" json:"-"` // Faction abbreviation
}
