package httpapi

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/guii55/music-library/internal/storage/memory"
	"github.com/guii55/music-library/internal/track"
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
	mux.HandleFunc("POST /tracks", s.HandlerCreateTrack)
	mux.HandleFunc("GET /tracks", s.HandlerListTracks)
	mux.HandleFunc("GET /tracks/{id}", s.HandlerGetTrack)
	mux.HandleFunc("PUT /tracks/{id}", s.HandlerUpdateTrack)
	mux.HandleFunc("DELETE /tracks/{id}", s.HandlerDeleteTrack)
	return mux
}

func (s *Server) HandlerHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"status":"ok"}`))
}

func parseID(r *http.Request) (int64, error) {
	id := r.PathValue("id")
	return strconv.ParseInt(id, 10, 64)
}

func (s *Server) HandlerCreateTrack(w http.ResponseWriter, r *http.Request) {

	t := track.Track{}
	err := json.NewDecoder(r.Body).Decode(&t)
	if err != nil {

		WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}
	newTrack, err := s.store.Create(t)
	if err != nil {
		WriteError(w, err)
		return
	}

	WriteJSON(w, http.StatusCreated, newTrack)
}

func (s *Server) HandlerListTracks(w http.ResponseWriter, r *http.Request) {
	tracks := s.store.List()
	WriteJSON(w, http.StatusOK, tracks)
}

func (s *Server) HandlerGetTrack(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)

	if err != nil {
		WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid track ID"})
		return
	}

	t, err := s.store.Get(id)
	if err != nil {
		WriteError(w, err)
		return
	}
	WriteJSON(w, http.StatusOK, t)

}

func (s *Server) HandlerUpdateTrack(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)

	if err != nil {
		WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid track ID"})
		return
	}

	t := track.Track{}
	DecoderErr := json.NewDecoder(r.Body).Decode(&t)
	if DecoderErr != nil {
		WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}
	t.ID = id
	updated, err := s.store.Update(t)
	if err != nil {
		WriteError(w, err)
		return
	}
	WriteJSON(w, http.StatusOK, updated)
}

func (s *Server) HandlerDeleteTrack(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)

	if err != nil {
		WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid track ID"})
		return
	}
	erro := s.store.Delete(id)
	if erro != nil {
		WriteError(w, erro)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
