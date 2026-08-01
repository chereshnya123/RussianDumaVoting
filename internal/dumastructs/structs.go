package dumastructs

type Deputy struct {
	ID      string      `json:"id"`
	Name    string      `json:"name"`
	Faction FactionInfo `json:"faction"`
}

type FactionInfo struct {
	ID    string `json:"id"`
	Title string `json:"title"` // For example: "United Russia", "KPRF"
}

type VoteRecord struct {
	Deputy Deputy `json:"deputy"`
	Result string `json:"result"` // "accept", "declice" (Typo in Gosduma API: decline), "abstain", "none"
}

type DumaVote struct {
	Code     string       `json:"code"`
	Title    string       `json:"title"`
	Datetime string       `json:"datetime"`
	Votes    []VoteRecord `json:"votes"`
}

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
	Votes       map[string]VoteResult
}

type Faction struct {
	Name string
}
