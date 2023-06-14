package main

import (
	"fmt"
	"go-discord/constant"
	"go-discord/handler"
	"go-discord/logger"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

func main() {
	discord, err := discordgo.New("Bot " + constant.TOKEN)
	if err != nil {
		logger.Log("Unable to connect to the bot: " + err.Error())
	}

	err = discord.Open()
	if err != nil {
		logger.Log("Unable to connect to the bot: " + err.Error())
	}

	fmt.Println("Bot connected")

	// Create a handler to handle the app
	h := handler.NewHandler()
	discord.AddHandler(h.ReadCommand)

	// Handle messages sent to Discord
	// When adding a new command to handle, add the function onto handler package
	h.Await("hello", h.SayHello)
	h.Await("join", h.JoinVoiceChannel)

	// Keep the bot alive until stopped
	Loop()
}

func Loop() {
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM)
	<-sc
}
