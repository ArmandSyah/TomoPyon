package commands

import (
	"fmt"
	"github.com/ArmandSyah/TomoPyon/anilist"
	"github.com/ArmandSyah/TomoPyon/config"
	"github.com/ArmandSyah/TomoPyon/misc"
	"github.com/bwmarrin/discordgo"
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
			fmt.Printf("Index: %v\n", index)
		})
	} else {
		_, err := session.ChannelMessageSendEmbed(message.ChannelID, embeds[0].MessageEmbed)
		if err != nil {
			return
		}
	}
}
