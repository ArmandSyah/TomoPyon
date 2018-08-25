package commands

import (
	"fmt"
	"github.com/ArmandSyah/TomoPyon/anilist"
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

func init() {
	add(&command{
		execute: searchAnime,
		trigger: "anime",
		aliases: []string{"a", "searchAnime", "sa"},
		desc:    "Search for anime using Anilist's Awesome API",
	})
}
