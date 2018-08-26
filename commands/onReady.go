package commands

import (
	//"encoding/json"
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
	// out, err := json.Marshal(session.State.Ready.Guilds[0])
	// fmt.Printf(string(out))
	mainServer, _ := session.Guild(session.State.Ready.Guilds[0].ID)
	fmt.Println(mainServer.Channels)
	sentMessage, err := session.ChannelMessageSend(session.State.Guilds[0].Channels[0].ID, "Meme")
	if err != nil {
		panic(err)
	}
	err = session.MessageReactionAdd(sentMessage.ChannelID, sentMessage.ID, "arrow_left")
	if err != nil {
		panic(err)
	}
}
