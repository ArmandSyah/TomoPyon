package anilist

import (
	"fmt"
	"github.com/ArmandSyah/TomoPyon/misc"
)

//SearchCharacter [title]
func SearchCharacter(title string) interface{} {
	query := `query ($search: String) {
		characterSearchResults: Page {
					characters(search: $search, sort: SEARCH_MATCH){
						id
						name {
							first
							last
							native
							alternative
						}
						image {
							large
							medium
						}
						description
						siteUrl
						media {
							nodes {
								id
								title {
									romaji
									english
									native
								}
							}
						}
					}
				}
			}`
	variables := map[string]string{"search": title}
	queryResults := runQuery(query, variables)
	if data, ok := queryResults.(Data); ok {
		characterSearchResults := data.CharacterSearchResults.Characters
		fmt.Printf("chartac: %v \n", len(characterSearchResults))
		for i, characterSearchResult := range characterSearchResults {
			characterSearchResult.Description = misc.StripHTML(characterSearchResult.Description)
			characterSearchResults[i] = characterSearchResult
		}
		return characterSearchResults
	}
	return nil
}
