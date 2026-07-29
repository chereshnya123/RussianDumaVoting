package dumastructs

// Структура соответствует открытым данным vote.duma.gov.ru и API АСОЗД
type Deputy struct {
	ID      string      `json:"id"`
	Name    string      `json:"name"`
	Faction FactionInfo `json:"faction"`
}

type FactionInfo struct {
	ID    string `json:"id"`
	Title string `json:"title"` // Например: "Единая Россия", "КПРФ"
}

type VoteRecord struct {
	Deputy Deputy `json:"deputy"`
	Result string `json:"result"` // "accept", "declice" (опечатка в API Госдумы означает decline), "abstain", "none"
}

type DumaVote struct {
	Code     string       `json:"code"`
	Title    string       `json:"title"`
	Datetime string       `json:"datetime"`
	Votes    []VoteRecord `json:"votes"`
}

// Наши внутренние структуры для отображения
type VoteResult struct {
	For       int
	Against   int
	Abstained int
	NotVoted  int
}

type Law struct {
	ID          string
	Title       string
	Description string
	Tags        []string
	Votes       map[string]VoteResult // Ключ - название фракции
}

type Faction struct {
	Name string
}
