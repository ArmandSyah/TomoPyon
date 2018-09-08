package anilist

import (
	"github.com/ArmandSyah/TomoPyon/misc"
)

func SearchStaff(title string) interface{} {
	query := `query($search: String, sort:SEARCH_MATCH){
		staffSearchResults: Page{
			staff(search: $search){
				id
				name {
					first
					last
					native
				}
				language
				image {
					large
					medium
				}
				description
				siteUrl
				staffMedia (sort: SCORE_DESC) {
					nodes {
						id
						title {
							english
							romaji
							native
						}
					}
				}
				characters {
					nodes {
						id
						name {
							first
							last
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
		staffSearchResults := data.StaffSearchResults.Staff
		for i, staffSearchResult := range staffSearchResults {
			staffSearchResult.Description = misc.StripHTML(staffSearchResult.Description)
			staffSearchResults[i] = staffSearchResult
		}
		return staffSearchResults
	}
	return nil
}
