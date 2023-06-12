package handler

import (
	"go-discord/constant"
	"go-discord/logger"
	"go-discord/reaction"
	"strings"

	"github.com/bwmarrin/discordgo"
)

type Handler struct {
	Command map[string]func(s *discordgo.Session, m *discordgo.MessageCreate, args interface{})
}

func NewHandler() *Handler {
	handler := &Handler{
		Command: make(map[string]func(s *discordgo.Session, m *discordgo.MessageCreate, args interface{})),
	}
	return handler
}

// Await : Adds a command and the reaction to the command
func (h Handler) Await(command string, reaction func(s *discordgo.Session, m *discordgo.MessageCreate, args interface{})) {
	h.Command[command] = reaction
}

// MessageHandler : Handle messages entered by user
func (h Handler) MessageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	message := m.Content
	if strings.HasPrefix(message, constant.BOT_PREFIX) {
		message = strings.TrimPrefix(message, constant.BOT_PREFIX)
		contents := strings.Split(message, " ")
		command, args := contents[0], contents[1:]
		h.HandleCommmand(s, m, command, args)
	}
}

// HandleCommmand : Handle command entered by user
func (h Handler) HandleCommmand(s *discordgo.Session, m *discordgo.MessageCreate, command string, args interface{}) {
	if fn, ok := h.Command[command]; ok {
		fn(s, m, args)
	} else {
		reaction.SendMessage(s, m, "Command entered doesn't exist")
		logger.Log("Command doesn't exist")
	}
}
