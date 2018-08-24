package main

import (
	"fmt"
	"github.com/ArmandSyah/TomoPyon/commands"
	"github.com/ArmandSyah/TomoPyon/config"
	"github.com/ArmandSyah/TomoPyon/misc"
	"os"
)

func main() {
	fmt.Println(len(os.Args))
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
	misc.SearchAnime("cowboy")
}
