package handler

import "fmt"

// Always receive arguments as interface{}
// Create the converter at helper package if not already exists

func (h *Handler) SayHello(args ...interface{}) {
	fmt.Println(h.m)
	fmt.Println(h.s)
	h.SendMessage("Hello!")
}

func (h *Handler) Join(args ...interface{}) {
	h.JoinVoiceChannel()
}
