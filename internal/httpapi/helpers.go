package httpapi

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/guii55/music-library/internal/track"
)

func WriteJSON(w http.ResponseWriter, status int, value any) {

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(value)

}

func WriteError(w http.ResponseWriter, err error) {

	if errors.Is(err, track.ErrNotFound) {
		WriteJSON(w, http.StatusNotFound, map[string]string{"error": "track not found"})
		return
	}
	if errors.Is(err, track.ErrInvalid) {
		WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "track invalid"})
		return

	}
	log.Printf("erro interno: %v", err)
	WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal server error"})

}
