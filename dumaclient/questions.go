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

// Separate meeting issue
type Question struct {
	Name    string `json:"name"`    // Question name
	Datez   string `json:"datez"`   // Meeting date
	Kodz    int    `json:"kodz"`    // Meeting code
	Kodvopr int    `json:"kodvopr"` // Sequence number of issue
	Nbegin  int    `json:"nbegin"`  // First issue num
	Nend    int    `json:"nend"`    // Second issue num
}

type QuestionResponse struct {
	PageSize   int        `json:"pageSize"`
	Page       int        `json:"page"`
	TotalCount int        `json:"totalCount"`
	Questions  []Question `json:"questions"`
}

type DumaVoteClient struct {
	appApiKey      string
	personalApiKey string
	logger         *slog.Logger
}

func NewDumaVoteClient(logger *slog.Logger) DumaVoteClient {
	appApiKey := os.Getenv("APP_API_KEY")
	personalApiKey := os.Getenv("PERSONAL_API_KEY")
	return DumaVoteClient{appApiKey: appApiKey, personalApiKey: personalApiKey, logger: logger}
}

func (d DumaVoteClient) GetLastMeetingQuestion() (Question, error) {
	now := time.Now()
	apiURL := fmt.Sprintf("http://api.duma.gov.ru/api/%s/questions.json?app_token=%s&limit=5&dateTo=%s", d.personalApiKey, d.appApiKey, now.Format("2006-01-02"))
	d.logger.Info(fmt.Sprintf("Calling: %s", apiURL))
	resp, err := utils.DoSimpleRequest(apiURL)
	if err != nil {
		d.logger.Error(fmt.Sprintf("Get an error while request. Error = %v", err.Error()))
		return Question{}, fmt.Errorf("get an error while request. Error = %v", err.Error())
	}

	if resp == nil {
		d.logger.Warn("Get null response")
		return Question{}, nil
	}
	bodyBytes, err := io.ReadAll(resp.Body)

	if err != nil {
		d.logger.Error(fmt.Sprintf("Can not read body from /questions response. Err = %v", err))
		return Question{}, fmt.Errorf("can not read body from /questions response. Err = %v", err)
	}
	defer func() { _ = resp.Body.Close() }()

	var response QuestionResponse
	err = json.Unmarshal([]byte(bodyBytes), &response)
	if err != nil {
		d.logger.Error(fmt.Sprintf("Can not parse JSON from /questions.json. Err = %v", err))
		return Question{}, err
	}

	sort.Slice(response.Questions, func(i, j int) (less bool) {
		return response.Questions[i].Datez > response.Questions[j].Datez
	})

	if len(response.Questions) == 0 {
		d.logger.Warn("Empty response")
		return Question{}, nil
	}

	return response.Questions[0], nil
}
