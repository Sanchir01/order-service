package order

import "net/http"

type Handler struct {
}

func NewHandler(s *Service) *Handler {
	return &Handler{}
}

func (h *Handler) HelloWorld(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello World"))
}
