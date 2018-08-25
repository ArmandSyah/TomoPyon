package commands

import (
	// "fmt"
	// "github.com/ArmandSyah/TomoPyon/anilist"
	"github.com/ArmandSyah/TomoPyon/misc"
	"github.com/bwmarrin/discordgo"
)

const (
	animeRegex = "<.*>"
	flagRegex  = "[.*]"
)

func searchAnime(session *discordgo.Session, message *discordgo.Message) {
	content := message.Content
	animeSearchQuery := misc.ExtractSubstr(content, animeRegex)
	flags := misc.ExtractSubstr(content, flagRegex)
	animeSearchQuery = misc.TrimSides(animeSearchQuery, "<", ">")
	flags = misc.TrimSides(flags, "[", "]")
}

func init() {
	add(&command{
		execute: searchAnime,
		trigger: "anime",
		aliases: []string{"a", "searchAnime", "sa"},
		desc:    "Search for anime using Anilist's Awesome API",
	})
}
