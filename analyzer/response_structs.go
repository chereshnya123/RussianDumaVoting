package analyzer

// Question represents a single meeting issue from the Duma API.
type Question struct {
	Name    string `json:"name"`    // Question name
	Datez   string `json:"datez"`   // Meeting date
	Kodz    int    `json:"kodz"`    // Meeting code
	Kodvopr int    `json:"kodvopr"` // Issue sequence number
	Nbegin  int    `json:"nbegin"`  // First page of transcript
	Nend    int    `json:"nend"`    // Last page of transcript
}

// VoteResponse represents the top-level Duma API response for voting results.
type VoteResponse struct {
	TotalCount int      `json:"totalCount"` // Total number of voting records
	Page       int      `json:"page"`       // Current page number
	PageSize   int      `json:"pageSize"`   // Number of items per page
	Wordings   []string `json:"wording"`    // Description of the voting period/query
	Votes      []Vote   `json:"votes"`      // Array of voting records
}

// Vote represents a single voting record from the Duma API.
type Vote struct {
	Id           int    `json:"id"`           // Voting identifier
	Subject      string `json:"subject"`      // Issue description with number
	VoteDate     string `json:"voteDate"`     // Voting timestamp
	PersonResult string `json:"personResult"` // Result for searched deputy (for/against/abstain/absent)
	VoteCount    int    `json:"voteCount"`    // Total votes cast (or faction size if searched by faction)
	ForCount     int    `json:"forCount"`     // Votes "for" (fraction-internal if searched by faction)
	AgainstCount int    `json:"againstCount"` // Votes "against" (fraction-internal if searched by faction)
	AbstainCount int    `json:"abstainCount"` // Votes "abstain" (fraction-internal if searched by faction)
	AbsentCount  int    `json:"absentCount"`  // Non-voting deputies (fraction-internal if searched by faction)
	ResultType   string `json:"resultType"`   // Type: количественное, рейтинговое, качественное, альтернативное
	Result       bool   `json:"result"`       // true — adopted, false — not adopted
}

// Deputy represents a brief deputy profile from search/list responses.
type Deputy struct {
	Id        string `json:"id"`        // Deputy/Federation Council member identifier
	Name      string `json:"name"`      // Full name
	Position  string `json:"position"`  // "Депутат ГД" or "Член СФ"
	IsCurrent bool   `json:"isCurrent"` // true if currently holding the position
}

// DeputiesResponse represents the top-level Duma API response for parliament member searches.
type DeputiesResponse struct {
	TotalCount int      `json:"totalCount"` // Total number of persons found
	Page       int      `json:"page"`       // Current page number
	PageSize   int      `json:"pageSize"`   // Number of items per page
	Wordings   []string `json:"wording"`    // Description of the search query
	People     []Deputy `json:"people"`     // Array of parliament members
}

// DeputyEducation represents one record of a deputy's higher education.
type DeputyEducation struct {
	Institution string   `json:"institution"` // University name
	Year        string   `json:"year"`        // Year of graduation
	Degrees     []string `json:"degrees"`     // Academic degrees
	Ranks       []string `json:"ranks"`       // Academic titles
}

// DeputyActivity represents one record of a deputy's activity (committee, group, etc.).
type DeputyActivity struct {
	Name                    string `json:"name"`                    // Activity name (e.g. "Член комитета")
	SubdivisionId           int64  `json:"subdivisionId"`           // Subdivision identifier
	SubdivisionNameGenitive string `json:"subdivisionNameGenitive"` // Subdivision name in genitive case
}

// DeputyInfo represents a full deputy profile from the Duma API.
type DeputyInfo struct {
	Id                string            `json:"id"`                // Deputy identifier
	Family            string            `json:"family"`            // Surname
	Name              string            `json:"name"`              // First name
	Patronymic        string            `json:"patronymic"`        // Patronymic
	Birthdate         string            `json:"birthdate"`         // Date of birth (YYYY-MM-DD)
	CredentialsStart  string            `json:"credentialsStart"`  // Start of mandate (YYYY-MM-DD)
	CredentialsEnd    *string           `json:"credentialsEnd"`    // End of mandate (nullable)
	FactionId         string            `json:"factionId"`         // Faction identifier
	FactionName       string            `json:"factionName"`       // Full faction name
	FactionRole       string            `json:"factionRole"`       // Role in faction (sorting prefix at start)
	PartyNameInstr    string            `json:"partyNameInstr"`    // Party name in instrumental case
	IsActual          string            `json:"isActual"`          // "true"/"false" — mandate active now
	HomePage          *string           `json:"homePage"`          // Personal website (nullable)
	FactionRegion     string            `json:"factionRegion"`     // Faction-region link
	NameGenitive      string            `json:"nameGenitive"`      // Full name in genitive case
	LawCount          int               `json:"lawcount"`          // Number of draft laws initiated
	Regions           []string          `json:"regions"`           // Regions linked to deputy (array of strings)
	FamilyAndInitials string            `json:"familyAndInitials"` // Surname + initials
	SpeachCount       int               `json:"speachCount"`       // Number of speeches
	VoteLink          string            `json:"voteLink"`          // Link to deputy's votes
	TranscriptLink    string            `json:"transcriptLink"`    // Link to deputy's transcripts
	Educations        []DeputyEducation `json:"educations"`        // Education records
	Activities        []DeputyActivity  `json:"activity"`          // Activity records
}
