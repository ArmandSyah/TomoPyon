package main

import (
	//"strings"
	//"encoding/json"
	"fmt"
	"github.com/ArmandSyah/TomoPyon/commands"
	"github.com/ArmandSyah/TomoPyon/config"
	"github.com/ArmandSyah/TomoPyon/misc"
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
	// animeListings := anilist.SearchAnime("cowboy")
	// if animes, ok := animeListings.([]anilist.Media); ok {
	// 	for i, anime := range animes {
	// 		genres := strings.Join(anime.Genres, ", ")
	// 		fmt.Println(i)
	// 		fmt.Println("English Title: " + anime.Title.English)
	// 		fmt.Println("Description: " + anime.Description)
	// 		fmt.Println("Genres: " + genres)
	// 	}
	// } else {
	// 	fmt.Println("XDDDDDD")
	// }
	fmt.Println(misc.TrimSides("<Cowboy Bebop>", "<", ">"))
}
