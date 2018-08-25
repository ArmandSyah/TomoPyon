package anilist

import (
	"github.com/ArmandSyah/TomoPyon/misc"
)

type FuzzyDate struct {
	Year  int `json:"year"`
	Month int `json:"month"`
	Day   int `json:"day"`
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
}

type AnimeSearchResults struct {
	Media []Media `json:"media"`
}

type Data struct {
	AnimeSearchResults AnimeSearchResults `json:"animeSearchResults"`
}

//SearchAnime [title]
func SearchAnime(title string) interface{} {
	query := `query ($search: String) {
				animeSearchResults: Page {
					media(search: $search, type: ANIME){
						id
						idMal
						title{
							english
							romaji
							native
						}
						type
						format
						status
						description
						startDate{
							year
        					month
        					day
						}
						endDate{
							year
        					month
        					day
						}
						season      
						episodes     
						duration    
						isLicensed  
						source     
						hashtag      
						genres      
						synonyms     
						averageScore 
						meanScore    
						popularity   
						trending
						siteUrl     
					}
				}
			}`
	variables := map[string]string{"search": title}
	queryResults := runQuery(query, variables)
	if data, ok := queryResults.(Data); ok {
		animeSearchResults := data.AnimeSearchResults.Media
		for i, animeSearchResult := range animeSearchResults {
			animeSearchResult.Description = misc.StripHTML(animeSearchResult.Description)
			animeSearchResults[i] = animeSearchResult
		}
		return animeSearchResults
	}
	return nil
}
