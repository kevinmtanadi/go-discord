package handler

import (
	"go-discord/logger"

	"github.com/bwmarrin/discordgo"
)

// SendMessage : Sends a message to a text channel
func (h Handler) SendMessage(message string) {
	_, err := h.s.ChannelMessageSend(h.m.ChannelID, message)
	if err != nil {
		logger.Log("Failed to send message: " + err.Error())
		return
	}
}

// SendEmbed : Sends a embed message to a text channel
func (h Handler) SendEmbed(embed *discordgo.MessageEmbed) {
	_, err := h.s.ChannelMessageSendEmbed(h.m.ChannelID, embed)
	if err != nil {
		logger.Log("Failed to send embed: " + err.Error())
		return
	}
}
