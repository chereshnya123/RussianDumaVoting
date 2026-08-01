package server

import (
	"dumaVote/db"
	"dumaVote/dumaclient"
	"html/template"
	"log"
	"net/http"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

type dumaVotesServer struct {
	apiDumaClient dumaclient.DumaVoteClient
	db            *db.Database
}

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
        <h1>📋 Повестка дня</h1>
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

var funcMap = template.FuncMap{
	"joinTags": func(tags []string) string {
		return strings.Join(tags, ", ")
	},
}

func NewDumaVotesServer(apiKey, personalKey string) dumaVotesServer {
	db, err := db.NewDatabase("RussianDumaVotes")
	if err != nil {
		log.Fatalf("Can not create duma database. Error = %v", err)
	}

	return dumaVotesServer{apiDumaClient: dumaclient.NewDumaVoteClient(), db: db}
}

func (d *dumaVotesServer) MainHandler(w http.ResponseWriter, r *http.Request) {
	question, err := d.apiDumaClient.GetLastMeetingQuestion()
	if err != nil {
		log.Printf("API error: %v.", err)
		// Reply user with error
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl, err := template.New("webpage").Funcs(funcMap).Parse(htmlTemplate)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, []dumaclient.Question{question})
	if err != nil {
		log.Print(err.Error())
	}
}
