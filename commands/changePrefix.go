package commands

import (
	"fmt"
	"github.com/ArmandSyah/TomoPyon/config"
	"github.com/bwmarrin/discordgo"
	"strings"
)

func changePrefixCommand(session *discordgo.Session, message *discordgo.Message) {
	splitMessage := strings.Split(message.Content, " ")
	newPrefix := splitMessage[len(splitMessage)-1]
	config.CommandPrefix = newPrefix
	fmt.Println("New Prefix: " + newPrefix)
	_, err := session.ChannelMessageSend(message.ChannelID, "Prefix changed to "+newPrefix)
	if err != nil {
		panic(err)
	}
}

func init() {
	add(&command{
		execute: changePrefixCommand,
		trigger: "changePrefix",
		aliases: []string{"change", "cp"},
		desc:    "change handler prefix",
	})
}
