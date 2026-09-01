package analyzer

import (
	"dumaVote/db"
	"encoding/json"
	"errors"
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

func (u *Updater) getActualFaction(factions []Faction) (Faction, error) {
	if len(factions) == 0 {
		return Faction{}, errors.New("can not get actual faction from empty list")
	}

	latest := factions[0]
	for _, f := range factions[1:] {
		if f.EndDate > latest.EndDate {
			latest = f
		}
	}

	return latest, nil
}

func (u *Updater) getActualFactionId(deputy Deputy) (int64, error) {
	currentFaction, err := u.getActualFaction(deputy.Factions)
	if err != nil {
		u.logger.Error("cannot get actual faction for deputy.", "deputyId", deputy.Id)
		return -1, err
	}
	currentFactionId, err := strconv.Atoi(currentFaction.Id)
	if err != nil {
		u.logger.Error("cannot parse current faction id to int.", "factionId", currentFaction.Id)
		return -1, err
	}

	return int64(currentFactionId), nil
}

func (u *Updater) updateFactions(factions []Faction) error {
	currentFaction, err := u.getActualFaction(factions)
	if err != nil {
		return nil
	}

	currentFactionId, err := strconv.Atoi(currentFaction.Id)
	if err != nil {
		u.logger.Error("cannot parse current faction id to int.", "factionId", currentFaction.Id)
		return err
	}

	var faction db.Faction
	faction.ApiId = int64(currentFactionId)
	faction.Name = currentFaction.Name

	err = u.db.SaveFaction(&faction)
	if err != nil {
		u.logger.Error("cannot save faction in database.", "factionId", faction.ApiId)
		return err
	}

	return nil
}

func (u *Updater) UpdateDeputiesAndFactions() error {
	deputies, err := u.fetcher.FetchAllDeputies()
	if err != nil {
		return err
	}

	factions := make(map[int64]db.Faction)
	for _, rawDeputy := range deputies {
		if rawDeputy.Position == "Член СФ" || !rawDeputy.IsCurrent {
			continue
		}
		u.logger.Info("Fetch deputy. ", "Id", rawDeputy.Id, "name", rawDeputy.Name)

		apiId, err := strconv.ParseInt(rawDeputy.Id, 10, 64)
		if err != nil {
			return fmt.Errorf("cannot parse deputy api ID %q: %w", rawDeputy.Id, err)
		}

		currentFactionId, err := u.getActualFactionId(rawDeputy)
		if err != nil {
			return err
		}
		var deputy db.Deputy
		deputy.ApiId = apiId
		deputy.FullName = rawDeputy.Name
		deputy.FactionId = int64(currentFactionId)
		deputy.Department = -1

		err = u.updateFactions(rawDeputy.Factions)
		if err != nil {
			u.logger.Error("cannot update factions.", "deputyId", rawDeputy.Id)
			return err
		}

		if err := u.db.SaveDeputyUpsert(&deputy); err != nil {
			return fmt.Errorf("cannot save deputy. API id = %d: %w", apiId, err)
		}
		deputyJson, err := json.Marshal(deputy)
		if err != nil {
			return fmt.Errorf("cannot marshall deputy. Error = %v", err)
		}
		u.logger.Debug("Update deputy. ", "deputy", string(deputyJson))

	}

	for _, faction := range factions {
		err = u.db.SaveFaction(&faction)
		if err != nil {
			return fmt.Errorf("can not save faction. Id = %d, name = %s. Error = %v", faction.ApiId, faction.Name, err)
		}
		u.logger.Info("Save faction. ", "id", faction.ApiId, "name", faction.Name)
	}

	return nil
}

// UpdateData refreshes the data if enough time has elapsed.
func (u *Updater) UpdateDatabase() error {
	update, err := u.ShouldUpdate()
	if err != nil || !update {
		return err
	}
	err = u.UpdateDeputiesAndFactions()
	if err != nil {
		return err
	}
	return nil
}
