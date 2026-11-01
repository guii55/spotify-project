package track

import (
	"errors"
	"testing"
	"time"

	"github.com/guii55/music-library/internal/track"
)

type validateCase struct {
	name    string
	tr      Track
	wantErr error
}

func loop(t *testing.T, tests []validateCase) {
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.tr.Validate()
			if !errors.Is(err, tc.wantErr) {
				t.Errorf("Validate() = %v, want %v", err, tc.wantErr)
			}
		})
	}
}

func TestValidate(t *testing.T) {
	tests := []validateCase{

		{"title vazio", Track{Title: "", Artist: "A"}, ErrInvalid},
		{"ano 1899", Track{Title: "S", Artist: "A", Year: 1899}, ErrInvalid},
		{"ano 0 e rating 0", Track{Title: "S", Artist: "A", Year: 0, Rating: 0}, nil},
		{"ano atual + 2", Track{Title: "S", Artist: "A", Year: time.Now().Year() + 2}, ErrInvalid},
		{"rating -1", Track{Title: "S", Artist: "A", Year: 2000, Rating: -1}, ErrInvalid},
		{"rating 6", Track{Title: "S", Artist: "A", Year: 2000, Rating: 6}, ErrInvalid},
	}

	loop(t, tests)
}

func TestCreateInvalidTrack(t *testing.T) {
	s := NewStore()

	_, err := s.Create(track.Track{})
	if !errors.Is(err, track.ErrInvalid) {
		t.Fatal("Create: erro = %v, esperado ErrInvalid", err)
	}

}
