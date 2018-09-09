package main

import (
	// "strings"
	//"encoding/json"
	"fmt"
	"github.com/ArmandSyah/TomoPyon/anilist"
	"github.com/ArmandSyah/TomoPyon/commands"
	"github.com/ArmandSyah/TomoPyon/config"
	"os"
)

func main() {
	if len(os.Args) > 1 {
		testing()
	} else {
		config.Startup()
		config.Discord.AddHandler(commands.HandleCommandOnMessageCreate)
		config.Discord.AddHandler(commands.OnReady)
		err := config.Discord.Open()
		if err != nil {
			panic(err)
		}
		defer config.Discord.Close()
		<-make(chan struct{})
	}
}

func testing() {
	animeListings := anilist.GetSchedule()
	if animes, ok := animeListings.([]anilist.AiringSchedule); ok {
		fmt.Printf("studios found: %v\n", len(animes))
		for i, anime := range animes {
			fmt.Println(i)
			fmt.Printf("Media ID: %d\n", anime.MediaID)
			fmt.Println("English: " + anime.Media.Title.English)
			fmt.Println("Romaji: " + anime.Media.Title.Romaji)
		}
	} else {
		fmt.Println("XDDDDDD")
	}
	// fmt.Println(misc.TrimSides("<Cowboy Bebop>", "<", ">"))
}
