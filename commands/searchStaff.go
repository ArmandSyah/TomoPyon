package commands

import (
	"fmt"
	"github.com/ArmandSyah/TomoPyon/anilist"
	"github.com/ArmandSyah/TomoPyon/misc"
	"github.com/bwmarrin/discordgo"
	// "strings"
	// "time"
)

func searchStaff(session *discordgo.Session, message *discordgo.Message) {
	content := message.Content
	flags, staffSearchQuery := parseMessageContent(content)
	staffSearchResults := anilist.SearchStaff(staffSearchQuery)
	if staffs, ok := staffSearchResults.([]anilist.Staff); ok {
		if misc.StringInSlice("v", flags) {
			makeStaffSearchEmbedsVerbose(session, message, staffs, staffSearchQuery)
		} else {
			makeStaffSearchEmbeds(session, message, staffs, staffSearchQuery)
		}
	} else {
		_, err := session.ChannelMessageSend(message.ChannelID, "Error while searching")
		if err != nil {
			panic(err)
		}
	}
}

func makeStaffSearchEmbeds(session *discordgo.Session, message *discordgo.Message, staffs []anilist.Staff, staffSearchQuery string) {
	authorName, avatarURL := message.Author.Username, message.Author.AvatarURL("")
	var embeds []*misc.Embed
	for i, staff := range staffs {
		embed := misc.NewEmbed()
		name, image, description, siteurl := staff.Name, staff.Image, staff.Description, staff.SiteURL
		title := fmt.Sprintf("Name: %s - Link: %s", getStaffName(name), siteurl)
		embed.SetTitle(title)
		embed.SetAuthor(authorName, avatarURL)
		embed.SetThumbnail(image.Medium)
		embed.SetImage(image.Large)
		embed.SetDescription(description)
		embed.SetColor(0x1855F5)
		footerMetadata := fmt.Sprintf("Current Page: %v - Total Pages: %v", i+1, len(staffs))
		embed.SetFooter(footerMetadata, avatarURL)
		embeds = append(embeds, embed)
	}
	sendEmbeddedMessage(session, message, embeds)
}

func makeStaffSearchEmbedsVerbose(session *discordgo.Session, message *discordgo.Message, staffs []anilist.Staff, staffSearchQuery string) {
	authorName, avatarURL := message.Author.Username, message.Author.AvatarURL("")
	var embeds []*misc.Embed
	for i, staff := range staffs {
		embed := misc.NewEmbed()
		name, image, description, siteurl := staff.Name, staff.Image, staff.Description, staff.SiteURL
		language, staffMedia, characters := staff.Language, staff.StaffMedia.Nodes, staff.Characters.Nodes
		fmt.Println(len(characters))
		title := fmt.Sprintf("Name: %s - Link: %s", getStaffName(name), siteurl)
		embed.SetTitle(title)
		embed.SetAuthor(authorName, avatarURL)
		embed.SetThumbnail(image.Medium)
		embed.SetImage(image.Large)
		embed.SetDescription(description)
		embed.SetColor(0x1855F5)
		footerMetadata := fmt.Sprintf("Current Page: %v - Total Pages: %v", i+1, len(staffs))
		embed.SetFooter(footerMetadata, avatarURL)
		embed.AddField("Language", language)
		if len(staffMedia) > 0 {
			key := "Worked on the following media"
			var size int
			var value string
			if len(staffMedia) > 10 {
				size = 10
			} else {
				size = len(staffMedia)
			}
			for _, media := range staffMedia[:size] {
				value += fmt.Sprintf("Media Title: %s - Anilist Link: %s\n", getTitle(media.Title), media.SiteURL)
			}
			embed.AddField(key, value)
		}
		if len(characters) > 0 {
			key := "Played the following Characters"
			var size int
			var value string
			if len(characters) > 10 {
				size = 10
			} else {
				size = len(characters)
			}
			for _, character := range characters[:size] {
				value += fmt.Sprintf("Character Name: %s - Anilist Link: %s\n", getCharacterName(character), character.SiteURL)
			}
			embed.AddField(key, value)
		}
		embeds = append(embeds, embed)
	}
	sendEmbeddedMessage(session, message, embeds)
}

func init() {
	add(&command{
		execute: searchStaff,
		trigger: "staff",
		aliases: []string{"s", "searchStaff", "st"},
		desc:    "Search for staff using Anilist's Awesome API",
	})
}
