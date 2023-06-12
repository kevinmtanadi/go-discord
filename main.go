package main

import (
	"fmt"
	"go-discord/constant"
	"go-discord/handler"
	"go-discord/logger"
	"go-discord/reaction"
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
	discord.AddHandler(h.MessageHandler)

	// Handle messages sent to Discord
	// When adding a new command to handle, add the function onto reaction package
	h.Await("hello", reaction.SayHello)

	// Keep the bot alive until stopped
	Loop()
}

func Loop() {
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM)
	<-sc
}
