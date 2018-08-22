package commands

import (
	"github.com/ArmandSyah/TomoPyon/config"
	"github.com/bwmarrin/discordgo"
	"log"
	"os"
	"strings"
)

var (
	CommandMap = make(map[string]*command)
	aliasMap   = make(map[string]string)
	l          = log.New(os.Stderr, "cmds: ", log.LstdFlags|log.Lshortfile)
)

type command struct {
	execute      func(*discordgo.Session, *discordgo.Message)
	trigger      string
	aliases      []string
	desc         string
	commandCount int
	deleteAfter  bool
}

func add(c *command) {
	CommandMap[c.trigger] = c
	for _, alias := range c.aliases {
		aliasMap[alias] = c.trigger
	}
	l.Printf("Added command %s | %d aliases", c.trigger, len(c.aliases))
}

func HandleCommandOnMessageCreate(session *discordgo.Session, message *discordgo.MessageCreate) {
	// defer func() {
	// 	if rec := recover(); rec != nil {
	// 		_, err := session.ChannelMessageSend(config.BotLogID, rec.(string))
	// 		if err != nil {

	// 			l.Println(err.Error())
	// 			l.Println(rec)
	// 		}
	// 	}
	// }()

	if user := message.Author; user.ID == config.BotID || user.Bot {
		return
	}

	if message.Content[0:len(config.CommandPrefix)] != config.CommandPrefix {
		return
	}

	trigger := strings.ToLower(strings.Split(message.Content, " ")[0][len(config.CommandPrefix):])
	cmd, ok := CommandMap[trigger]
	if !ok {
		cmd, ok = CommandMap[aliasMap[trigger]]
		if !ok {
			return
		}
	}
	cmd.execute(session, message.Message)
	cmd.commandCount++
	if cmd.deleteAfter {
		session.ChannelMessageDelete(message.ChannelID, message.ID)
	}
}
