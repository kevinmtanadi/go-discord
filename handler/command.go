package handler

import (
	"fmt"
	"go-discord/helper"
	"go-discord/logger"
	"go-discord/service"
	"go-discord/song"
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

	searchResult, err := service.SearchYoutube(title)
	if err != nil {
		logger.Log("Error searching youtube: " + err.Error())
	}

	songList := song.GetSongListInstance()
	songList.AddSong(*searchResult)

	var wg sync.WaitGroup
	wg.Add(1)

	go service.DownloadAudio(searchResult.URL, &wg)

	wg.Wait()

	embed := discordgo.MessageEmbed{
		Description: fmt.Sprintf("Currently playing **%s**", searchResult.Title),
	}
	h.SendEmbed(&embed)
}
