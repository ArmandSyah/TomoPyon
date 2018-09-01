package anilist

import (
	"github.com/ArmandSyah/TomoPyon/misc"
)

func SearchUser(title string) interface{} {
	query := `query ($search: String) {
		mangaSearchResults: Page {
			media(search: $search, type: MANGA){
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
				chapters
				volumes    
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
		mangaSearchResults := data.MangaSearchResults.Media
		for i, mangaSearchResult := range mangaSearchResults {
			mangaSearchResult.Description = misc.StripHTML(mangaSearchResult.Description)
			mangaSearchResults[i] = mangaSearchResult
		}
		return mangaSearchResults
	}
	return nil
}
