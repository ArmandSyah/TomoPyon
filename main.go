package main

import (
	"encoding/json"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"os"
)

var (
	commandPrefix string
	botID         string
	token         string
)

func main() {
	token = getEnv("TomoyoKey")
	if token == "" {
		os.Exit(3)
	}
	discord, err := discordgo.New("Bot " + token)
	errorCheck("Error creating discord session", err)
	user, err := discord.User("@me")
	errorCheck("Error retrieving account", err)
	botID = user.ID
	discord.AddHandler(pingCommand)
	discord.AddHandler(func(session *discordgo.Session, ready *discordgo.Ready) {
		err = session.UpdateStatus(0, "Tomoyo Ready!")
		errorCheck("Error attempting to set status", err)
		servers := session.State.Guilds
		fmt.Printf("Tomoyo is running on %d servers \n", len(servers))
	})

	err = discord.Open()
	errorCheck("Error opening connection to Discord", err)
	defer discord.Close()
	commandPrefix = "!"
	<-make(chan struct{})
}

func pingCommand(session *discordgo.Session, message *discordgo.MessageCreate) {
	if user := message.Author; user.ID == botID || user.Bot {
		fmt.Println("Yeah no, this'll lead to spam")
		return
	}
	if string(message.Content[0]) != commandPrefix {
		return
	}
	unpackedMessage, _ := json.Marshal(message)
	fmt.Println(string(unpackedMessage))
	_, err := session.ChannelMessageSend(message.ChannelID, "I'm Alive")
	if err != nil {
		return
	}
}

func getEnv(key string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	fmt.Println("No token was found")
	return ""
}

func errorCheck(msg string, err error) {
	if err != nil {
		fmt.Printf("%s: $+v", msg, err)
		panic(err)
	}
}
