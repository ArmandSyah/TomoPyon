package anilist

type FuzzyDate struct {
	Year  int `json:"year"`
	Month int `json:"month"`
	Day   int `json:"day"`
}

type CharacterName struct {
	First       string   `json:"first"`
	Last        string   `json:"last"`
	Native      string   `json:"native"`
	Alternative []string `json:"alternative"`
}

type CharacterImage struct {
	Large  string `json:"large"`
	Medium string `json:"medium"`
}

type Character struct {
	ID          int             `json:"id"`
	Name        CharacterName   `json:"name"`
	Image       CharacterImage  `json:"image"`
	Description string          `json:"description"`
	SiteURL     string          `json:"siteUrl"`
	Media       MediaConnection `json:"media"`
}

type MediaConnection struct {
	Nodes []Media `json:"nodes"`
}

type MediaTitle struct {
	Romaji  string `json:"romaji"`
	English string `json:"english"`
	Native  string `json:"native"`
}

type MediaCoverImage struct {
	Large  string `json:"large"`
	Medium string `json:"medium"`
}

type Media struct {
	ID           int             `json:"id"`
	IDMal        int             `json:"idMal"`
	Title        MediaTitle      `json:"title"`
	Type         string          `json:"type"`
	Format       string          `json:"format"`
	Status       string          `json:"status"`
	Description  string          `json:"description"`
	StartDate    FuzzyDate       `json:"startDate"`
	EndDate      FuzzyDate       `json:"endDate"`
	Season       string          `json:"season"`
	Episodes     int             `json:"episodes"`
	Chapters     int             `json:"chapters"`
	Volumes      int             `json:"volumes"`
	Duration     int             `json:"duration"`
	IsLicensed   bool            `json:"isLicensed"`
	Source       string          `json:"source"`
	Hashtag      string          `json:"hashtag"`
	Genres       []string        `json:"genres"`
	Synonyms     []string        `json:"synonyms"`
	AverageScore int             `json:"averageScore"`
	MeanScore    int             `json:"meanScore"`
	Popularity   int             `json:"popularity"`
	Trending     int             `json:"trending"`
	SiteURL      string          `json:"siteUrl"`
	CoverImage   MediaCoverImage `json:"coverImage"`
}

type CharacterSearchResults struct {
	Characters []Character `json:"characters"`
}

type AnimeSearchResults struct {
	Media []Media `json:"media"`
}

type MangaSearchResults struct {
	Media []Media `json:"media"`
}

type Data struct {
	MangaSearchResults     MangaSearchResults     `json:"mangaSearchResults"`
	AnimeSearchResults     AnimeSearchResults     `json:"animeSearchResults"`
	CharacterSearchResults CharacterSearchResults `json:"characterSearchResults"`
}
