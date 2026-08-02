package analyzer

import (
	"dumaVote/db"
	"fmt"
	"log/slog"
	"strconv"
	"time"
)

type Updater struct {
	db      *db.Database
	fetcher Fetcher
	// parser  Parser
	logger *slog.Logger
}

func NewUpdater(appApiKey, personApiKey string, db *db.Database, logger *slog.Logger) *Updater {
	return &Updater{
		db:      db,
		fetcher: *NewFetcher(appApiKey, personApiKey, logger),
		logger:  logger,
	}
}

var unixEpochStart = time.Unix(0, 0)

func (u *Updater) ShouldUpdate() (bool, error) {
	lastUpdate, err := u.db.GetLastUpdateTime()
	if err != nil {
		return false, err
	}

	if lastUpdate.Equal(unixEpochStart) {
		return true, nil
	}

	return time.Since(lastUpdate) > time.Hour, nil
}

func (u *Updater) UpdateDeputies() error {
	deputies, err := u.fetcher.FetchAllDeputies()
	if err != nil {
		return err
	}

	for _, rawDeputy := range deputies {
		if rawDeputy.Position == "Член СФ" || !rawDeputy.IsCurrent {
			continue
		}
		u.logger.Info("Update deputy. ", "Id", rawDeputy.Id)
		deputyInfo, err := u.fetcher.FetchDeputyInfo(rawDeputy.Id)
		if err != nil {
			return err
		}

		apiId, err := strconv.ParseInt(rawDeputy.Id, 10, 64)
		if err != nil {
			return fmt.Errorf("cannot parse deputy api ID %q: %w", rawDeputy.Id, err)
		}

		var deputy db.Deputy
		deputy.ApiId = apiId

		if len(deputyInfo.Activities) > 0 {
			deputy.Department = deputyInfo.Activities[0].SubdivisionId
		}
		if deputyInfo.FactionId == "" {
			deputyInfo.FactionId = "-1"
		}

		factionId, err := strconv.ParseInt(deputyInfo.FactionId, 10, 64)
		if err != nil {
			return fmt.Errorf("cannot parse faction ID %q: %w", deputyInfo.FactionId, err)
		}
		deputy.FactionId = factionId
		deputy.FullName = deputyInfo.Name

		if err := u.db.SaveDeputyUpsert(&deputy); err != nil {
			return fmt.Errorf("cannot save deputy. API id = %d: %w", apiId, err)
		}
	}

	return nil
}

// UpdateData refreshes the data if enough time has elapsed.
func (u *Updater) UpdateDatabase() error {
	update, err := u.ShouldUpdate()
	if err != nil || !update {
		return err
	}
	return u.UpdateDeputies()
}
