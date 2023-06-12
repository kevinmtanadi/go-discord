package handler

import (
	"go-discord/logger"

	"github.com/bwmarrin/discordgo"
)

// JoinVoiceChannel : connect to a voice channel of the caller
func (h *Handler) JoinVoiceChannel(s *discordgo.Session, m *discordgo.MessageCreate) {
	vc, err := s.ChannelVoiceJoin(m.GuildID, m.ChannelID, false, false)
	if err != nil {
		logger.Log("Failed to join voice channel: " + err.Error())
		return
	}

	defer vc.Disconnect()
}
