package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/guii55/music-library/internal/httpapi"
	"github.com/guii55/music-library/internal/storage/memory"
)

func main() {
	store := memory.NewStore()
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
	log.Fatal(server.ListenAndServe())
}
