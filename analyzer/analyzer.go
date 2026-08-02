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

func (a *Analyzer) GetLastQuestion() (db.Question, error) {
	err := a.updater.UpdateDatabase()
	if err != nil {
		a.logger.Error("Cannot get latest voting.", " Error", err)
		return db.Question{}, err
	}

	voting, err := a.db.GetLatestVoting()
	if err != nil {
		a.logger.Error("Cannot get latest voting.", " Error", err)
		return db.Question{}, err
	}

	return a.db.GetQuestionByID(voting.QuestionId)
}
