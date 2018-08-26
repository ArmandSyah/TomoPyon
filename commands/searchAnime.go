package commands

import (
	"fmt"
	"github.com/ArmandSyah/TomoPyon/anilist"
	"github.com/ArmandSyah/TomoPyon/config"
	"github.com/ArmandSyah/TomoPyon/misc"
	"github.com/bwmarrin/discordgo"
)

const (
	animeRegex = "<.*>"
	flagRegex  = "[.*]"
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
	content, authorName, avatarURL := message.Content, message.Author.Username, message.Author.AvatarURL("")
	animeSearchQuery := misc.ExtractSubstr(content, animeRegex)
	//flags := misc.ExtractSubstr(content, flagRegex)
	animeSearchQuery = misc.TrimSides(animeSearchQuery, "<", ">")
	//flags = misc.TrimSides(flags, "[", "]")
	animeSearchResults := anilist.SearchAnime(animeSearchQuery)
	if animes, ok := animeSearchResults.([]anilist.Media); ok {
		if len(animes) > 25 {
			animes = animes[:25]
		}
		embed := misc.NewEmbed()
		title := fmt.Sprintf("Search Results for: %s", animeSearchQuery)
		embed.SetTitle(title)
		embed.SetAuthor(authorName, avatarURL)
		for _, anime := range animes {
			title, anilistLink, score, episodes, status, popularity := getTitle(anime.Title), anime.SiteURL, anime.MeanScore, anime.Episodes, anime.Status, anime.Popularity
			key := fmt.Sprintf("%s: %s", title, anilistLink)
			value := fmt.Sprintf("Score: %v - Eps: %v - Status: %s - Popularity: %v", score, episodes, status, popularity)
			embed.AddField(key, value)
		}
		embed.SetColor(0x1855F5)
		session.ChannelMessageSendEmbed(message.ChannelID, embed.MessageEmbed)
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
	for _, searchResultChunk := range searchResultChunks {
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
		embeds = append(embeds, embed)
		sendEmbeddedMessage(session, message, embeds)
	}
}

func sendEmbeddedMessage(session *discordgo.Session, message *discordgo.Message, embeds []*misc.Embed) {
	if len(embeds) > 1 {
		sentMessage, err := session.ChannelMessageSendEmbed(message.ChannelID, embeds[0].MessageEmbed)
		if err != nil {
			return
		}
		err = session.MessageReactionAdd(sentMessage.ChannelID, sentMessage.ID, ":arrow_left:")
		if err != nil {
			return
		}
		err = session.MessageReactionAdd(sentMessage.ChannelID, sentMessage.ID, ":arrow_right:")
		if err != nil {
			return
		}
		config.Discord.AddHandler(func(session *discordgo.Session, reaction *discordgo.MessageReactionAdd) {

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
