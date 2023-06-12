package handler

import (
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
