package main

import (
	"fmt"
	"github.com/ArmandSyah/TomoPyon/commands"
	"github.com/ArmandSyah/TomoPyon/config"
	"github.com/bwmarrin/discordgo"
)

func main() {
	config.Startup()
	config.Discord.AddHandler(commands.HandleCommandOnMessageCreate)
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
	<-make(chan struct{})
}
