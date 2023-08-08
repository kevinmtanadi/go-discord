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
	author := h.m.Author

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

	searchResult.Requester = author
	searchResult.ChannelID = h.m.ChannelID

	songList := song.GetSongListInstance()
	songList.AddSong(*searchResult)

	embed = discordgo.MessageEmbed{
		Description: fmt.Sprintf("**Added Song**\n\n**Song Title**\n%s\n\n**Track Length**\n%s", searchResult.Title, helper.FormatTime(searchResult.Duration)),
		Footer: &discordgo.MessageEmbedFooter{
			Text:    fmt.Sprintf("Requested by %s", author.Username),
			IconURL: author.AvatarURL("64"),
		},
	}

	h.SendEmbed(&embed)
}

func (h *Handler) StopPlayingSong(args ...interface{}) {
	songList := song.GetSongListInstance()
	songList.Clear()

	guildID := h.s.State.Application.GuildID
	voiceConn := h.s.VoiceConnections
	if _, ok := voiceConn[guildID]; ok {
		voiceConn[guildID].Close()
	}
	helper.DeleteFileExists("audio.webm")
}

func (h *Handler) ClearSongQueue(args ...interface{}) {
	songList := song.GetSongListInstance()
	songList.Clear()
}

func (h *Handler) Skip(args ...interface{}) {

}

func (h *Handler) PrintQueueList(args ...interface{}) {
	songList := song.GetSongListInstance()

	songString := "**Song Queue**\n"
	for _, song := range songList.Songs {
		songString += fmt.Sprintf("%s - %s\n\n", song.Requester.Username, song.Title)
	}

	embed := &discordgo.MessageEmbed{
		Description: songString,
	}
	h.SendEmbed(embed)
}
