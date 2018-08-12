package main

import (
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
	if token == "exit" {
		os.Exit(3)
	}
	discord, err := discordgo.New("Bot " + token)
	errorCheck("Error creating discord session", err)
	user, err := discord.User("@me")
	errorCheck("Error retrieving account", err)
	botID = user.ID
	discord.AddHandler(func(session *discordgo.Session, ready *discordgo.Ready) {
		err = session.UpdateStatus(0, "Tomoyo Ready!")
		errorCheck("Error attempting to set status", err)
		servers := session.State.Guilds
		fmt.Printf("Tomoyo is running on %d servers", len(servers))
	})

	err = discord.Open()
	errorCheck("Error opening connection to Discord", err)
	defer discord.Close()
	commandPrefix = "!"
	<-make(chan struct{})
}

func getEnv(key string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	fmt.Println("No token was found")
	return "exit"
}

func errorCheck(msg string, err error) {
	if err != nil {
		fmt.Printf("%s: $+v", msg, err)
		panic(err)
	}
}
