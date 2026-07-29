package parse

import (
	"dumaVote/internal/dumastructs"
	"dumaVote/internal/tag"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"sort"
	"strings"
	"time"
)

type TagStat struct {
	Name  string
	Count int
}

func FetchVoteInfo(appKey, personalKey string, voteId int) ([]dumastructs.Law, []dumastructs.Faction, error) {
	// ЗАМЕНИТЕ этот URL на точный endpoint из вашей документации по API-ключу
	// Например: "https://api.duma.gov.ru/api/v1/vote/{vote_ID}.{json/xml}?app_token={token}&limit=50" или аналогичный для законопроектов
	apiURL := fmt.Sprintf("http://api.duma.gov.ru/api/%s/vote/80988.json?app_token=%s", personalKey, appKey)
	log.Printf("Calling: %s\n", apiURL)
	req, err := http.NewRequest("GET", apiURL, nil)
	log.Printf("Do GET request")
	if err != nil {
		return nil, nil, err
	}

	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Set("Accept-Encoding", "gzip, deflate")
	req.Header.Set("Accept-Language", "ru-RU,ru;q=0.9,en;q=0.8")
	req.Header.Set("Cache-Control", "max-age=0")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Cookie", "symfony=ssqkfsohai78j66e03eag3a3g6")
	req.Header.Set("Dnt", "1")
	req.Header.Set("Host", "api.duma.gov.ru")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/146.0.0.0 YaBrowser/26.4.0.0 Safari/537.36")

	client := &http.Client{Timeout: 30 * time.Minute}
	resp, err := client.Do(req)
	if err != nil && !errors.Is(err, io.EOF) {
		return nil, nil, fmt.Errorf("ошибка запроса к API: %w", err)
	}
	log.Printf("API request succeeded")
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, nil, fmt.Errorf("API вернул статус %d: %s", resp.StatusCode, string(body))
	}
	log.Printf("API request is OK")

	defer func() { _ = resp.Body.Close() }()
	var rawVotes []dumastructs.DumaVote
	if err := json.NewDecoder(resp.Body).Decode(&rawVotes); err != nil {
		return nil, nil, fmt.Errorf("ошибка парсинга JSON: %w", err)
	}
	log.Printf("Parse response")

	// Собираем уникальные фракции и преобразуем в наши модели
	factionMap := make(map[string]bool)
	var laws []dumastructs.Law

	for _, rv := range rawVotes {
		factionVotes := aggregateVotesByFaction(rv.Votes)

		for fName := range factionVotes {
			factionMap[fName] = true
		}

		law := dumastructs.Law{
			ID:          rv.ID,
			Title:       rv.Title,
			Description: "Голосование по вопросу", // В расширенной версии можно подтянуть из /api/bill/{id}
			Votes:       factionVotes,
		}
		tag.AssignTags(&law)
		laws = append(laws, law)
	}

	var factions []dumastructs.Faction
	for fName := range factionMap {
		factions = append(factions, dumastructs.Faction{Name: fName})
	}

	// Сортируем фракции для стабильного отображения (опционально)
	sort.Slice(factions, func(i, j int) bool {
		return factions[i].Name < factions[j].Name
	})

	return laws, factions, nil
}

func aggregateVotesByFaction(votes []dumastructs.VoteRecord) map[string]dumastructs.VoteResult {
	factionVotes := make(map[string]dumastructs.VoteResult)

	for _, v := range votes {
		factionName := v.Deputy.Faction.Title
		if factionName == "" {
			factionName = "Беспартийные"
		}

		res := factionVotes[factionName]
		switch strings.ToLower(v.Result) {
		case "accept":
			res.For++
		case "declice", "decline": // Учитываем известную опечатку в API Госдумы ("declice")
			res.Against++
		case "abstain":
			res.Abstained++
		case "none":
			res.NotVoted++
		}
		factionVotes[factionName] = res
	}

	return factionVotes
}

// ==================== АНАЛИТИКА ====================

func CalculateFactionDirections(laws []dumastructs.Law, factions []dumastructs.Faction) map[string][]TagStat {
	result := make(map[string][]TagStat)

	for _, f := range factions {
		tagCounts := make(map[string]int)
		for _, law := range laws {
			vote, exists := law.Votes[f.Name]
			if !exists {
				continue
			}
			// Если фракция в целом поддержала закон (За > Против)
			if vote.For > vote.Against {
				for _, tag := range law.Tags {
					tagCounts[tag]++
				}
			}
		}

		var stats []TagStat
		for tag, count := range tagCounts {
			stats = append(stats, TagStat{Name: tag, Count: count})
		}
		// Сортируем по убыванию частоты
		sort.Slice(stats, func(i, j int) bool {
			return stats[i].Count > stats[j].Count
		})
		result[f.Name] = stats
	}
	return result
}
