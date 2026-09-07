package dumastructs

// Law represents a draft law from the Duma API.
type Law struct {
	Id          int      `json:"id"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Number      string   `json:"number"`
	Tags        []string `json:"tags,omitempty"`
}
