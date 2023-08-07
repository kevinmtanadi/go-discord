package handler

import (
	"fmt"
	"go-discord/helper"
	"go-discord/logger"
	"go-discord/service"
	"go-discord/song"

	"github.com/bwmarrin/discordgo"
)

// Always receive arguments as interface{}
// Create the converter at helper package if not already exists

func (h *Handler) PlaySong(args ...interface{}) {
	title := helper.GetArgs(args)

	var embed discordgo.MessageEmbed
	searchResult, err := service.SearchYoutube(title)
	if err != nil {
		logger.Log("Error searching youtube: " + err.Error())
	}

	if searchResult == nil {
		embed := discordgo.MessageEmbed{
			Description: "No song found",
		}
		h.SendEmbed(&embed)
		return
	}

	songList := song.GetSongListInstance()
	songList.AddSong(*searchResult)

	embed = discordgo.MessageEmbed{
		Description: fmt.Sprintf("**Added Song**\n\n**Song Title**\n%s\n**Track Length**\n%s", searchResult.Title, helper.FormatTime(searchResult.Duration)),
	}

	h.SendEmbed(&embed)
}
