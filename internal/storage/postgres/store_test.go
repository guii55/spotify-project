package postgres

import (
	"context"
	"database/sql"
	"errors"
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

func TestGetNotFound(t *testing.T) {
	store := newTestStore(t)
	ctx := context.Background()
	_, err := store.Get(ctx, 999)
	if !errors.Is(err, track.ErrNotFound) {
		t.Errorf("Get(999) erro = %v, esperado ErrNotFound", err)

	}
}

func TestDelete(t *testing.T) {
	store := newTestStore(t)
	ctx := context.Background()
	created, err := store.Create(ctx, track.Track{
		Title:  "Paranoid Android",
		Artist: "Radiohead",
	})
	if err != nil {
		t.Fatalf("Create: %v", err)
	}

	err = store.Delete(ctx, created.ID)
	if err != nil {
		t.Fatalf("Delete: %v", err)
	}
	_, err = store.Get(ctx, created.ID)
	if !errors.Is(err, track.ErrNotFound) {
		t.Errorf("Get(id) erro = %v, esperado ErrNotFound", err)
	}
	err = store.Delete(ctx, created.ID)
	if !errors.Is(err, track.ErrNotFound) {
		t.Errorf("Delete repetido erro = %v, esperado ErrNotFound", err)
	}

}

func TestUpdate(t *testing.T) {
	store := newTestStore(t)
	ctx := context.Background()

	// 1. CRIAR, momento em que a track é criada.
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

	// 2. PREPARAR, oque voce quer gravar (o updated)
	want := created
	want.Title = "Karma Police"
	want.Rating = 4
	//3. ATUALIZAR, mandando o want
	updated, err := store.Update(ctx, want)
	if err != nil {
		t.Fatalf("Update: %v", err)
	}
	// 4. CONFERIR oque o update devolveu, contrao want
	if !updated.CreatedAt.Equal(want.CreatedAt) {
		t.Errorf("Update = %v, esperado %v", updated.CreatedAt, want.CreatedAt)
	}
	updated.CreatedAt = want.CreatedAt
	if updated != want {
		t.Errorf("Update = %+v, esperado %+v", updated, want)
	}
	//5. BUSCAR no banco, o got vai entrar aqui
	got, err := store.Get(ctx, created.ID)
	if err != nil {
		t.Fatalf("Get: %v", err)
	}

	//6. CONFERIR oque esta gravado, contra o want
	if !got.CreatedAt.Equal(want.CreatedAt) {
		t.Errorf("Get = %v, esperado %v", got.CreatedAt, want.CreatedAt)
	}
	got.CreatedAt = want.CreatedAt
	if got != want {
		t.Errorf("Get = %+v, esperado %+v", got, want)
	}

}

func TestUpdateNotFound(t *testing.T) {
	store := newTestStore(t)
	ctx := context.Background()
	_, err := store.Update(ctx, track.Track{
		ID:     999,
		Title:  "Karma Police",
		Artist: "Radiohead",
	})
	if !errors.Is(err, track.ErrNotFound) {
		t.Errorf("Update(999) erro = %v, esperado ErrNotFound", err)
	}
}

func TestListEmpty(t *testing.T) {
	store := newTestStore(t)
	ctx := context.Background()

	tracks, err := store.List(ctx)
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if tracks == nil {
		t.Error("List devolveu nil, esperado slice vazia")
	}
	if len(tracks) != 0 {
		t.Errorf("List devolveu %d tracks, esperado 0", len(tracks))
	}
}
func TestList(t *testing.T) {
	store := newTestStore(t)
	ctx := context.Background()

	// 1. Criar tracks
	first, err := store.Create(ctx, track.Track{
		Title:  "Paranoid Android",
		Artist: "Radiohead",
	})
	if err != nil {
		t.Fatalf("Create first: %v", err)
	}
	second, err := store.Create(ctx, track.Track{
		Title:  "Creep",
		Artist: "Radiohead",
	})
	if err != nil {
		t.Fatalf("Create second: %v", err)
	}
	//2. Listar
	tracks, err := store.List(ctx)
	if err != nil {
		t.Fatalf("List: %v", err)
	}

	//3. Conferir o tamanho
	if len(tracks) != 2 {
		t.Fatalf("List devolveu %d tracks, esperado 2", len(tracks))
	}

	//4. Conferir a ordem:

	if tracks[0].ID != first.ID {
		t.Errorf("tracks[0].ID = %d, esperado %d", tracks[0].ID, first.ID)
	}

	if tracks[1].ID != second.ID {
		t.Errorf("tracks[1].ID = %d, esperado %d", tracks[1].ID, second.ID)
	}
}
