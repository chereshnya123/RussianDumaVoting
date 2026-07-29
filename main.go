package main

import (
	"dumaVote/internal/dumastructs"
	"dumaVote/internal/parse"
	"dumaVote/internal/tag"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"
)

// ==================== ВЕБ-СЕРВЕР ====================

const htmlTemplate = `
<!DOCTYPE html>
<html lang="ru">
<head>
    <meta charset="UTF-8">
    <title>Аналитика Госдумы РФ</title>
    <style>
        body { font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Arial, sans-serif; line-height: 1.6; margin: 40px; background: #f8f9fa; color: #333; }
        h1, h2 { color: #2c3e50; text-align: center; }
        .container { max-width: 1000px; margin: 0 auto; }
        .faction-card { background: #fff; padding: 25px; margin-bottom: 20px; border-radius: 8px; box-shadow: 0 2px 8px rgba(0,0,0,0.08); }
        .faction-name { font-size: 1.4em; font-weight: bold; color: #2c3e50; border-bottom: 3px solid #3498db; padding-bottom: 10px; margin-bottom: 15px; }
        .tags { display: flex; flex-wrap: wrap; gap: 10px; margin-top: 10px; }
        .tag { background: #e3f2fd; color: #1565c0; padding: 6px 14px; border-radius: 20px; font-weight: 600; font-size: 0.9em; }
        .tag-count { font-size: 0.85em; color: #555; margin-left: 4px; }
        .law-card { background: #fff; padding: 20px; margin-bottom: 15px; border-radius: 8px; border-left: 5px solid #3498db; box-shadow: 0 1px 4px rgba(0,0,0,0.05); }
        .law-title { font-size: 1.1em; font-weight: bold; margin-bottom: 8px; }
        .law-tags { color: #666; font-size: 0.9em; margin-bottom: 12px; font-style: italic; }
        .vote-table { width: 100%; border-collapse: collapse; font-size: 0.9em; margin-top: 10px; }
        .vote-table th, .vote-table td { padding: 8px; text-align: left; border-bottom: 1px solid #eee; }
        .vote-table th { background: #f1f3f5; color: #495057; }
        .alert { background: #fff3cd; color: #856404; padding: 15px; border-radius: 5px; text-align: center; margin-bottom: 20px; border: 1px solid #ffeaa7; }
    </style>
</head>
<body>
    <div class="container">
        <h1>Политические приоритеты фракций Госдумы</h1>
        {{if .ErrorMessage}}
            <div class="alert">
                ⚠️ <strong>Внимание:</strong> Не удалось загрузить данные из API. 
                <br>Проверьте API-ключ и URL. Ошибка: {{.ErrorMessage}}
                <br><em>Показаны демонстрационные данные.</em>
            </div>
        {{else}}
            <p style="text-align:center; color:#666;">Данные получены в реальном времени через API Госдумы РФ</p>
        {{end}}

        <h2>Основные направления законотворчества</h2>
        {{range .Factions}}
        <div class="faction-card">
            <div class="faction-name">{{.Name}}</div>
            <strong>Поддерживаемые темы:</strong>
            <div class="tags">
                {{range (index $.Directions .Name)}}
                    <span class="tag">{{.Name}} <span class="tag-count">({{.Count}})</span></span>
                {{end}}
                {{if not (index $.Directions .Name)}}
                    <span style="color:#999; font-size:0.9em;">Нет данных о поддержанных законах</span>
                {{end}}
            </div>
        </div>
        {{end}}

        <h2 style="margin-top: 40px;">Последние голосования</h2>
        {{range .Laws}}
        <div class="law-card">
            <div class="law-title">{{.Title}}</div>
            <div class="law-tags">Теги: {{joinTags .Tags}}</div>
            <table class="vote-table">
                <tr>
                    <th>Фракция</th>
                    <th>✅ За</th>
                    <th>❌ Против</th>
                    <th>⬜ Воздерж.</th>
                    <th>➖ Не голос.</th>
                </tr>
                {{range $factionName, $vote := .Votes}}
                <tr>
                    <td><strong>{{$factionName}}</strong></td>
                    <td style="color:#27ae60; font-weight:bold;">{{$vote.For}}</td>
                    <td style="color:#c0392b; font-weight:bold;">{{$vote.Against}}</td>
                    <td>{{$vote.Abstained}}</td>
                    <td>{{$vote.NotVoted}}</td>
                </tr>
                {{end}}
            </table>
        </div>
        {{end}}
    </div>
</body>
</html>
`

var funcMap = template.FuncMap{
	"joinTags": func(tags []string) string {
		return strings.Join(tags, ", ")
	},
}

type PageData struct {
	Factions     []dumastructs.Faction
	Laws         []dumastructs.Law
	Directions   map[string][]parse.TagStat
	ErrorMessage string
}

func isExist(env_name string) bool {
	_, isExists := os.LookupEnv("APP_API_KEY")
	return isExists
}

func writeApiError(w http.ResponseWriter) {
	w.WriteHeader(500)
	_, _ = w.Write([]byte("Got an error while processing request to api.duma.gov.ru"))
}

func mainHandler(w http.ResponseWriter, r *http.Request) {
	if !isExist("APP_API_KEY") || !isExist("PERSONAL_API_KEY") {
		log.Fatal("APP_API_KEY or PERSONAL_API_KEY are not set")
	}
	appKey := os.Getenv("APP_API_KEY")
	pesonalKey := os.Getenv("PERSONAL_API_KEY")
	var laws []dumastructs.Law
	var factions []dumastructs.Faction
	var errMsg string

	voteId := 80988
	realLaws, realFactions, err := parse.FetchVoteInfo(appKey, pesonalKey, voteId)
	if err != nil {
		log.Printf("API error: %v.", err)
		errMsg = err.Error()
		// Reply user with error
		writeApiError(w)
		return
	}
	laws = realLaws
	factions = realFactions
	directions := parse.CalculateFactionDirections(laws, factions)

	data := PageData{
		Factions:     factions,
		Laws:         laws,
		Directions:   directions,
		ErrorMessage: errMsg,
	}

	tmpl, err := template.New("webpage").Funcs(funcMap).Parse(htmlTemplate)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		log.Print(err.Error())
	}
}

// getMockData оставлен как надежный фоллбэк, если API-ключ не подходит или эндпоинт изменился
func getMockData() ([]dumastructs.Law, []dumastructs.Faction) {
	// ... (тот же код моковых данных из Итерации 1, сокращен для краткости, но он должен быть здесь) ...
	factions := []dumastructs.Faction{{Name: "Единая Россия"}, {Name: "КПРФ"}, {Name: "ЛДПР"}, {Name: "Справедливая Россия"}}
	laws := []dumastructs.Law{
		{ID: "1", Title: "О повышении пенсий для инвалидов", Description: "Дополнительные льготы",
			Votes: map[string]dumastructs.VoteResult{"Единая Россия": {For: 300, Against: 10}, "КПРФ": {For: 40, Against: 0}}},
	}
	for i := range laws {
		tag.AssignTags(&laws[i])
	}
	return laws, factions
}

func main() {
	http.HandleFunc("/", mainHandler)
	port := ":8080"
	log.Printf("Started server on http://localhost%s", port)
	if !isExist("APP_API_KEY") {
		log.Fatal("App api key is not set. APP_API_KEY env is empty")
	}
	if !isExist("PERSONAL_API_KEY") {
		log.Fatal("App api key is not set. PERSONAL_API_KEY env is empty")
	}
	log.Fatal(http.ListenAndServe(port, nil))
}
