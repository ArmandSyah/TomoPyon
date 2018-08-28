package commands

import (
	"fmt"
	"github.com/ArmandSyah/TomoPyon/anilist"
	"github.com/ArmandSyah/TomoPyon/config"
	"github.com/ArmandSyah/TomoPyon/misc"
	"github.com/bwmarrin/discordgo"
	"strings"
	"time"
)

const (
	animeRegex      = "<.*>"
	flagRegex       = `\[.*\]`
	flagRemoveRegex = `(.*) (\[.*\])`
)

func getTitle(mediaTitle anilist.MediaTitle) string {
	if mediaTitle.English != "" {
		return mediaTitle.English
	} else if mediaTitle.Romaji != "" {
		return mediaTitle.Romaji
	} else if mediaTitle.Native != "" {
		return mediaTitle.Native
	} else {
		return "Why is there no title, how does this happen?"
	}
}

func searchAnime(session *discordgo.Session, message *discordgo.Message) {
	content := message.Content
	cleanedContent := misc.ReplaceSubstr(content, flagRemoveRegex)
	seperatedContent := strings.Split(cleanedContent, " ")
	animeSearchQuery := strings.Join(seperatedContent[1:], " ")
	flags := misc.ExtractSubstr(content, flagRegex)
	flags = misc.TrimSides(flags, "[", "]")
	flags = misc.StripWhitespace(flags)
	flagList := strings.Split(flags, ",")
	fmt.Println(flagList)
	animeSearchQuery = misc.TrimSides(animeSearchQuery, "<", ">")
	animeSearchResults := anilist.SearchAnime(animeSearchQuery)
	if animes, ok := animeSearchResults.([]anilist.Media); ok {
		makeAnimeSearchEmbeds(session, message, animes, animeSearchQuery)
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
	for _, anime := range animes {
		embed := misc.NewEmbed()
		title := fmt.Sprintf("Search Results for: %s", animeSearchQuery)
		engTitle, romajiTitle, japTitle, anilistLink, status, format := anime.Title.English, anime.Title.Romaji, anime.Title.Native, anime.SiteURL, anime.Status, anime.Format
		animeType, startDate, endDate, season, episodes := anime.Type, anime.StartDate, anime.EndDate, anime.Season, anime.Episodes
		duration, source, genres, popularity, coverImages := anime.Duration, anime.Source, anime.Genres, anime.Popularity, anime.CoverImage
		synonyms, malID, description := anime.Synonyms, anime.IDMal, anime.Description
		sd, ed := time.Date(startDate.Year, MonthstartDate.Month, startDate.Day, 0, 0, 0, 0, time.UTC), time.Date(endDate.Year, endDate.Month, endDate.Day, 0, 0, 0, 0, time.UTC)
		sdy, sdm, sdd := sd.Date()
		edy, edm, edd := ed.Date()
		startDateStr := fmt.Sprintf("%v %v, %v")
		embed.SetTitle(title)
		embed.SetAuthor(authorName, avatarURL)
		embed.SetThumbnail(coverImages.Medium)
		embed.SetImage(coverImages.Large)
		embed.InlineAllFields()
		embed.AddField("Names", fmt.Sprintf("English: %s - Romaji: %s - Japanese: %s", engTitle, romajiTitle, japTitle))
		embed.AddField("Synonyms", strings.Join(synonyms, ", "))
		embed.AddField("Anilist Link", anilistLink)
		embed.AddField("MAL Link", fmt.Sprintf("myanimelist.net/anime/%v", malID))
		embed.AddField("Description", description)
		embed.AddField("Status", status)
		embed.AddField("Format", format)
		embed.AddField("Type", animeType)
		embed.AddField("Start Date")
	}
	sendEmbeddedMessage(session, message, embeds)
}

func sendEmbeddedMessage(session *discordgo.Session, message *discordgo.Message, embeds []*misc.Embed) {
	if len(embeds) > 1 {
		index, totalPages := 0, len(embeds)-1
		sentMessage, err := session.ChannelMessageSendEmbed(message.ChannelID, embeds[index].MessageEmbed)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		err = session.MessageReactionAdd(sentMessage.ChannelID, sentMessage.ID, "⬅")
		if err != nil {
			fmt.Println("arrow_left error: " + err.Error())
			return
		}
		err = session.MessageReactionAdd(sentMessage.ChannelID, sentMessage.ID, "➡")
		if err != nil {
			fmt.Println("arrow_right error: " + err.Error())
			return
		}
		config.Discord.AddHandler(func(session *discordgo.Session, reaction *discordgo.MessageReactionAdd) {
			if reaction.UserID == config.BotID || reaction.UserID != sentMessage.Author.ID {
				return
			}
			if reaction.Emoji.Name == "➡" {
				if index == totalPages {
					index = 0
				} else {
					index++
				}
				session.ChannelMessageEditEmbed(sentMessage.ChannelID, sentMessage.ID, embeds[index].MessageEmbed)
				session.MessageReactionRemove(reaction.ChannelID, reaction.MessageID, reaction.Emoji.Name, reaction.UserID)
			} else if reaction.Emoji.Name == "⬅" {
				if index == 0 {
					index = totalPages
				} else {
					index--
				}
				session.ChannelMessageEditEmbed(sentMessage.ChannelID, sentMessage.ID, embeds[index].MessageEmbed)
				session.MessageReactionRemove(reaction.ChannelID, reaction.MessageID, reaction.Emoji.Name, reaction.UserID)
			}
		})
	} else {
		_, err := session.ChannelMessageSendEmbed(message.ChannelID, embeds[0].MessageEmbed)
		if err != nil {
			return
		}
	}
}

func init() {
	add(&command{
		execute: searchAnime,
		trigger: "anime",
		aliases: []string{"a", "searchAnime", "sa"},
		desc:    "Search for anime using Anilist's Awesome API",
	})
}
