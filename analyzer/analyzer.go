package analyzer

import (
	"dumaVote/db"
	"log/slog"
)

type Analyzer struct {
	db      *db.Database
	updater *Updater
	logger  *slog.Logger
}

func NewAnalyzer(appApiKey, personApiKey string, db *db.Database, logger *slog.Logger) *Analyzer {
	return &Analyzer{
		db:      db,
		updater: NewUpdater(appApiKey, personApiKey, db, logger),
		logger:  logger,
	}
}

func (a *Analyzer) GetLastVoting() (db.LawDraft, error) {
	err := a.updater.UpdateDatabase()
	if err != nil {
		a.logger.Error("Cannot get latest voting. Can not update database", " Error", err)
		return db.LawDraft{}, err
	}

	return db.LawDraft{}, nil
}
