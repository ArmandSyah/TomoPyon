package commands

import (
	"fmt"
	"github.com/ArmandSyah/TomoPyon/anilist"
	"github.com/ArmandSyah/TomoPyon/misc"
	"github.com/bwmarrin/discordgo"
	"strings"
)

func searchCharacter(session *discordgo.Session, message *discordgo.Message) {
	content := message.Content
	flags, characterSearchQuery := parseMessageContent(content)
	characterSearchResults := anilist.SearchCharacter(characterSearchQuery)
	if characters, ok := characterSearchResults.([]anilist.Character); ok {
		if misc.StringInSlice("v", flags) {
			makeCharacterSearchEmbedsVerbose(session, message, characters, characterSearchQuery)
		} else {
			makeCharacterSearchEmbed(session, message, characters, characterSearchQuery)
		}
	} else {
		_, err := session.ChannelMessageSend(message.ChannelID, "Error while searching")
		if err != nil {
			panic(err)
		}
	}
}

func makeCharacterSearchEmbed(session *discordgo.Session, message *discordgo.Message, characters []anilist.Character, characterSearchQuery string) {
	authorName, avatarURL := message.Author.Username, message.Author.AvatarURL("")
	var searchResultChunks [][]anilist.Character
	var embeds []*misc.Embed
	if len(characters) > 25 {
		for i := 0; i < len(characters); i += 25 {
			end := i + 25
			if end > len(characters) {
				end = len(characters)
			}
			searchResultChunks = append(searchResultChunks, characters[i:end])
		}
	} else {
		searchResultChunks = append(searchResultChunks, characters[:len(characters)])
	}
	for i, searchResultChunk := range searchResultChunks {
		embed := misc.NewEmbed()
		title := fmt.Sprintf("Search Results for: %s", characterSearchQuery)
		embed.SetTitle(title)
		embed.SetAuthor(authorName, avatarURL)
		for _, character := range searchResultChunk {
			var appearedIn []string
			for _, media := range character.Media.Nodes {
				title := fmt.Sprintf("%s (%s)", getTitle(media.Title), media.Type)
				appearedIn = append(appearedIn, title)
			}
			var key string
			if character.Name.Last == "" {
				key = fmt.Sprintf("%s - %s", character.Name.First, character.SiteURL)
			} else {
				key = fmt.Sprintf("%s, %s - %s", character.Name.Last, character.Name.First, character.SiteURL)
			}
			value := fmt.Sprintf("Appeared in: %s", strings.Join(appearedIn, ", "))
			embed.AddField(key, value)
		}
		embed.SetColor(0x1855F5)
		footerMetadata := fmt.Sprintf("Current Page: %v - Total Pages: %v - Results on Page: %v - Total # of Results: %v", i+1, len(searchResultChunks), len(searchResultChunk), len(characters))
		embed.SetFooter(footerMetadata, avatarURL)
		embeds = append(embeds, embed)
	}
	sendEmbeddedMessage(session, message, embeds)
}

func makeCharacterSearchEmbedsVerbose(session *discordgo.Session, message *discordgo.Message, characters []anilist.Character, characterSearchQuery string) {
	authorName, avatarURL := message.Author.Username, message.Author.AvatarURL("")
	var embeds []*misc.Embed
	for i, character := range characters {
		embed := misc.NewEmbed()
		title := fmt.Sprintf("Search Results for: %s", characterSearchQuery)
		name, images, description, anilistLink, media := character.Name, character.Image, character.Description, character.SiteURL, character.Media
		alternateNames := strings.Join(name.Alternative, ", ")
		nameField := fmt.Sprintf("%s, %s\nNative: %s\nAlternate Names: %s", name.Last, name.First, name.Native, alternateNames)
		var appearedIn []string
		for _, media := range media.Nodes {
			title := fmt.Sprintf("%s (%s)", getTitle(media.Title), media.Type)
			appearedIn = append(appearedIn, title)
		}
		embed.SetTitle(title)
		embed.SetAuthor(authorName, avatarURL)
		embed.SetThumbnail(images.Medium)
		embed.SetImage(images.Large)
		embed.InlineAllFields()
		embed.AddField("Names", nameField)
		embed.AddField("Anilist Link", anilistLink)
		embed.AddField("Description", description)
		embed.AddField("Appeared in:", strings.Join(appearedIn, ", "))
		embed.SetColor(0x1855F5)
		footerMetadata := fmt.Sprintf("Current Page: %v - Total Pages: %v", i+1, len(characters))
		embed.SetFooter(footerMetadata, avatarURL)
		truncatedEmbed := embed.Truncate()
		embeds = append(embeds, truncatedEmbed)
	}
	sendEmbeddedMessage(session, message, embeds)
}

func init() {
	add(&command{
		execute: searchCharacter,
		trigger: "character",
		aliases: []string{"c", "char", "searchCharacter", "sc"},
		desc:    "Search for characters using Anilist's Awesome API",
	})
}
