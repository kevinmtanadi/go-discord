package handler

// Always receive arguments as interface{}
// Create the converter at helper package if not already exists

func (h *Handler) SayHello(args ...interface{}) {
	h.SendMessage("Hello!")
}

func (h *Handler) Join(args ...interface{}) {
	h.JoinVoiceChannel()
}
