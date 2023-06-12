package reaction

import (
	"go-discord/helper"
	"go-discord/logger"
	"strings"

	"github.com/bwmarrin/discordgo"
)

// This package is used to create reaction to a message

func SayHello(s *discordgo.Session, m *discordgo.MessageCreate, args interface{}) {
	argList := helper.ConvertInterfaceToString(args)
	name := strings.Join(argList, " ")
	SendMessage(s, m, name)
}

// SendMessage : Sends a message to a text channel
func SendMessage(s *discordgo.Session, m *discordgo.MessageCreate, message string) {
	_, err := s.ChannelMessageSend(m.ChannelID, message)
	if err != nil {
		logger.Log("Failed to send message: " + err.Error())
		return
	}
}

// SendEmbed : Sends a embed message to a text channel
func SendEmbed(s *discordgo.Session, m *discordgo.MessageCreate, embed *discordgo.MessageEmbed) {
	_, err := s.ChannelMessageSendEmbed(m.ChannelID, embed)
	if err != nil {
		logger.Log("Failed to send embed: " + err.Error())
		return
	}
}
