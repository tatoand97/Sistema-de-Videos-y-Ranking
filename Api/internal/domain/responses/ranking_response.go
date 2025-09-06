package responses

// RankingEntry is the public API schema for /api/public/rankings
// Note: position is assigned per returned page starting at 1.
type RankingEntry struct {
	Position int     `json:"position"`
	Username string  `json:"username"`
	City     *string `json:"city,omitempty"`
	Votes    int     `json:"votes"`
}

// RankingItem is an internal representation without position
// used by repositories/services before pagination position numbering.
type RankingItem struct {
	Username string  `json:"username"`
	City     *string `json:"city,omitempty"`
	Votes    int     `json:"votes"`
}
