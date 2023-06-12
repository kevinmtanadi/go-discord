package handler

import (
	"go-discord/logger"

	"github.com/bwmarrin/discordgo"
)

func (h Handler) JoinVoiceChannel(s *discordgo.Session, v *discordgo.VoiceStateUpdate) {
	if v.UserID == s.State.User.ID && v.ChannelID != "" {
		// Connect to the voice channel
		vc, err := s.ChannelVoiceJoin(v.GuildID, v.ChannelID, false, false)
		if err != nil {
			logger.Log("Failed to join voice channel: " + err.Error())
			return
		}

		defer vc.Disconnect()
	}
}
