package postgres

import (
	"context"
	"database/sql"
	"os"
	"testing"

	"github.com/guii55/music-library/internal/track"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func newTestStore(t *testing.T) *Store {
	t.Helper()
	url := os.Getenv("TEST_DATABASE_URL")
	if url == "" {
		t.Skip("TEST_DATABASE_URL nao definida; pulando teste de integracao")
	}
	db, err := sql.Open("pgx", url)
	if err != nil {
		t.Fatalf("abrir conexao: %v", err)
	}
	t.Cleanup(func() { db.Close() })
	if _, err := db.Exec("TRUNCATE tracks RESTART IDENTITY"); err != nil {
		t.Fatalf("limpar tabela: %v", err)
	}
	return New(db)
}

func TestCreateAndGet(t *testing.T) {
	store := newTestStore(t)
	ctx := context.Background()
	created, err := store.Create(ctx, track.Track{
		Title:  "Paranoid Android",
		Artist: "Radiohead",
		Album:  "OK Computer",
		Year:   1997,
		Rating: 5,
	})
	if err != nil {
		t.Fatalf("Create: %v", err)
	}
	if created.ID == 0 {
		t.Error("Create nao preencheu o ID")
	}
	if created.CreatedAt.IsZero() {
		t.Error("Create nao preencheu o CreatedAt")
	}
	got, err := store.Get(ctx, created.ID)
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if !got.CreatedAt.Equal(created.CreatedAt) {
		t.Errorf("CreatedAt = %v, esperado %v", got.CreatedAt, created.CreatedAt)
	}

	got.CreatedAt = created.CreatedAt
	if got != created {
		t.Errorf("Get = %+v, esperado %+v", got, created)
	}
}
