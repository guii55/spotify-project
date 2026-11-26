package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/guii55/music-library/internal/httpapi"
	"github.com/guii55/music-library/internal/storage/memory"
	"github.com/guii55/music-library/internal/storage/postgres"
	"github.com/guii55/music-library/internal/track"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
	var store track.Repository = memory.NewStore()
	if url := os.Getenv("DATABASE_URL"); url != "" {
		db, err := sql.Open("pgx", url)
		if err != nil {
			log.Fatalf("abrir conexao com o postgres: %v", err)
		}

		defer db.Close()

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := db.PingContext(ctx); err != nil {
			log.Fatalf("conectar ao postgres: %v", err)
		}

		store = postgres.New(db)
		log.Println("o storage é Postgres")
	} else {
		log.Println("o storage esta em memoria, irao sumir dados ao reiniciar")
	}
	srv := httpapi.NewServer(store)
	port := os.Getenv("PORT")
	if port == "" {

		port = "8080"
	}

	server := &http.Server{
		Addr:              ":" + port,
		Handler:           srv.Handler(),
		ReadHeaderTimeout: 5 * time.Second,
	}
	log.Printf("ouvindo na porta: %v", port)
	log.Fatal(server.ListenAndServe())
}
