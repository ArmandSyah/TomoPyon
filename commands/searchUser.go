package commands

import (
	"fmt"
	"github.com/ArmandSyah/TomoPyon/anilist"
	"github.com/ArmandSyah/TomoPyon/misc"
	"github.com/bwmarrin/discordgo"
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

func init() {
	add(&command{
		execute: searchUser,
		trigger: "user",
		aliases: []string{"u", "searchUser", "su"},
		desc:    "Search for anime using Anilist's Awesome API",
	})
}
