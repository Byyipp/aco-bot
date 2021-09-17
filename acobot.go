package main

import (
	"../modules"
	"fmt"
	"github.com/bwmarrin/discordgo"
)

const token string = "<DISCORD BOT TOKEN>"

var BotID string

func main() {

	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	u, err := dg.User("@me")
	if err != nil {
		fmt.Println(err.Error())
	}

	BotID = u.ID


	dg.AddHandler(modules.HandReaction)

	_ = dg.Open()


	println("running")
	<-make(chan struct{})
}
