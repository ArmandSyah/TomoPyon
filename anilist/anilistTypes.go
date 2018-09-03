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

type Character struct {
	ID          int             `json:"id"`
	Name        CharacterName   `json:"name"`
	Image       Image           `json:"image"`
	Description string          `json:"description"`
	SiteURL     string          `json:"siteUrl"`
	Media       MediaConnection `json:"media"`
}

type Favourites struct {
	Anime      MediaConnection     `json:"anime"`
	Manga      MediaConnection     `json:"manga"`
	Characters CharacterConnection `json:"characters"`
}

type Image struct {
	Large  string `json:"large"`
	Medium string `json:"medium"`
}

type MediaConnection struct {
	Nodes []Media `json:"nodes"`
}

type CharacterConnection struct {
	Nodes []Character `json:"nodes"`
}

type MediaTitle struct {
	Romaji  string `json:"romaji"`
	English string `json:"english"`
	Native  string `json:"native"`
}

type Media struct {
	ID           int        `json:"id"`
	IDMal        int        `json:"idMal"`
	Title        MediaTitle `json:"title"`
	Type         string     `json:"type"`
	Format       string     `json:"format"`
	Status       string     `json:"status"`
	Description  string     `json:"description"`
	StartDate    FuzzyDate  `json:"startDate"`
	EndDate      FuzzyDate  `json:"endDate"`
	Season       string     `json:"season"`
	Episodes     int        `json:"episodes"`
	Chapters     int        `json:"chapters"`
	Volumes      int        `json:"volumes"`
	Duration     int        `json:"duration"`
	IsLicensed   bool       `json:"isLicensed"`
	Source       string     `json:"source"`
	Hashtag      string     `json:"hashtag"`
	Genres       []string   `json:"genres"`
	Synonyms     []string   `json:"synonyms"`
	AverageScore int        `json:"averageScore"`
	MeanScore    int        `json:"meanScore"`
	Popularity   int        `json:"popularity"`
	Trending     int        `json:"trending"`
	SiteURL      string     `json:"siteUrl"`
	CoverImage   Image      `json:"coverImage"`
}

type StatusDistribution struct {
	Status string `json:"status"`
	Amount int    `json:"amount"`
}

type UserActivityHistory struct {
	Date   int `json:"date"`
	Amount int `json:"amount"`
	Level  int `json:"level"`
}

type ScoreDistribution struct {
	Score  int `json:"score"`
	Amount int `json:"amount"`
}

type ListScoreStats struct {
	MeanScore         int `json:"meanScore"`
	StandardDeviation int `json:"standardDeviation"`
}

type UserStats struct {
	WatchedTime             int                   `json:"watchedTime"`
	ChaptersRead            int                   `json:"chaptersRead"`
	ActivityHistory         []UserActivityHistory `json:"activityHistory"`
	AnimeStatusDistribution []StatusDistribution  `json:"animeStatusDistribution"`
	MangaStatusDistribution []StatusDistribution  `json:"mangaStatusDistribution"`
	AnimeScoreDistribution  []ScoreDistribution   `json:"animeScoreDistrubtion"`
	MangaScoreDistribution  []ScoreDistribution   `json:"mangaScoreDistribution"`
	AnimeListScores         ListScoreStats        `json:"animeListScores"`
	MangaListScores         ListScoreStats        `json:"mangaListScores"`
}

type User struct {
	ID         int        `json:"id"`
	Name       string     `json:"name"`
	About      string     `json:"about"`
	Avatar     Image      `json:"avatar"`
	Favourites Favourites `json:"favourites"`
	Stats      UserStats  `json:"stats"`
	SiteURL    string     `json:"siteUrl"`
	UpdatedAt  int        `json:"updatedAt"`
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

type UserSearchResults struct {
	Users []User `json:"users"`
}

type Data struct {
	MangaSearchResults     MangaSearchResults     `json:"mangaSearchResults"`
	AnimeSearchResults     AnimeSearchResults     `json:"animeSearchResults"`
	CharacterSearchResults CharacterSearchResults `json:"characterSearchResults"`
	UserSearchResults      UserSearchResults      `json:"userSearchResults"`
}
