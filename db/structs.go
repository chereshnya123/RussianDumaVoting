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

// VoteStage represents a row in the "vote_stages" table.
//
// Columns:
//
//	id               – INTEGER PRIMARY KEY AUTOINCREMENT
//	law_draft_id     – INTEGER (FK to law_drafts.id)
//	stage_name       – TEXT
//	for_count        – INTEGER
//	against_count    – INTEGER
//	abstained_count  – INTEGER
//	no_vote_count    – INTEGER
type VoteStage struct {
	Id             int64  `db:"id" json:"id"`
	LawDraftId     int64  `db:"law_draft_id" json:"law_draft_id"`
	StageName      string `db:"stage_name" json:"stage_name"`
	ForCount       int    `db:"for_count" json:"for_count"`
	AgainstCount   int    `db:"against_count" json:"against_count"`
	AbstainedCount int    `db:"abstained_count" json:"abstained_count"`
	NoVoteCount    int    `db:"no_vote_count" json:"no_vote_count"`
}

// VoteStageBullentin represents a row in the "vote_stage_bulletins" table.
//
// Columns:
//
//	id            – INTEGER PRIMARY KEY
//	verdict       – VARCHAR(10) CHECK(in ('FOR', 'AGAINST', 'ABSTAINED', 'NO_VOTE'))
//	deputy_id     – INTEGER (FK to deputies.id)
//	vote_stage_id – INTEGER (FK to vote_stages.id)
type VoteStageBulletin struct {
	Id          int64  `db:"id" json:"id"`
	Verdict     string `db:"verdict" json:"verdict"`
	DeputyId    int64  `db:"deputy_id" json:"deputy_id"`
	VoteStageId int64  `db:"vote_stage_id" json:"vote_stage_id"`
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
