package commands

import (
	"fmt"
	"github.com/ArmandSyah/TomoPyon/anilist"
	"github.com/ArmandSyah/TomoPyon/misc"
	"github.com/bwmarrin/discordgo"
)

func searchStudio(session *discordgo.Session, message *discordgo.Message) {
	content := message.Content
	flags, studioSearchQuery := parseMessageContent(content)
	studioSearchResults := anilist.SearchStudio(studioSearchQuery)
	if studios, ok := studioSearchResults.([]anilist.Studio); ok {
		if misc.StringInSlice("v", flags) {
			//makeUserSearchEmbedsVerbose(session, message, users, studioSearchQuery)
		} else {
			makeStudioSearchEmbeds(session, message, studios, studioSearchQuery)
		}
	} else {
		_, err := session.ChannelMessageSend(message.ChannelID, "Error while searching")
		if err != nil {
			panic(err)
		}
	}
}

func makeStudioSearchEmbeds(session *discordgo.Session, message *discordgo.Message, studios []anilist.Studio, studioSearchQuery string) {
	if len(studios) < 1 {
		_, err := session.ChannelMessageSend(message.ChannelID, "No Studios were found with that name")
		if err != nil {
			panic(err)
		}
		return
	}
	authorName, avatarURL := message.Author.Username, message.Author.AvatarURL("")
	studio := studios[0]
	embed := misc.NewEmbed()
	title := fmt.Sprintf("Search Results for: %s", studioSearchQuery)
	embed.SetTitle(title)
	embed.SetAuthor(authorName, avatarURL)
	embed.AddField(studio.Name, studio.SiteURL)
	embed.SetColor(0x1855F5)
	var embeds []*misc.Embed
	embeds = append(embeds, embed)
	sendEmbeddedMessage(session, message, embeds)
}

func init() {
	add(&command{
		execute: searchStudio,
		trigger: "studio",
		aliases: []string{"s", "searchStudio", "ss"},
		desc:    "Search for studios using Anilist's Awesome API",
	})
}
