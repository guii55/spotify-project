package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/guii55/music-library/internal/track"
)

type Store struct {
	db *sql.DB
}

var _ track.Repository = (*Store)(nil)

func New(db *sql.DB) *Store {
	return &Store{
		db: db,
	}

}

func (s *Store) Create(ctx context.Context, t track.Track) (track.Track, error) {
	if err := t.Validate(); err != nil {
		return track.Track{}, fmt.Errorf("create: %w", err)
	}

	query := `
		INSERT INTO tracks (title, artist, album, year, rating)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at
	`
	err := s.db.QueryRowContext(ctx, query, t.Title, t.Artist, t.Album, t.Year, t.Rating).Scan(&t.ID, &t.CreatedAt)
	if err != nil {
		return track.Track{}, fmt.Errorf("create: %w", err)
	}
	return t, nil

}
func (s *Store) Get(ctx context.Context, id int64) (track.Track, error) {

	query := `
		SELECT id, title, artist, album, year, rating, created_at
		FROM tracks
		WHERE id = $1
	`
	var t track.Track
	err := s.db.QueryRowContext(ctx, query, id).Scan(&t.ID, &t.Title, &t.Artist, &t.Album, &t.Year, &t.Rating, &t.CreatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return track.Track{}, fmt.Errorf("get: %w: %d", track.ErrNotFound, id)
	}
	if err != nil {
		return track.Track{}, fmt.Errorf("get: %w", err)
	}
	return t, nil
}

func (s *Store) List(ctx context.Context) ([]track.Track, error) {
	query := `
		SELECT id, title, artist, album, year, rating, created_at
		FROM tracks
		ORDER BY id
	`
	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("list: %w", err)
	}
	defer rows.Close()
	tracks := make([]track.Track, 0)
	for rows.Next() {
		var t track.Track
		if err := rows.Scan(&t.ID, &t.Title, &t.Artist, &t.Album, &t.Year, &t.Rating, &t.CreatedAt); err != nil {
			return nil, fmt.Errorf("list: %w", err)
		}
		tracks = append(tracks, t)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("list: %w", err)
	}
	return tracks, nil
}
func (s *Store) Update(ctx context.Context, t track.Track) (track.Track, error) {
	if err := t.Validate(); err != nil {
		return track.Track{}, fmt.Errorf("update: %w", err)
	}

	query := `
		UPDATE tracks
		set title = $1, artist = $2, album = $3, year = $4, rating = $5
		WHERE id = $6
		RETURNING id, created_at
	`
	err := s.db.QueryRowContext(ctx, query, t.Title, t.Artist, t.Album, t.Year, t.Rating, t.ID).Scan(&t.ID, &t.CreatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return track.Track{}, fmt.Errorf("update: %w: %d", track.ErrNotFound, t.ID)
	}
	if err != nil {
		return track.Track{}, fmt.Errorf("update: %w", err)
	}
	return t, nil
}

func (s *Store) Delete(ctx context.Context, id int64) error {
	query := `
		DELETE FROM tracks
		WHERE id = $1
	`
	result, err := s.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("delete: %w", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("delete: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("delete: %w: %d", track.ErrNotFound, id)
	}
	return nil
}
