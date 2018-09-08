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
			// makeUserSearchEmbedsVerbose(session, message, users, staffSearchQuery)
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

func init() {
	add(&command{
		execute: searchStaff,
		trigger: "staff",
		aliases: []string{"s", "searchStaff", "st"},
		desc:    "Search for staff using Anilist's Awesome API",
	})
}
