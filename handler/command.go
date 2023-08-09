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

	songSearched, err := service.Searchyoutube(title, h.m.GuildID)
	if err != nil {
		embed = discordgo.MessageEmbed{
			Description: fmt.Sprintf("Error searching for song **%s**", title),
		}
		h.SendEmbed(&embed)
		return
	}

	songList := song.GetSongListInstance()
	voiceChannelID := h.getUserVoiceState(h.m.GuildID, author.ID)
	songSearched.SearchQuery = title
	songSearched.VoiceChannelID = voiceChannelID
	songSearched.RequesterChannelID = h.m.ChannelID
	songSearched.Requester = author

	songList.AddSong(*songSearched, h.m.GuildID)

	embed = discordgo.MessageEmbed{
		Description: fmt.Sprintf("**Added Song**\n\n**Song Title**\n%s\n\n**Song Duration**\n%s", songSearched.Title, helper.FormatTime(songSearched.Duration)),
		Footer: &discordgo.MessageEmbedFooter{
			Text:    fmt.Sprintf("Requested by %s", author.Username),
			IconURL: author.AvatarURL("64"),
		},
	}

	h.SendEmbed(&embed)
}

func (h *Handler) StopPlayingSong(args ...interface{}) {
	songList := song.GetSongListInstance()
	songList.Clear(h.m.GuildID)

	guildID := h.s.State.Application.GuildID
	voiceConn := h.s.VoiceConnections
	if _, ok := voiceConn[guildID]; ok {
		voiceConn[guildID].Disconnect()
	}
	helper.DeleteFileExists(guildID + ".webm")
}

func (h *Handler) ClearSongQueue(args ...interface{}) {
	songList := song.GetSongListInstance()
	songList.Clear(h.m.GuildID)
}

func (h *Handler) Skip(args ...interface{}) {

}

func (h *Handler) PrintQueueList(args ...interface{}) {
	songList := song.GetSongListInstance()

	fmt.Println(h.m.GuildID)

	songString := "**Song Queue**\n==========================\n"
	for _, song := range songList.Songs[h.m.GuildID] {
		songString += fmt.Sprintf("**%s**\n - %s\n\n", song.Requester.Username, song.Title)
	}

	embed := &discordgo.MessageEmbed{
		Description: songString,
	}
	h.SendEmbed(embed)
}

func (h *Handler) getUserVoiceState(guildID string, userID string) string {
	voiceState, err := h.s.State.VoiceState(guildID, userID)
	if err != nil {
		logger.Log("Error getting voice state: " + err.Error())
		return ""
	}

	if voiceState != nil && voiceState.ChannelID != "" {
		return voiceState.ChannelID
	}

	embed := discordgo.MessageEmbed{
		Description: "User is not connected to a voice channel",
	}
	h.SendEmbed(&embed)
	return ""
}
