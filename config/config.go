// Package Config provides easy start up for TomoPyon bot, and
// provides easy access to important variables
package config

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"os"
)

//Exported Variables
var (
	Discord       *discordgo.Session //current discord session instance
	Token         string             //Discord secret token for bot
	BotID         string             //Bot user's userID
	CommandPrefix string             //CommandPrefix
	BotUser       *discordgo.User    //current User struct
)

// Startup performs initialization of TomoPyon
func Startup() error {
	fmt.Println("Setting up bot")
	Token = getEnv("TomoyoKey")
	if Token == "" {
		panic("No token set in your environment variables for key \"TomoyoKey\"")
	}
	discord, err := discordgo.New("Bot " + Token)
	errorCheck("Error creating discord session", err)
	botUser, err := discord.User("@me")
	errorCheck("Error retrieving account", err)
	Discord, BotUser = discord, botUser
	BotID = BotUser.ID
	CommandPrefix = "!"
	return nil
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
