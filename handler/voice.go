package handler

import (
	"go-discord/logger"
)

// JoinVoiceChannel : connect to a voice channel of the caller
func (h *Handler) JoinVoiceChannel(args ...interface{}) {
	callerID := h.m.Author.ID
	callerVoiceState, err := h.s.State.VoiceState(h.m.GuildID, callerID)
	if err != nil {
		logger.Log("Failed retrieving author ID: " + err.Error())
		return
	}

	if callerVoiceState == nil || callerVoiceState.ChannelID == "" {
		h.SendMessage("Connect to a voice channel before calling this command")
		return
	}

	_, err = h.s.ChannelVoiceJoin(h.m.GuildID, callerVoiceState.ChannelID, false, true)
	if err != nil {
		logger.Log("Failed to join voice channel: " + err.Error())
		return
	}
}
