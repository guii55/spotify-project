package httpapi

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/guii55/music-library/internal/storage/memory"
	"github.com/guii55/music-library/internal/track"
)

func TestHealth(t *testing.T) {
	server := NewServer(memory.NewStore())
	req := httptest.NewRequest("GET", "/health", nil)
	rec := httptest.NewRecorder()
	server.Handler().ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, esperado %d", rec.Code, http.StatusOK)
	}
	if rec.Body.String() != `{"status":"ok"}` {
		t.Errorf("body = %q, esperado %q", rec.Body.String(), `{"status":"ok"}`)
	}
}
func TestCreateTrack(t *testing.T) {
	server := NewServer(memory.NewStore())
	body := strings.NewReader(`{"title":"Song Title","artist":"Artist Name","album":"Album Name","year":2023,"rating":4}`)
	req := httptest.NewRequest("POST", "/tracks", body)
	rec := httptest.NewRecorder()
	server.Handler().ServeHTTP(rec, req)
	if rec.Code != http.StatusCreated {
		t.Fatalf("status = %d, esperado %d", rec.Code, http.StatusCreated)
	}

	created := track.Track{}
	err := json.NewDecoder(rec.Body).Decode(&created)
	if err != nil {
		t.Fatalf("decode da resposta: %v", err)
	}

	if created.Title != "Song Title" || created.ID != 1 {
		t.Errorf("track = %+v, esperado titulo %q e id 1", created, "Song Title")
	}

}

func TestListTracks(t *testing.T) {
	store := memory.NewStore()
	server := NewServer(store)
	track1, err := store.Create(track.Track{Title: "Song 1", Artist: "Artist 1", Year: 2023, Rating: 4})
	if err != nil {
		t.Fatalf("erro ao criar track 1: %v", err)
	}
	req := httptest.NewRequest("GET", "/tracks", nil)
	rec := httptest.NewRecorder()
	server.Handler().ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, esperado %d", rec.Code, http.StatusOK)
	}

	tracks := []track.Track{}
	err = json.NewDecoder(rec.Body).Decode(&tracks)
	if err != nil {
		t.Fatalf("decode da resposta: %v", err)
	}
	if len(tracks) != 1 || tracks[0].ID != track1.ID {
		t.Errorf("tracks = %+v, esperado 1 track com ID %d", tracks, track1.ID)
	}

}
func TestGetTrack(t *testing.T) {
	store := memory.NewStore()
	server := NewServer(store)
	track1, err := store.Create(track.Track{Title: "Song 1", Artist: "Artist 1", Year: 2023, Rating: 4})
	if err != nil {
		t.Fatalf("erro ao criar track 1: %v", err)
	}
	req := httptest.NewRequest("GET", "/tracks/1", nil)
	rec := httptest.NewRecorder()
	server.Handler().ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, esperado %d", rec.Code, http.StatusOK)
	}
	got := track.Track{}
	err = json.NewDecoder(rec.Body).Decode(&got)
	if err != nil {
		t.Fatalf("decode da resposta: %v", err)
	}
	if got.ID != track1.ID || got.Title != track1.Title {
		t.Errorf("track = %+v, esperado %+v", got, track1)
	}

}

func TestUpdateTrack(t *testing.T) {
	store := memory.NewStore()
	server := NewServer(store)
	track1, err := store.Create(track.Track{Title: "Song 1", Artist: "Artist 1", Year: 2023, Rating: 4})
	if err != nil {
		t.Fatalf("erro ao criar track 1: %v", err)
	}
	body := strings.NewReader(`{"title":"Updated Song","artist":"Updated Artist","album":"Updated Album","year":2024,"rating":5}`)
	req := httptest.NewRequest("PUT", "/tracks/1", body)
	rec := httptest.NewRecorder()
	server.Handler().ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, esperado %d", rec.Code, http.StatusOK)
	}
	got := track.Track{}
	err = json.NewDecoder(rec.Body).Decode(&got)
	if err != nil {
		t.Fatalf("decode da resposta: %v", err)
	}
	if got.ID != track1.ID || got.Title != "Updated Song" {
		t.Errorf("track = %+v, esperado ID %d e título %q", got, track1.ID, "Updated Song")
	}
}

func TestDeleteTrack(t *testing.T) {
	store := memory.NewStore()
	server := NewServer(store)
	track1, err := store.Create(track.Track{Title: "Song 1", Artist: "Artist 1", Year: 2023, Rating: 4})
	if err != nil {
		t.Fatalf("erro ao criar track 1: %v", err)
	}
	req := httptest.NewRequest("DELETE", "/tracks/1", nil)
	rec := httptest.NewRecorder()
	server.Handler().ServeHTTP(rec, req)
	if rec.Code != http.StatusNoContent {
		t.Fatalf("status = %d, esperado %d", rec.Code, http.StatusNoContent)
	}

	_, err = store.Get(track1.ID)
	if !errors.Is(err, track.ErrNotFound) {
		t.Fatalf("track com ID %d não foi deletada", track1.ID)
	}
	if rec.Body.String() != "" {
		t.Fatalf("body = %q, esperado vazio", rec.Body.String())
	}

}

func TestGetTrackInvalidID(t *testing.T) {

	store := memory.NewStore()
	server := NewServer(store)
	req := httptest.NewRequest("GET", "/tracks/abc", nil)
	rec := httptest.NewRecorder()
	server.Handler().ServeHTTP(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, esperado %d", rec.Code, http.StatusBadRequest)
	}

}

func TestGetTrackNotFound(t *testing.T) {

	store := memory.NewStore()
	server := NewServer(store)

	req := httptest.NewRequest("GET", "/tracks/999", nil)
	rec := httptest.NewRecorder()
	server.Handler().ServeHTTP(rec, req)
	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d, esperado %d", rec.Code, http.StatusNotFound)
	}
}

func TestCreateTrackInvalid(t *testing.T) {

	store := memory.NewStore()
	server := NewServer(store)
	body := strings.NewReader(`{"title":"","artist":"Artist Name","album":"Album Name","year":2023,"rating":4}`)
	req := httptest.NewRequest("POST", "/tracks", body)
	rec := httptest.NewRecorder()
	server.Handler().ServeHTTP(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, esperado %d", rec.Code, http.StatusBadRequest)
	}

}
