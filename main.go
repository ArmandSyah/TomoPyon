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
	animeListings := anilist.SearchCharacter("rias")
	if animes, ok := animeListings.([]anilist.Character); ok {
		fmt.Printf("characters found: %v\n", len(animes))
		for i, anime := range animes {
			fmt.Println(i)
			fmt.Println("English Title: " + anime.Name.Native)
		}
	} else {
		fmt.Println("XDDDDDD")
	}
	// fmt.Println(misc.TrimSides("<Cowboy Bebop>", "<", ">"))
}
