package anilist

type Anime struct {
	ID          int        `json:"id"`
	Title       AnimeTitle `json:"title"`
	Description string     `json:"description"`
	CoverImage  CoverImage `json:"coverImage"`
	Episodes    *int       `json:"episodes,omitempty"`
	Status      string     `json:"status,omitempty"`
	Format      string     `json:"format,omitempty"`
	Season      string     `json:"season,omitempty"`
	SeasonYear  *int       `json:"seasonYear,omitempty"`
	Genres      []string   `json:"genres,omitempty"`
	Score       *int       `json:"averageScore,omitempty"`
}

type AnimeTitle struct {
	Romaji  string `json:"romaji,omitempty"`
	English string `json:"english,omitempty"`
	Native  string `json:"native,omitempty"`
}

type CoverImage struct {
	Medium string `json:"medium,omitempty"`
	Large  string `json:"large,omitempty"`
}

type PageInfo struct {
	CurrentPage int  `json:"currentPage"`
	HasNextPage bool `json:"hasNextPage"`
	LastPage    int  `json:"lastPage"`
	PerPage     int  `json:"perPage"`
	Total       int  `json:"total"`
}

type AnimePage struct {
	PageInfo PageInfo `json:"pageInfo"`
	Anime    []Anime  `json:"media"`
}

type SearchParams struct {
	Title      string
	ExternalID int
	Season     string
	SeasonYear int
	Genre      string
	Status     string
	Format     string
	Page       int
	PerPage    int
}
