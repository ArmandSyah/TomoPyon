package commands

import (
	"fmt"
	"github.com/ArmandSyah/TomoPyon/anilist"
	"github.com/ArmandSyah/TomoPyon/misc"
	"github.com/bwmarrin/discordgo"
	"strconv"
	"strings"
	"time"
)

func searchAnime(session *discordgo.Session, message *discordgo.Message) {
	content := message.Content
	flags, animeSearchQuery := parseMessageContent(content)
	animeSearchResults := anilist.SearchAnime(animeSearchQuery)
	if animes, ok := animeSearchResults.([]anilist.Media); ok {
		if misc.StringInSlice("v", flags) {
			makeAnimeSearchEmbedsVerbose(session, message, animes, animeSearchQuery)
		} else {
			makeAnimeSearchEmbeds(session, message, animes, animeSearchQuery)
		}
	} else {
		_, err := session.ChannelMessageSend(message.ChannelID, "Error while searching")
		if err != nil {
			panic(err)
		}
	}
}

func makeAnimeSearchEmbeds(session *discordgo.Session, message *discordgo.Message, animes []anilist.Media, animeSearchQuery string) {
	authorName, avatarURL := message.Author.Username, message.Author.AvatarURL("")
	var searchResultChunks [][]anilist.Media
	var embeds []*misc.Embed
	if len(animes) > 25 {
		for i := 0; i < len(animes); i += 25 {
			end := i + 25
			if end > len(animes) {
				end = len(animes)
			}
			searchResultChunks = append(searchResultChunks, animes[i:end])
		}
	} else {
		searchResultChunks = append(searchResultChunks, animes[:len(animes)])
	}
	for i, searchResultChunk := range searchResultChunks {
		embed := misc.NewEmbed()
		title := fmt.Sprintf("Search Results for: %s", animeSearchQuery)
		embed.SetTitle(title)
		embed.SetAuthor(authorName, avatarURL)
		for _, anime := range searchResultChunk {
			title, anilistLink, score, episodes, status, popularity := getTitle(anime.Title), anime.SiteURL, anime.MeanScore, anime.Episodes, anime.Status, anime.Popularity
			key := fmt.Sprintf("%s: %s", title, anilistLink)
			value := fmt.Sprintf("Score: %v - Eps: %v - Status: %s - Popularity: %v", score, episodes, status, popularity)
			embed.AddField(key, value)
		}
		embed.SetColor(0x1855F5)
		footerMetadata := fmt.Sprintf("Current Page: %v - Total Pages: %v - Results on Page: %v - Total # of Results: %v", i+1, len(searchResultChunks), len(searchResultChunk), len(animes))
		embed.SetFooter(footerMetadata, avatarURL)
		embeds = append(embeds, embed)
	}
	sendEmbeddedMessage(session, message, embeds)
}

func makeAnimeSearchEmbedsVerbose(session *discordgo.Session, message *discordgo.Message, animes []anilist.Media, animeSearchQuery string) {
	authorName, avatarURL := message.Author.Username, message.Author.AvatarURL("")
	var embeds []*misc.Embed
	for i, anime := range animes {
		embed := misc.NewEmbed()
		title := fmt.Sprintf("Search Results for: %s", animeSearchQuery)
		engTitle, romajiTitle, anilistLink, status, format := anime.Title.English, anime.Title.Romaji, anime.SiteURL, anime.Status, anime.Format
		animeType, startDate, endDate, season, episodes := anime.Type, anime.StartDate, anime.EndDate, anime.Season, anime.Episodes
		duration, genres, popularity, coverImages := anime.Duration, anime.Genres, anime.Popularity, anime.CoverImage
		synonyms, malID, description := anime.Synonyms, anime.IDMal, anime.Description
		sd, ed := time.Date(startDate.Year, time.Month(startDate.Month), startDate.Day, 0, 0, 0, 0, time.UTC), time.Date(endDate.Year, time.Month(endDate.Month), endDate.Day, 0, 0, 0, 0, time.UTC)
		sdy, sdm, sdd := sd.Date()
		edy, edm, edd := ed.Date()
		dateStr := fmt.Sprintf("%v %v, %v - %v %v, %v", sdm, sdd, sdy, edm, edd, edy)
		embed.SetTitle(title)
		embed.SetAuthor(authorName, avatarURL)
		embed.SetThumbnail(coverImages.Medium)
		embed.SetImage(coverImages.Large)
		embed.InlineAllFields()
		embed.AddField("Names", fmt.Sprintf("English: %s - Romaji: %s", engTitle, romajiTitle))
		embed.AddField("Synonyms", strings.Join(synonyms, ", "))
		embed.AddField("Anilist Link", anilistLink)
		embed.AddField("MAL Link", fmt.Sprintf("https://myanimelist.net/anime/%v", malID))
		embed.AddField("Description", description)
		embed.AddField("Status", status)
		embed.AddField("Format", format)
		embed.AddField("Type", animeType)
		embed.AddField("Start Date - End Date", dateStr)
		embed.AddField("Season", season)
		embed.AddField("Episode", strconv.Itoa(episodes))
		embed.AddField("Duration", strconv.Itoa(duration))
		embed.AddField("Genres", strings.Join(genres, ", "))
		embed.AddField("Popularity", strconv.Itoa(popularity))
		embed.SetColor(0x1855F5)
		footerMetadata := fmt.Sprintf("Current Page: %v - Total Pages: %v", i+1, len(animes))
		embed.SetFooter(footerMetadata, avatarURL)
		truncatedEmbed := embed.Truncate()
		embeds = append(embeds, truncatedEmbed)
	}
	sendEmbeddedMessage(session, message, embeds)
}

func init() {
	add(&command{
		execute: searchAnime,
		trigger: "anime",
		aliases: []string{"a", "searchAnime", "sa"},
		desc:    "Search for anime using Anilist's Awesome API",
	})
}
