package track

import (
	"fmt"
	"time"
)

type Track struct {
	ID     int64 //atribuido pelo storage, zero antes de salvar
	Title  string
	Artist string
	Album  string
	Year   int
	Rating int //0 = sem nota; 1 a 5
}

func (t Track) Validate() error {
	if t.Title == "" {
		return fmt.Errorf("title is required: %w", ErrInvalid)
	}
	if t.Artist == "" {
		return fmt.Errorf("artist is required: %w", ErrInvalid)
	}

	if t.Year != 0 {
		if t.Year < 1900 || t.Year > time.Now().Year()+1 {
			return fmt.Errorf("year %d out of range: %w", t.Year, ErrInvalid)
		}
	}
	if t.Rating < 0 || t.Rating > 5 {
		return fmt.Errorf("rating %d out of range (0-5): %w", t.Rating, ErrInvalid)
	}
	return nil
}
