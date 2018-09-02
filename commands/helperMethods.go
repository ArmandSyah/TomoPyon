package commands

import (
	"fmt"
	"github.com/ArmandSyah/TomoPyon/anilist"
	"github.com/ArmandSyah/TomoPyon/config"
	"github.com/ArmandSyah/TomoPyon/misc"
	"github.com/bwmarrin/discordgo"
	"strings"
)

const (
	animeRegex      = "<.*>"
	flagRegex       = `^\-.*`
	flagRemoveRegex = `(.*) (\[.*\])`
)

func parseMessageContent(content string) (flags []string, searchQuery string) {
	seperatedContent := strings.Split(content, " ")
	if len(seperatedContent) < 2 {
		return make([]string, 0), ""
	}
	seperatedContent = seperatedContent[1:]
	var startTitleIndex int
	for i, word := range seperatedContent {
		if string(word[0]) == "-" {
			flagString := word[1:]
			for _, c := range flagString {
				flags = append(flags, string(c))
			}
		} else {
			startTitleIndex = i
			break
		}
	}
	remainingQuery := seperatedContent[startTitleIndex:]
	searchQuery = strings.Join(remainingQuery, " ")
	return
	// mangaSearchQuery := strings.Join(seperatedContent[1:], " ")
	// flags := misc.ExtractSubstr(content, flagRegex)
	// flags = misc.TrimSides(flags, "[", "]")
	// flags = misc.StripWhitespace(flags)
}

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
			if reaction.UserID == config.BotID || reaction.UserID != message.Author.ID {
				return
			}
			if reaction.Emoji.Name == "➡" {
				if index == totalPages {
					index = 0
				} else {
					index++
				}
				_, err := session.ChannelMessageEditEmbed(sentMessage.ChannelID, sentMessage.ID, embeds[index].MessageEmbed)
				if err != nil {
					index--
					if index < 0 {
						index = totalPages
					}
					fmt.Println(err.Error())
					return
				}
				session.MessageReactionRemove(reaction.ChannelID, reaction.MessageID, reaction.Emoji.Name, reaction.UserID)
			} else if reaction.Emoji.Name == "⬅" {
				if index == 0 {
					index = totalPages
				} else {
					index--
				}
				_, err := session.ChannelMessageEditEmbed(sentMessage.ChannelID, sentMessage.ID, embeds[index].MessageEmbed)
				if err != nil {
					index++
					if index > totalPages {
						index = 0
					}
					fmt.Println(err.Error())
					return
				}
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
