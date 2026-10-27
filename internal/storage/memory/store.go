package memory

import (
	"fmt"
	"sort"

	"github.com/guii55/music-library/internal/track"
)

type Store struct {
	m      map[int64]track.Track
	nextID int64
}

func NewStore() *Store {
	m := make(map[int64]track.Track)
	return &Store{
		m:      m,
		nextID: 1,
	}
}

func (s *Store) Create(t track.Track) (track.Track, error) {
	err := t.Validate()
	if err != nil {
		return track.Track{}, fmt.Errorf("create: %w", err)
	}
	t.ID = s.nextID
	s.m[s.nextID] = t
	s.nextID++
	return t, nil
}

func (s *Store) Get(id int64) (track.Track, error) {
	t, exists := s.m[id]
	if !exists {
		return track.Track{}, fmt.Errorf("get: %w: %d", track.ErrNotFound, id)
	}
	return t, nil
}

func (s *Store) List() []track.Track {
	tracks := make([]track.Track, 0, len(s.m))
	for _, t := range s.m {
		tracks = append(tracks, t)
	}
	sort.Slice(tracks, func(i, j int) bool {
		return tracks[i].ID < tracks[j].ID
	})
	return tracks
}

func (s *Store) Update(t track.Track) (track.Track, error) {
	err := t.Validate()
	if err != nil {
		return track.Track{}, fmt.Errorf("update: %w", err)
	}
	_, exists := s.m[t.ID]
	if !exists {
		return track.Track{}, fmt.Errorf("update: %w: %d", track.ErrNotFound, t.ID)
	}
	s.m[t.ID] = t
	return t, nil
}

func (s *Store) Delete(id int64) error {
	_, exists := s.m[id]
	if !exists {
		return fmt.Errorf("delete: %w: %d", track.ErrNotFound, id)
	}
	delete(s.m, id)
	return nil
}
