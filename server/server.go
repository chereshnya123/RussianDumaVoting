// Package server provides the HTTP server for the dumaVote application.
package server

import (
	"dumaVote/db"
	"dumaVote/dumaclient"
	"html/template"
	"log/slog"
	"net/http"
	"strings"
	"time"
)

const htmlTemplate = `
<!DOCTYPE html>
<html lang="ru">
<head>
    <meta charset="UTF-8">
    <title>Повестка заседания</title>
    <style>
        body { font-family: sans-serif; background: #f4f6f8; padding: 20px; }
        .container { max-width: 800px; margin: 0 auto; }
        .card { background: white; padding: 15px; margin-bottom: 15px; border-radius: 8px; box-shadow: 0 2px 4px rgba(0,0,0,0.1); border-left: 4px solid #3498db; }
        .header { display: flex; justify-content: space-between; color: #7f8c8d; font-size: 0.9em; margin-bottom: 8px; }
        .title { font-weight: bold; font-size: 1.1em; margin-bottom: 10px; }
        .meta { font-size: 0.85em; color: #95a5a6; }
    </style>
</head>
<body>
    <div class="container">
        <h1>Повестка дня</h1>
        {{range .}}
        <div class="card">
            <div class="header">
                <span>Заседание #{{.Kodz}}</span>
                <span>Вопрос №{{.Kodvopr}}</span>
            </div>
            <div class="title">{{.Name}}</div>
            <div class="meta">
                📅 {{.Datez}} | 📝 Стенограмма: стр. {{.Nbegin}}–{{.Nend}}
            </div>
        </div>
        {{end}}
    </div>
</body>
</html>`

var htmlTmpl = template.Must(template.New("webpage").Funcs(template.FuncMap{
	"joinTags": func(tags []string) string {
		return strings.Join(tags, ", ")
	},
}).Parse(htmlTemplate))

// DumaVotesServer is the HTTP server handler.
type DumaVotesServer struct {
	apiDumaClient dumaclient.DumaVoteClient
	db            *db.Database
	logger        *slog.Logger
}

// NewDumaVotesServer creates a new server instance with database and API client.
func NewDumaVotesServer(apiKey, personalKey string, logger *slog.Logger) *DumaVotesServer {
	db, err := db.NewDatabase("RussianDumaVote")
	if err != nil {
		logger.Error("Failed to create database", "error", err)
		panic(err)
	}

	return &DumaVotesServer{
		apiDumaClient: dumaclient.NewDumaVoteClient(logger),
		db:            db,
		logger:        logger,
	}
}

// MainHandler is the HTTP handler for the main page.
func (s *DumaVotesServer) MainHandler(w http.ResponseWriter, r *http.Request) {
	question, err := s.apiDumaClient.GetLastMeetingQuestion()
	if err != nil {
		s.logger.Error("API error", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := htmlTmpl.Execute(w, []dumaclient.Question{question}); err != nil {
		s.logger.Error("Template execution failed", "error", err)
	}
}

// ShouldUpdate checks if a data refresh is needed based on the last update time.
func (s *DumaVotesServer) ShouldUpdate() (bool, error) {
	lastUpdate, err := s.db.GetLastUpdateTime()
	if err != nil {
		return false, err
	}

	if lastUpdate.IsZero() {
		return true, nil
	}

	return time.Since(lastUpdate) > time.Hour, nil
}

// UpdateData refreshes the data if enough time has elapsed.
func (s *DumaVotesServer) UpdateData() error {
	update, err := s.ShouldUpdate()
	if err != nil || !update {
		return err
	}

	return nil
}

// Close closes the underlying database connection.
func (s *DumaVotesServer) Close() error {
	return s.db.Close()
}
