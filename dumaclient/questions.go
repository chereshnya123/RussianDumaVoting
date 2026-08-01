// Package dumaclient provides a client for the Gosduma API.
package dumaclient

import (
	"dumaVote/internal/utils"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"os"
	"sort"
	"time"
)

// Question represents a single meeting issue from the Gosduma API.
type Question struct {
	Name    string `json:"name"`    // Question name
	Datez   string `json:"datez"`   // Meeting date
	Kodz    int    `json:"kodz"`    // Meeting code
	Kodvopr int    `json:"kodvopr"` // Issue sequence number
	Nbegin  int    `json:"nbegin"`  // First page of transcript
	Nend    int    `json:"nend"`    // Last page of transcript
}

// QuestionResponse is the response structure for the questions API endpoint.
type QuestionResponse struct {
	PageSize   int        `json:"pageSize"`
	Page       int        `json:"page"`
	TotalCount int        `json:"totalCount"`
	Questions  []Question `json:"questions"`
}

// DumaVoteClient interacts with the Gosduma API.
type DumaVoteClient struct {
	appApiKey      string
	personalApiKey string
	logger         *slog.Logger
}

// NewDumaVoteClient creates a new client using environment variables for authentication.
func NewDumaVoteClient(logger *slog.Logger) DumaVoteClient {
	return DumaVoteClient{
		appApiKey:      os.Getenv("APP_API_KEY"),
		personalApiKey: os.Getenv("PERSONAL_API_KEY"),
		logger:         logger,
	}
}

// GetLastMeetingQuestion fetches the most recent meeting question from the API.
func (c DumaVoteClient) GetLastMeetingQuestion() (Question, error) {
	now := time.Now()
	apiURL := fmt.Sprintf(
		"http://api.duma.gov.ru/api/%s/questions.json?app_token=%s&limit=5&dateTo=%s",
		c.personalApiKey, c.appApiKey, now.Format("2006-01-02"),
	)

	c.logger.Info("Calling Gosduma API", "host", "api.duma.gov.ru")

	resp, err := utils.DoSimpleRequest(apiURL)
	if err != nil {
		return Question{}, fmt.Errorf("get request failed: %w", err)
	}

	if resp == nil {
		c.logger.Warn("Received nil response")
		return Question{}, nil
	}
	defer func() { _ = resp.Body.Close() }()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		c.logger.Error("Failed to read response body", "error", err)
		return Question{}, fmt.Errorf("read response body: %w", err)
	}

	var response QuestionResponse
	if err := json.Unmarshal(bodyBytes, &response); err != nil {
		c.logger.Error("Failed to parse JSON response", "error", err)
		return Question{}, fmt.Errorf("parse JSON: %w", err)
	}

	sort.Slice(response.Questions, func(i, j int) bool {
		return response.Questions[i].Datez > response.Questions[j].Datez
	})

	if len(response.Questions) == 0 {
		c.logger.Warn("Empty questions list")
		return Question{}, nil
	}

	return response.Questions[0], nil
}
