package tracemoe

type SearchResult struct {
	AnilistID  int     `json:"anilist"`
	Filename   string  `json:"filename"`
	Episode    *int    `json:"episode,omitempty"`
	From       float64 `json:"from"`
	To         float64 `json:"to"`
	Similarity float64 `json:"similarity"`
	Video      string  `json:"video"`
	Image      string  `json:"image"`
}

type SearchResponse struct {
	FrameCount int            `json:"frameCount"`
	Error      string         `json:"error,omitempty"`
	Result     []SearchResult `json:"result"`
}
