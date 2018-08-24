package misc

import (
	"bytes"
	"encoding/json"
	"fmt"
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
func SearchAnime(title string) {
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
	runQuery(query, variables)
}

func runQuery(query string, variables map[string]string) {
	url := "https://graphql.anilist.co"
	values := map[string]interface{}{"query": query, "variables": variables}
	jsonValue, err := json.Marshal(values)
	if err != nil {
		fmt.Println("Test1")
		return
	}
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		fmt.Println("Test2")
		return
	}
	defer resp.Body.Close()
	anilistData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err.Error())
	}
	anime, err := getResponse([]byte(anilistData))
	a, err := json.Marshal(&anime.Data.AnimeSearchResults.Media)
	if err != nil {
		fmt.Printf("There was an error encoding the json. err = %s", err)
		return
	}
	fmt.Printf("encoded json = %s\r\n", string(a))
}
