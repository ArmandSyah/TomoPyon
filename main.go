package main

import (
	"encoding/json"
	"fmt"
	"github.com/ArmandSyah/TomoPyon/config"
	"github.com/bwmarrin/discordgo"
)

func main() {
	config.Startup()
	config.Discord.AddHandler(pingCommand)
	config.Discord.AddHandler(func(session *discordgo.Session, ready *discordgo.Ready) {
		err := session.UpdateStatus(0, "Tomoyo Ready!")
		if err != nil {
			panic(err)
		}
		servers := session.State.Guilds
		fmt.Printf("Tomoyo is running on %d servers \n", len(servers))
	})

	err := config.Discord.Open()
	if err != nil {
		panic(err)
	}
	defer config.Discord.Close()
	fmt.Println("test")
	<-make(chan struct{})
}

func pingCommand(session *discordgo.Session, message *discordgo.MessageCreate) {
	fmt.Println("Get this message?")
	if user := message.Author; user.ID == config.BotID || user.Bot {
		fmt.Println("Yeah no, this'll lead to spam")
		return
	}
	if string(message.Content[0]) != config.CommandPrefix {
		return
	}
	unpackedMessage, _ := json.Marshal(message)
	fmt.Println(string(unpackedMessage))
	_, err := session.ChannelMessageSend(message.ChannelID, "I'm Alive")
	if err != nil {
		panic(err)
	}
}
