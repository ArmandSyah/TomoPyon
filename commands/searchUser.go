package commands

import (
	"fmt"
	"github.com/ArmandSyah/TomoPyon/anilist"
	"github.com/ArmandSyah/TomoPyon/misc"
	"github.com/bwmarrin/discordgo"
	"strings"
	// "strconv"
	// "strings"
	"time"
)

func searchUser(session *discordgo.Session, message *discordgo.Message) {
	content := message.Content
	flags, userSearchQuery := parseMessageContent(content)
	userSearchResults := anilist.SearchUser(userSearchQuery)
	if users, ok := userSearchResults.([]anilist.User); ok {
		if misc.StringInSlice("v", flags) {
			// makeAnimeSearchEmbedsVerbose(session, message, users, userSearchQuery)
		} else {
			makeUserSearchEmbeds(session, message, users, userSearchQuery)
		}
	} else {
		_, err := session.ChannelMessageSend(message.ChannelID, "Error while searching")
		if err != nil {
			panic(err)
		}
	}
}

func makeUserSearchEmbeds(session *discordgo.Session, message *discordgo.Message, users []anilist.User, userSearchQuery string) {
	authorName, avatarURL := message.Author.Username, message.Author.AvatarURL("")
	var searchResultChunks [][]anilist.User
	var embeds []*misc.Embed
	if len(users) > 25 {
		for i := 0; i < len(users); i += 25 {
			end := i + 25
			if end > len(users) {
				end = len(users)
			}
			searchResultChunks = append(searchResultChunks, users[i:end])
		}
	} else {
		searchResultChunks = append(searchResultChunks, users[:len(users)])
	}
	for i, searchResultChunk := range searchResultChunks {
		embed := misc.NewEmbed()
		title := fmt.Sprintf("Search Results for: %s", userSearchQuery)
		embed.SetTitle(title)
		embed.SetAuthor(authorName, avatarURL)
		for _, user := range searchResultChunk {
			name, anilistLink, updatedAt := user.Name, user.SiteURL, user.UpdatedAt
			dateUpdate := time.Unix(int64(updatedAt), 0)
			key := fmt.Sprintf("%s: %s", name, anilistLink)
			value := fmt.Sprintf("Updated on: %v", dateUpdate)
			embed.AddField(key, value)
		}
		embed.SetColor(0x1855F5)
		footerMetadata := fmt.Sprintf("Current Page: %v - Total Pages: %v - Results on Page: %v - Total # of Results: %v", i+1, len(searchResultChunks), len(searchResultChunk), len(users))
		embed.SetFooter(footerMetadata, avatarURL)
		embeds = append(embeds, embed)
	}
	sendEmbeddedMessage(session, message, embeds)
}

func makeUserSearchEmbedsVerbose(session *discordgo.Session, message *discordgo.Message, users []anilist.User, userSearchQuery string) {
	authorName, avatarURL := message.Author.Username, message.Author.AvatarURL("")
	var embeds []*misc.Embed
	for i, user := range users {
		embed := misc.NewEmbed()
		title := fmt.Sprintf("Search Results for: %s", userSearchQuery)
		embed.SetTitle(title)
		embed.SetAuthor(authorName, avatarURL)
		name, anilistLink, updatedAt, avatar, about := user.Name, user.SiteURL, time.Unix(int64(user.UpdatedAt), 0), user.Avatar, user.About
		favourites, stats := user.Favourites, user.Stats
		animeFavourites, mangaFavourites, characterFavourites := favourites.Anime.Nodes, favourites.Manga.Nodes, favourites.Characters.Nodes
		watchedTime, chaptersRead, mostRecentWatch, animeStatusDistribution, mangaStatusDistribution := stats.WatchedTime, stats.ChaptersRead, stats.ActivityHistory[len(stats.ActivityHistory)-1], stats.AnimeStatusDistribution, stats.MangaStatusDistribution
		animeScoreDistribution, mangaScoreDistribution, animeListScores, mangaListScores := stats.AnimeScoreDistribution, stats.MangaScoreDistribution, stats.AnimeListScores, stats.MangaListScores
		animeFavouritesValue, mangaFavouritesValue, characterFavouritesValue := strings.Join(getTitlesFromList(animeFavourites), "|"), strings.Join(getTitlesFromList(mangaFavourites), "|"), strings.Join(getNamesFromList(characterFavourites), "|")
		embed.SetThumbnail(avatar.Medium)
		embed.SetImage(avatar.Large)
		embed.InlineAllFields()
		embed.AddField("User - LastUpdated", fmt.Sprintf("%s - %s", name, updatedAt))
		embed.AddField("Anilist Link", anilistLink)
		embed.AddField("About", about)
		embed.AddField("Favourite Anime", animeFavouritesValue)
		embed.AddField("Favourite Manga", mangaFavouritesValue)
		embed.AddField("Favourite Characters", characterFavouritesValue)
		embed.AddField("Stats", fmt.Sprintf("Total Watch Time: %s - Total Chapters Read: %s", watchedTime, chaptersRead))
		embed.AddField("Recent Activity History", fmt.Sprintf("Date: %s - Amount: %s - Level: %s", mostRecentWatch.Date, mostRecentWatch.Amount, mostRecentWatch.Level))

		// embed.AddField("Status", status)
		// embed.AddField("Format", format)
		// embed.AddField("Type", animeType)
		// embed.AddField("Start Date - End Date", dateStr)
		// embed.AddField("Season", season)
		// embed.AddField("Episode", strconv.Itoa(episodes))
		// embed.AddField("Duration", strconv.Itoa(duration))
		// embed.AddField("Genres", strings.Join(genres, ", "))
		// embed.AddField("Popularity", strconv.Itoa(popularity))
		embed.SetColor(0x1855F5)
		footerMetadata := fmt.Sprintf("Current Page: %v - Total Pages: %v", i+1, len(users))
		embed.SetFooter(footerMetadata, avatarURL)
		truncatedEmbed := embed.Truncate()
		embeds = append(embeds, truncatedEmbed)
	}
	sendEmbeddedMessage(session, message, embeds)
}

func init() {
	add(&command{
		execute: searchUser,
		trigger: "user",
		aliases: []string{"u", "searchUser", "su"},
		desc:    "Search for anime using Anilist's Awesome API",
	})
}
