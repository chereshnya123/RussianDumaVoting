package analyzer

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"slices"

	"dumaVote/internal/utils"
)

// Fetcher makes HTTP requests to the Duma API and unmarshals the response.
type Fetcher struct {
	appApiKey      string
	personApiToken string
	logger         *slog.Logger
}

func (f Fetcher) getAllDeputiesApiUrl() string {
	const deputiesURL = "http://api.duma.gov.ru/api/%s/deputies.json?app_token=%s"
	return fmt.Sprintf(deputiesURL, f.personApiToken, f.appApiKey)
}

func (f Fetcher) getDeputyInfoApiUrl(deputyId string) string {
	const deputyInfoURL = "http://api.duma.gov.ru/api/%s/deputy.json?app_token=%s&id=%s"
	return fmt.Sprintf(deputyInfoURL, f.personApiToken, f.appApiKey, deputyId)
}

func (f Fetcher) getVotingsApiUrl(pageNum, limit int) string {
	const votingsApiUrl = "http://api.duma.gov.ru/api/%s/voteSearch.json?app_token=%s&page=%d&limit=%d"
	return fmt.Sprintf(votingsApiUrl, f.personApiToken, f.appApiKey, pageNum, limit)
}

// NewFetcher creates a new Fetcher.
func NewFetcher(appApiKey, personApiKey string, logger *slog.Logger) *Fetcher {
	return &Fetcher{appApiKey: appApiKey, personApiToken: personApiKey, logger: logger}
}

// FetchAllDeputies fetches all deputies from the Duma API.
// The API returns a bare JSON array [{...}], not a wrapper object.
func (f *Fetcher) FetchAllDeputies() ([]Deputy, error) {

	apiURL := f.getAllDeputiesApiUrl()
	f.logger.Debug("Fetch all deputy data", "url", apiURL)
	resp, err := utils.DoSimpleRequest(apiURL)
	if err != nil {
		return nil, fmt.Errorf("fetch deputies request failed: %w", err)
	}
	if resp == nil {
		return nil, nil
	}
	defer func() { _ = resp.Body.Close() }()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read deputies response body: %w", err)
	}

	var deputies []Deputy
	if err := json.Unmarshal(bodyBytes, &deputies); err != nil {
		return nil, fmt.Errorf("unmarshal deputies response: %w", err)
	}

	return deputies, nil
}

// FetchDeputyInfo fetches detailed deputy profile from the Duma API.
// The API returns a single JSON object { ... } with the deputy's data.
func (f *Fetcher) FetchDeputyInfo(deputyId string) (DeputyInfo, error) {
	var empty DeputyInfo

	apiURL := f.getDeputyInfoApiUrl(deputyId)
	resp, err := utils.DoSimpleRequest(apiURL)
	if err != nil {
		return empty, fmt.Errorf("fetch deputy info request failed: %w", err)
	}

	if resp == nil {
		return empty, nil
	}
	defer func() { _ = resp.Body.Close() }()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return empty, fmt.Errorf("read deputy info response body: %w", err)
	}

	var deputyResp DeputyInfo
	if err := json.Unmarshal(bodyBytes, &deputyResp); err != nil {
		// Fallback: try unmarshalling as array [{...}]
		var deputies []DeputyInfo
		if arrErr := json.Unmarshal(bodyBytes, &deputies); arrErr == nil && len(deputies) > 0 {
			return deputies[0], nil
		}
		return empty, fmt.Errorf("unmarshal deputy info response (tried object and array): %w", err)
	}

	return deputyResp, nil
}

// Requests `limit` votings from page with given number
func (f *Fetcher) FetchVotings(pageNum, limit int) (VoteResponse, error) {
	if !slices.Contains([]int{5, 10, 20, 50, 100}, limit) {
		return VoteResponse{}, fmt.Errorf("Can not fetch votings. Get unexpected `limit` parameter. Available values = [5, 10, 20, 50, 100]")
	}
	votingsApiUrl := f.getVotingsApiUrl(pageNum, limit)
	resp, err := utils.DoSimpleRequest(votingsApiUrl)

	if err != nil {
		return VoteResponse{}, fmt.Errorf("Can not fetch votes request: %w", err)
	}

	if resp == nil {
		return VoteResponse{}, nil
	}
	defer func() { _ = resp.Body.Close() }()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return VoteResponse{}, fmt.Errorf("Can not read votes response body: %w", err)
	}

	var votesResp VoteResponse
	if err := json.Unmarshal(bodyBytes, &votesResp); err != nil {
		return VoteResponse{}, fmt.Errorf("Can not unmarshal votes response: %w", err)
	}

	return votesResp, nil
}
