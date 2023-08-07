package main

import (
	"flag"
	"fmt"
	"go-discord/handler"
	"go-discord/logger"
	"go-discord/service"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/subosito/gotenv"
)

func main() {
	discord, err := discordgo.New("Bot " + os.Getenv("DISCORD_TOKEN"))
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

	fmt.Println(time.Now())

	// Handle messages sent to Discord
	// When adding a new command to handle, add the function onto handler package
	h.Await("hello", h.SayHello)
	h.Await("join", h.JoinVoiceChannel)

	// Daily call functions at 08:00 AM
	service.DailyCall(discord)

	// Keep the bot alive until stopped
	Loop()
}

var projectFolder *string

func init() {
	projectFolder = flag.String("folder", "./", "absolute path of project folder")
	flag.Parse()

	if *projectFolder == "" {
		// Get file path
		dir, err := os.Getwd()
		if err != nil {
			log.Fatal(err)
		}
		a := "%s" + dir + "/"
		projectFolder = &a
	}

	// Load ENV Config
	gotenv.Load(*projectFolder + ".env")
	log.Printf("Project initialized")
}

func Loop() {
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM)
	<-sc
}
