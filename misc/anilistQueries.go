package misc

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
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
	ID          int        `json:"id"`
	IDMal       int        `json:"idMal"`
	Title       MediaTitle `json:"title"`
	Status      string     `json:"status"`
	Description string     `json:"description"`
	StartDate   FuzzyDate  `json:"startDate"`
	EndDate     FuzzyDate  `json:"endDate"`
}

type AnimeSearchResults struct {
	Media []Media `json:"media"`
}

type Data struct {
	AnimeSearchResults AnimeSearchResults `json:"animeSearchResults"`
}

type MediaAPIResponse struct {
	Data Data `json:"data"`
}

func getResponse(body []byte) (*MediaAPIResponse, error) {
	var s = new(MediaAPIResponse)
	err := json.Unmarshal(body, &s)
	if err != nil {
		panic(err)
	}
	return s, err
}

//SearchAnime [title]
func SearchAnime(title string) interface{} {
	query := `query ($search: String) {
				animeSearchResults: Page {
					media(search: $search, type: ANIME){
						id
						title{
							english
							romaji
							native
						}
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
					}
				}
			}`
	variables := map[string]string{"search": title}
	queryResults := runQuery(query, variables)
	if data, ok := queryResults.(Data); ok {
		animeSearchResults := data.AnimeSearchResults
		return animeSearchResults
	}
	return nil
}

func runQuery(query string, variables map[string]string) interface{} {
	url := "https://graphql.anilist.co"
	values := map[string]interface{}{"query": query, "variables": variables}
	jsonValue, err := json.Marshal(values)
	if err != nil {
		panic(err.Error())
	}
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		panic(err.Error())
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err.Error())
	}
	anilistData, err := getResponse([]byte(body))
	if err != nil {
		panic(err.Error())
	}
	return anilistData.Data
}
