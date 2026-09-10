package analyzer

import (
	"cmp"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"slices"
	"strconv"

	"dumaVote/internal/utils"
)

// Fetcher makes HTTP requests to the Duma API and unmarshals the response.
type Fetcher struct {
	appApiKey      string
	personApiToken string
	logger         *slog.Logger
}

func (f Fetcher) getAllDeputiesApiUrl() string {
	const deputiesUrlTemplate = "http://api.duma.gov.ru/api/%s/deputies.json?app_token=%s"
	return fmt.Sprintf(deputiesUrlTemplate, f.personApiToken, f.appApiKey)
}

func (f Fetcher) getDeputyInfoApiUrl(deputyId string) string {
	const deputyInfoUrlTemplate = "http://api.duma.gov.ru/api/%s/deputy.json?app_token=%s&id=%s"
	return fmt.Sprintf(deputyInfoUrlTemplate, f.personApiToken, f.appApiKey, deputyId)
}

func (f Fetcher) getVotingsApiUrl(params map[string]string) string {
	const votingsApiUrlTemplate = "http://api.duma.gov.ru/api/%s/voteSearch.json?app_token=%s&page=%s&limit=%s"
	votingsApiUrl := fmt.Sprintf(votingsApiUrlTemplate, f.personApiToken, f.appApiKey, "0", "0")
	paramsInserted := 0
	for key, value := range params {
		if paramsInserted == 0 {
			votingsApiUrl = fmt.Sprintf("%s?%s=%s", votingsApiUrl, key, value)
			paramsInserted++
			continue
		}
		votingsApiUrl = fmt.Sprintf("%s&%s=%s", votingsApiUrl, key, value)
	}

	return votingsApiUrl
}

// NewFetcher creates a new Fetcher.
func NewFetcher(appApiKey, personApiKey string, logger *slog.Logger) *Fetcher {
	return &Fetcher{appApiKey: appApiKey, personApiToken: personApiKey, logger: logger}
}

// Generic function for fetching responses from duma api
func (f *Fetcher) FetchDataFromApi(parsedOutput any, apiUrl string) error {
	resp, err := utils.DoSimpleRequest(apiUrl)
	if err != nil {
		return fmt.Errorf("fetch deputies request failed: %w", err)
	}
	if resp == nil {
		return nil
	}
	defer func() { _ = resp.Body.Close() }()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read deputies response body: %w", err)
	}

	if err := json.Unmarshal(bodyBytes, parsedOutput); err != nil {
		return fmt.Errorf("unmarshal deputies response: %w", err)
	}

	return nil
}

// FetchAllDeputies fetches all deputies from the Duma API.
// The API returns a bare JSON array [{...}], not a wrapper object.
func (f *Fetcher) FetchAllDeputies() ([]Deputy, error) {

	apiURL := f.getAllDeputiesApiUrl()
	f.logger.Debug("Fetch all deputy data", "url", apiURL)
	var deputies []Deputy
	err := f.FetchDataFromApi(&deputies, apiURL)
	return deputies, err
}

// FetchDeputyInfo fetches detailed deputy profile from the Duma API.
// The API returns a single JSON object { ... } with the deputy's data.
func (f *Fetcher) FetchDeputyInfo(deputyId string) (DeputyInfo, error) {
	apiURL := f.getDeputyInfoApiUrl(deputyId)
	var deputyInfo DeputyInfo
	err := f.FetchDataFromApi(&deputyInfo, apiURL)
	return deputyInfo, err
}

// Requests `limit` votings from page with given number
func (f *Fetcher) FetchVotings(pageNum, limit int) (VoteResponse, error) {
	if !slices.Contains([]int{5, 10, 20, 50, 100}, limit) {
		return VoteResponse{}, fmt.Errorf("can not fetch votings. Get unexpected `limit` parameter. Available values = [5, 10, 20, 50, 100]")
	}
	params := map[string]string{"page_num": strconv.Itoa(pageNum), "limit": strconv.Itoa(limit)}
	votingsApiUrl := f.getVotingsApiUrl(params)

	var voteResponse VoteResponse
	err := f.FetchDataFromApi(&voteResponse, votingsApiUrl)

	return voteResponse, err
}

func (f *Fetcher) GetLastLawVoting(lawNumber string) (Vote, error) {
	allVotingsInfos, err := f.fetchAllLawVotings(lawNumber)
	if err != nil {
		return Vote{}, fmt.Errorf("can not get last law voting: %w", err)
	}
	return slices.MaxFunc(allVotingsInfos, func(a, b Vote) int { return cmp.Compare(a.VoteDate, b.VoteDate) }), nil
}

// Must be generic: I) return struct II) func to generate url
func (f *Fetcher) fetchAllLawVotings(lawNumber string) ([]Vote, error) {
	params := map[string]string{"number": lawNumber}
	votingsApiUrl := f.getVotingsApiUrl(params)

	var votesResp VoteResponse
	err := f.FetchDataFromApi(&votesResp, votingsApiUrl)
	return votesResp.Votes, err
}

func (f *Fetcher) fetchVotingInfo(votingId int) (VoteDetailResponse, error) {
	const votingInfoUrlTemplate = "http://api.duma.gov.ru/api/%s/vote/%d.json?app_token=%s"
	apiURL := fmt.Sprintf(votingInfoUrlTemplate, f.personApiToken, votingId, f.appApiKey)

	var votingInfo VoteDetailResponse
	err := f.FetchDataFromApi(&votingInfo, apiURL)

	return votingInfo, err
}
