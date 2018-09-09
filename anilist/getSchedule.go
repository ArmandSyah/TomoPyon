package anilist

import (
	"fmt"
	"strconv"
	"time"
)

func GetSchedule() interface{} {
	airingSearchResults := make([]AiringSchedule, 10)
	currentTime := time.Now()
	timeEndpoint := currentTime.AddDate(0, 0, 7)
	currentTimeUnixSeconds, timeEndpointUnixSeconds := currentTime.Unix(), timeEndpoint.Unix()
	query := `query ($pageNum:Int, $airingGreater:Int, $airingLess:Int){
		airingSearchResults: Page (page: $pageNum, perPage: 50) {
			airingSchedules (airingAt_greater: $airingGreater, airingAt_lesser: $airingLess) {
				id
				airingAt
				timeUntilAiring
				mediaId
				media {
					title {
						english
						romaji
						native
					}
					siteUrl
				}
				episode
			}
		}
	}`
	for i := 1; ; i++ {
		variables := map[string]string{"pageNum": strconv.Itoa(i), "airingGreater": strconv.Itoa(int(currentTimeUnixSeconds)), "airingLess": strconv.Itoa(int(timeEndpointUnixSeconds))}
		queryResults := runQuery(query, variables)
		if data, ok := queryResults.(Data); ok {
			searchResults := data.AiringSearchResults.AiringSchedules
			if len(searchResults) <= 0 {
				break
			}
			for _, searchResult := range searchResults {
				airingSearchResults = append(airingSearchResults, searchResult)
			}
		}
		fmt.Println(i)
		time.Sleep(1 * time.Second)
	}
	return airingSearchResults
}
