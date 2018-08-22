package commands

import (
	"fmt"
	"github.com/ArmandSyah/TomoPyon/misc"
	"github.com/bwmarrin/discordgo"
)

func helpCommand(session *discordgo.Session, message *discordgo.Message) {
	embed := misc.NewEmbed()
	embed.SetTitle("Command Listings")
	embed.SetDescription("A list of all available commands at the moment")
	for k, v := range CommandMap {
		value := fmt.Sprintf("Trigger: %s - Aliases: %v - Description: %s", v.trigger, v.aliases, v.desc)
		embed.AddField(k, value)
	}
	embed.SetColor(0x00ff00)
	session.ChannelMessageSendEmbed(message.ChannelID, embed.MessageEmbed)
}

func init() {
	add(&command{
		execute: helpCommand,
		trigger: "help",
		aliases: []string{"h"},
		desc:    "for all the help you need",
	})
}
