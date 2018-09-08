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
}

func makeScoreDistribution(scoreDistribution []anilist.ScoreDistribution) (scores []string) {
	for _, score := range scoreDistribution {
		scores = append(scores, fmt.Sprintf("%d: %d", score.Score, score.Amount))
	}
	return
}

func makeStatusDistribution(statusDistribution []anilist.StatusDistribution) (statuses []string) {
	for _, status := range statusDistribution {
		statuses = append(statuses, fmt.Sprintf("%s: %d", status.Status, status.Amount))
	}
	return
}

func getCharacterName(character anilist.Character) (chacterName string) {
	if character.Name.First != "" && character.Name.Last != "" {
		chacterName = fmt.Sprintf("%s, %s", character.Name.Last, character.Name.First)
	} else if character.Name.First != "" {
		chacterName = character.Name.First
	} else if character.Name.Last != "" {
		chacterName = character.Name.Last
	} else {
		chacterName = character.Name.Native
	}
	return
}

func getCharacterNamesFromList(characterList []anilist.Character) (characterNames []string) {
	for _, character := range characterList {
		characterNames = append(characterNames, getCharacterName(character))
	}
	return
}

func getTitlesFromList(mediaList []anilist.Media) (titles []string) {
	for _, media := range mediaList {
		titles = append(titles, getTitle(media.Title))
	}
	return
}

func getTitle(mediaTitle anilist.MediaTitle) string {
	if mediaTitle.English != "" {
		return mediaTitle.English
	} else if mediaTitle.Romaji != "" {
		return mediaTitle.Romaji
	} else if mediaTitle.Native != "" {
		return mediaTitle.Native
	} else {
		return ""
	}
}

func getStaffName(staffName anilist.StaffName) (name string) {
	if staffName.First != "" && staffName.Last != "" {
		name = fmt.Sprintf("%s, %s", staffName.Last, staffName.First)
	} else if staffName.First != "" {
		name = staffName.First
	} else if staffName.Last != "" {
		name = staffName.Last
	} else {
		name = staffName.Native
	}
	return
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
