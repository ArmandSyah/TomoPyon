package anilist

func SearchStudio(title string) interface{} {
	query := `query($search: String){
		studioSearchResults: Page{
			studios(search: $search, sort: SEARCH_MATCH) {
				id
				name
				media (sort: SCORE_DESC) {
					nodes {
						id
						title {
							english
							romaji
							native
						}
						type
						status
						siteUrl
						averageScore
						meanScore
					}
				}
				siteUrl
			}
		}
	}`
	variables := map[string]string{"search": title}
	queryResults := runQuery(query, variables)
	if data, ok := queryResults.(Data); ok {
		studioSearchResults := data.StudioSearchResults.Studios
		return studioSearchResults
	}
	return nil
}
