package anilist

import (
	"github.com/ArmandSyah/TomoPyon/misc"
)

func SearchUser(title string) interface{} {
	query := `query ($search: String) {
		userSearchResults: Page {
			users(search: $search) {
				id
				name
				about
				avatar {
				  large
				  medium
				}
				favourites {
				  anime {
					nodes {
					  id
					  title {
						english
						romaji
						native
					  }
					}
				  }
				  manga {
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
						alternative
					  }
					}
				  }
				}
				siteUrl
				updatedAt
			}
		}
	}`
	variables := map[string]string{"search": title}
	queryResults := runQuery(query, variables)
	if data, ok := queryResults.(Data); ok {
		userSearchResults := data.UserSearchResults.Users
		for i, userSearchResult := range userSearchResults {
			userSearchResult.About = misc.StripHTML(userSearchResult.About)
			userSearchResults[i] = userSearchResult
		}
		return userSearchResults
	}
	return nil
}
