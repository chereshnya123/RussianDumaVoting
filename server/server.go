// Package server provides the HTTP server for the dumaVote application.
package server

import (
	"dumaVote/analyzer"
	"dumaVote/db"
	"dumaVote/dumaclient"
	"fmt"
	"html/template"
	"log"
	"log/slog"
	"net/http"
	"strings"
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

// var htmlTmpl = template.Must(template.New("webpage").Funcs(template.FuncMap{
// 	"joinTags": func(tags []string) string {
// 		return strings.Join(tags, ", ")
// 	},
// }).Parse(htmlTemplate))

// DumaVotesServer is the HTTP server handler.
type DumaVotesServer struct {
	apiDumaClient dumaclient.DumaVoteClient
	db            *db.Database
	logger        *slog.Logger
	analyzer      analyzer.Analyzer
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
		analyzer:      *analyzer.NewAnalyzer(apiKey, personalKey, db, logger),
	}
}

var funcMap = template.FuncMap{
	"joinTags": func(tags []string) string {
		return strings.Join(tags, ", ")
	},
}

// MainHandler is the HTTP handler for the main page.
func (s *DumaVotesServer) MainHandler(w http.ResponseWriter, r *http.Request) {
	s.logger.Debug(fmt.Sprintf("Accept request. URL = %s, URI = %s, method = %s\n", r.URL, r.RequestURI, r.Method))
	s.route(w, r)
}

func (s *DumaVotesServer) route(w http.ResponseWriter, r *http.Request) {
	uri := r.RequestURI
	method := r.Method
	switch {
	case method == http.MethodGet && uri == "/":
		s.HandleLastMeetingQuestion(w, r)
	default:
		s.logger.Error("Get invalid request. URI = %s, method = %s\n", uri, method)
		w.WriteHeader(404)
	}
}

func (s *DumaVotesServer) HandleLastMeetingQuestion(w http.ResponseWriter, r *http.Request) {
	lastQuestion, err := s.analyzer.GetLastQuestion()
	if err != nil {
		s.logger.Error("Can not get last question.", "Error", err)
		w.WriteHeader(500)
		return
	}

	tmpl, err := template.New("webpage").Funcs(funcMap).Parse(htmlTemplate)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, []db.Question{lastQuestion})
	if err != nil {
		log.Print(err.Error())
		s.logger.Error(err.Error())
	}
}

// Close closes the underlying database connection.
func (s *DumaVotesServer) Close() error {
	return s.db.Close()
}
