package commands

import (
	"fmt"
	"github.com/ArmandSyah/TomoPyon/anilist"
	"github.com/ArmandSyah/TomoPyon/misc"
	"github.com/bwmarrin/discordgo"
	"strings"
	"time"
)

func searchUser(session *discordgo.Session, message *discordgo.Message) {
	content := message.Content
	flags, userSearchQuery := parseMessageContent(content)
	userSearchResults := anilist.SearchUser(userSearchQuery)
	if users, ok := userSearchResults.([]anilist.User); ok {
		if misc.StringInSlice("v", flags) {
			makeUserSearchEmbedsVerbose(session, message, users, userSearchQuery)
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
		watchedTime, chaptersRead, mostRecentWatch, animeStatusDistribution, mangaStatusDistribution := stats.WatchedTime, stats.ChaptersRead, stats.ActivityHistory[len(stats.ActivityHistory)-1], makeStatusDistribution(stats.AnimeStatusDistribution), makeStatusDistribution(stats.MangaStatusDistribution)
		animeScoreDistribution, mangaScoreDistribution, animeListScores, mangaListScores := makeScoreDistribution(stats.AnimeScoreDistribution), makeScoreDistribution(stats.MangaScoreDistribution), stats.AnimeListScores, stats.MangaListScores
		embed.SetThumbnail(avatar.Medium)
		embed.SetImage(avatar.Large)
		embed.InlineAllFields()
		embed.AddField("User - LastUpdated", fmt.Sprintf("%s - %s", name, updatedAt))
		embed.AddField("Anilist Link", anilistLink)
		embed.AddField("About", about)
		embed.AddField("Favourite Anime", strings.Join(getTitlesFromList(animeFavourites), " | "))
		embed.AddField("Favourite Manga", strings.Join(getTitlesFromList(mangaFavourites), " | "))
		embed.AddField("Favourite Characters", strings.Join(getNamesFromList(characterFavourites), " | "))
		embed.AddField("Stats", fmt.Sprintf("Total Watch Time: %d - Total Chapters Read: %d", watchedTime, chaptersRead))
		embed.AddField("Recent Activity History", fmt.Sprintf("Date: %s - Amount: %d - Level: %d", time.Unix(int64(mostRecentWatch.Date), 0), mostRecentWatch.Amount, mostRecentWatch.Level))
		embed.AddField("Anime Status Distribution", strings.Join(animeStatusDistribution, " | "))
		embed.AddField("Manga Status Distribution", strings.Join(mangaStatusDistribution, " | "))
		embed.AddField("Anime Score Distribution", strings.Join(animeScoreDistribution, " | "))
		embed.AddField("Manga Score Distribution", strings.Join(mangaScoreDistribution, " | "))
		embed.AddField("Anime List Scores", fmt.Sprintf("Mean Score: %d - Standard Deviation: %d", animeListScores.MeanScore, animeListScores.StandardDeviation))
		embed.AddField("Manga List Scores", fmt.Sprintf("Mean Score: %d - Standard Deviation: %d", mangaListScores.MeanScore, mangaListScores.StandardDeviation))
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
