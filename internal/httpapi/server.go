package httpapi

import (
	"net/http"

	"github.com/guii55/music-library/internal/storage/memory"
)

type Server struct {
	store *memory.Store
}

func NewServer(store *memory.Store) *Server {
	return &Server{store: store}
}

func (s *Server) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", s.HandlerHealth)

	return mux
}

func (s *Server) HandlerHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
