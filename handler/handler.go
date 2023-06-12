package handler

import (
	"go-discord/constant"
	"go-discord/logger"
	"strings"

	"github.com/bwmarrin/discordgo"
)

type Handler struct {
	s       *discordgo.Session
	m       *discordgo.MessageCreate
	Command map[string]func(args ...interface{})
}

func NewHandler() *Handler {
	handler := &Handler{
		Command: make(map[string]func(args ...interface{})),
	}
	return handler
}

// Await : Adds a command and the reaction to the command
func (h *Handler) Await(command string, reaction func(args ...interface{})) {
	h.Command[command] = reaction
}

// ReadCommand : Handle messages entered by user
func (h *Handler) ReadCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	message := m.Content
	if !strings.HasPrefix(message, constant.BOT_PREFIX) {
		return
	}

	h.ResetSession(s, m)
	message = strings.TrimPrefix(message, constant.BOT_PREFIX)
	contents := strings.Split(message, " ")
	command, args := contents[0], contents[1:]

	h.HandleCommmand(command, args)
}

// ResetSession : Reset session and message create whenever a new command entered
func (h *Handler) ResetSession(s *discordgo.Session, m *discordgo.MessageCreate) {
	h.s = s
	h.m = m
}

// HandleCommmand : Handle command entered by user
func (h *Handler) HandleCommmand(command string, args ...interface{}) {
	if fn, ok := h.Command[command]; ok {
		fn(args)
	} else {
		h.SendMessage("Command entered doesn't exist")
		logger.Log("Command doesn't exist: " + command)
	}
}
