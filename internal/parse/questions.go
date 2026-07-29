package parse

import (
	"dumaVote/internal/utils"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"sort"
	"time"
)

// Separate meeting issue
type Issue struct {
	Name    string `json:"name"`    // Question name
	Datez   string `json:"datez"`   // Meeting date
	Kodz    int    `json:"kodz"`    // Meeting code
	Kodvopr int    `json:"kodvopr"` // Sequence number of issue
	Nbegin  int    `json:"nbegin"`  // First issue num
	Nend    int    `json:"nend"`    // Second issue num
}

type QuestionResponse struct {
	PageSize   int     `json:"pageSize"`
	Page       int     `json:"page"`
	TotalCount int     `json:"totalCount"`
	Issues     []Issue `json:"questions"`
}

func GetLastMeetingIssue(appKey, personalKey string) (Issue, error) {
	now := time.Now()
	apiURL := fmt.Sprintf("http://api.duma.gov.ru/api/%s/questions.json?api_token=%s&limit=5&dateFrom=%s", personalKey, appKey, now.Format("2006-01-02"))
	log.Printf("Calling: %s\n", apiURL)
	resp, err := utils.DoSimpleRequest(apiURL)
	if err != nil {
		log.Printf("Get an error while request. Error = %v", err.Error())
		return Issue{}, fmt.Errorf("get an error while request. Error = %v", err.Error())
	}
	defer func() { _ = resp.Body.Close() }()

	bodyBytes, err := io.ReadAll(resp.Body)

	if err != nil {
		log.Printf("Can not read body from /questions response. Err = %v", err)
		return Issue{}, fmt.Errorf("can not read body from /questions response. Err = %v", err)
	}

	var response QuestionResponse
	err = json.Unmarshal([]byte(bodyBytes), &response)
	if err != nil {
		log.Printf("Can not parse JSON from /questions.json. Err = %v", err)
		return Issue{}, err
	}

	sort.Slice(response.Issues, func(i, j int) (less bool) {
		return response.Issues[i].Datez > response.Issues[j].Datez
	})

	return response.Issues[0], nil
}
