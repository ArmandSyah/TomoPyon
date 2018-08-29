package anilist

import (
	"github.com/ArmandSyah/TomoPyon/misc"
)

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
						coverImage{
							large
							medium
						}    
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
