package handler

import (
	"fmt"
	"go-discord/helper"
	"go-discord/logger"
	"go-discord/service"
	"sync"

	"github.com/bwmarrin/discordgo"
)

// Always receive arguments as interface{}
// Create the converter at helper package if not already exists

func (h *Handler) SayHello(args ...interface{}) {
	h.SendMessage("Hello!")
}

func (h *Handler) Join(args ...interface{}) {
	h.JoinVoiceChannel()
}

func (h *Handler) PlaySong(args ...interface{}) {
	title := helper.GetArgs(args)

	song, err := service.SearchYoutube(title)
	if err != nil {
		logger.Log("Error searching youtube: " + err.Error())
	}

	var wg sync.WaitGroup
	wg.Add(1)

	go service.DownloadAudio(song.URL, &wg)

	wg.Wait()

	embed := discordgo.MessageEmbed{
		Description: fmt.Sprintf("Currently playing **%s**", song.Title),
	}
	h.SendEmbed(&embed)
}
