package commands

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
)

func OnReady(session *discordgo.Session, ready *discordgo.Ready) {
	err := session.UpdateStatus(0, "Tomoyo Ready!")
	if err != nil {
		panic(err)
	}
	servers := session.State.Guilds
	fmt.Printf("Tomoyo is running on %d servers \n", len(servers))
}
