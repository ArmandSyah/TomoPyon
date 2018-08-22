package commands

import (
	"encoding/json"
	"fmt"
	"github.com/bwmarrin/discordgo"
)

func pingCommand(session *discordgo.Session, message *discordgo.Message) {
	unpackedMessage, _ := json.Marshal(message)
	fmt.Println(string(unpackedMessage))
	_, err := session.ChannelMessageSend(message.ChannelID, "I'm Alive")
	if err != nil {
		panic(err)
	}
}

func init() {
	add(&command{
		execute: pingCommand,
		trigger: "ping",
		aliases: []string{"pingme", "p"},
		desc:    "Am I alive?",
	})
}
