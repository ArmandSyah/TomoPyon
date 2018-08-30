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

func searchManga(session *discordgo.Session, message *discordgo.Message) {
	content := message.Content
	cleanedContent := misc.ReplaceSubstr(content, flagRemoveRegex)
	seperatedContent := strings.Split(cleanedContent, " ")
	mangaSearchQuery := strings.Join(seperatedContent[1:], " ")
	flags := misc.ExtractSubstr(content, flagRegex)
	flags = misc.TrimSides(flags, "[", "]")
	flags = misc.StripWhitespace(flags)
	animeSearchResults := anilist.SearchManga(mangaSearchQuery)
	if animes, ok := animeSearchResults.([]anilist.Media); ok {
		if flags == "v" {
			makeMangaSearchEmbedsVerbose(session, message, animes, mangaSearchQuery)
		} else {
			makeMangaSearchEmbeds(session, message, animes, mangaSearchQuery)
		}
	} else {
		_, err := session.ChannelMessageSend(message.ChannelID, "Error while searching")
		if err != nil {
			panic(err)
		}
	}
}

func makeMangaSearchEmbeds(session *discordgo.Session, message *discordgo.Message, animes []anilist.Media, mangaSearchQuery string) {
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
		title := fmt.Sprintf("Search Results for: %s", mangaSearchQuery)
		embed.SetTitle(title)
		embed.SetAuthor(authorName, avatarURL)
		for _, anime := range searchResultChunk {
			title, anilistLink, score, chapters, status, popularity := getTitle(anime.Title), anime.SiteURL, anime.MeanScore, anime.Chapters, anime.Status, anime.Popularity
			key := fmt.Sprintf("%s: %s", title, anilistLink)
			value := fmt.Sprintf("Score: %v - Chapters: %v - Status: %s - Popularity: %v", score, chapters, status, popularity)
			embed.AddField(key, value)
		}
		embed.SetColor(0x1855F5)
		footerMetadata := fmt.Sprintf("Current Page: %v - Total Pages: %v - Results on Page: %v - Total # of Results: %v", i+1, len(searchResultChunks), len(searchResultChunk), len(animes))
		embed.SetFooter(footerMetadata, avatarURL)
		embeds = append(embeds, embed)
	}
	sendEmbeddedMessage(session, message, embeds)
}

func makeMangaSearchEmbedsVerbose(session *discordgo.Session, message *discordgo.Message, animes []anilist.Media, mangaSearchQuery string) {
	authorName, avatarURL := message.Author.Username, message.Author.AvatarURL("")
	var embeds []*misc.Embed
	for i, manga := range animes {
		embed := misc.NewEmbed()
		title := fmt.Sprintf("Search Results for: %s", mangaSearchQuery)
		engTitle, romajiTitle, anilistLink, status, format := manga.Title.English, manga.Title.Romaji, manga.SiteURL, manga.Status, manga.Format
		mangaType, startDate, endDate, volumes, chapters := manga.Type, manga.StartDate, manga.EndDate, manga.Volumes, manga.Chapters
		genres, popularity, coverImages := manga.Genres, manga.Popularity, manga.CoverImage
		synonyms, malID, description := manga.Synonyms, manga.IDMal, manga.Description
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
		embed.AddField("MAL Link", fmt.Sprintf("https://myanimelist.net/manga/%v", malID))
		embed.AddField("Description", description)
		embed.AddField("Status", status)
		embed.AddField("Format", format)
		embed.AddField("Type", mangaType)
		embed.AddField("Start Date - End Date", dateStr)
		embed.AddField("Volumes", strconv.Itoa(volumes))
		embed.AddField("Chapters", strconv.Itoa(chapters))
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
		execute: searchManga,
		trigger: "manga",
		aliases: []string{"m", "searchManga", "sma"},
		desc:    "Search for manga using Anilist's Awesome API",
	})
}
