package memory

import (
	"context"
	"errors"
	"testing"

	"github.com/guii55/music-library/internal/track"
)

// helper: devolve uma Track que passa no Validate()
func validTrack(t *testing.T) track.Track {
	t.Helper()

	tr := track.Track{
		Title:  "Song Title",
		Artist: "Artist Name",
		Album:  "Album Name",
		Year:   2023,
		Rating: 4,
	}
	if err := tr.Validate(); err != nil {
		t.Fatalf("validTrack: track inválida: %v", err)
	}
	return tr
}

func TestCreate(t *testing.T) {
	s := NewStore()

	created, err := s.Create(context.Background(), validTrack(t))
	if err != nil {
		t.Fatalf("Create: erro inesperado: %v", err)
	}
	if created.ID != 1 {
		t.Errorf("ID = %d, esperado 1", created.ID)
	}

	second, _ := s.Create(context.Background(), validTrack(t))
	if second.ID != 2 {
		t.Errorf("segundo ID = %d, esperado 2", second.ID)
	}
}
func TestCreateInvalidTrack(t *testing.T) {
	s := NewStore()

	_, err := s.Create(context.Background(), track.Track{})
	if !errors.Is(err, track.ErrInvalid) {
		t.Fatalf("Create: erro = %v, esperado ErrInvalid", err)
	}

}

func TestGet(t *testing.T) {
	s := NewStore()
	tr, err := s.Create(context.Background(), validTrack(t))
	if err != nil {
		t.Fatalf("Create: erro inesperado: %v", err)
	}
	created, err := s.Get(context.Background(), tr.ID)
	if err != nil {
		t.Fatalf("Get: erro inesperado: %v", err)
	}
	if created.ID != tr.ID {
		t.Errorf("Get: ID = %d, esperado %d", created.ID, tr.ID)
	}
}

func TestGetNotFound(t *testing.T) {
	s := NewStore()

	_, err := s.Get(context.Background(), 999)
	if !errors.Is(err, track.ErrNotFound) {
		t.Fatalf("Get: erro = %v, esperado ErrNotFound", err)
	}

}

func TestList(t *testing.T) {
	s := NewStore()

	for i := 0; i <= 12; i++ {
		tr, _ := s.Create(context.Background(), validTrack(t))
		if tr.ID != int64(i+1) {
			t.Errorf("Create: ID = %d, esperado %d", tr.ID, i+1)
		}
	}

	list, err := s.List(context.Background())
	if err != nil {
		t.Fatalf("List: erro inesperado: %v", err)
	}

	if len(list) != 13 {
		t.Errorf("List: len = %d, esperado 13", len(list))
	}
	for i, tr := range list {
		if tr.ID != int64(i+1) {
			t.Errorf("List: ID = %d, esperado %d", tr.ID, i+1)
		}
	}
}

func TestUpdate(t *testing.T) {
	s := NewStore()
	tr, erro := s.Create(context.Background(), validTrack(t))
	if erro != nil {
		t.Fatalf("Create: erro inesperado: %v", erro)
	}
	tr.Title = "Updated Title"
	_, err := s.Update(context.Background(), tr)
	if err != nil {
		t.Fatalf("Update: erro inesperado: %v", err)
	}
	updatedGET, _ := s.Get(context.Background(), tr.ID)
	if updatedGET.Title != "Updated Title" {
		t.Errorf("Get: Title = %s, esperado 'Updated Title'", updatedGET.Title)
	}

}

func TestUpdateNotFound(t *testing.T) {

	s := NewStore()
	tr := validTrack(t)
	tr.ID = -1

	_, err := s.Update(context.Background(), tr)
	if !errors.Is(err, track.ErrNotFound) {
		t.Fatalf("Update: erro = %v, esperado ErrNotFound", err)
	}

}

func TestDelete(t *testing.T) {
	s := NewStore()
	tr, err := s.Create(context.Background(), validTrack(t))
	if err != nil {
		t.Fatalf("Create: erro inesperado: %v", err)
	}
	err = s.Delete(context.Background(), tr.ID)
	if err != nil {
		t.Fatalf("Delete: erro inesperado: %v", err)
	}
	_, err = s.Get(context.Background(), tr.ID)
	if !errors.Is(err, track.ErrNotFound) {
		t.Fatalf("Get após Delete: erro = %v, esperado ErrNotFound", err)
	}
}

func TestDeleteNotFound(t *testing.T) {

	s := NewStore()
	err := s.Delete(context.Background(), -1)
	if !errors.Is(err, track.ErrNotFound) {
		t.Fatalf("Delete: erro = %v, esperado ErrNotFound", err)
	}

}
