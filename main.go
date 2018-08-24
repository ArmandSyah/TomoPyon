package main

import (
	"encoding/json"
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
	animeListing := misc.SearchAnime("cowboy")
	if str, ok := animeListing.(misc.AnimeSearchResults); ok {
		b, _ := json.Marshal(str)
		fmt.Println(string(b))
	} else {
		fmt.Println("something went wrong")
	}
}
