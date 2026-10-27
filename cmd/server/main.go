package main

import (
	"log"
	"net/http"
	"time"

	"github.com/guii55/music-library/internal/httpapi"
	"github.com/guii55/music-library/internal/storage/memory"
)

func main() {
	store := memory.NewStore()
	srv := httpapi.NewServer(store)

	server := &http.Server{
		Addr:        ":8080",
		Handler:     srv.Handler(),
		ReadTimeout: 5 * time.Second,
	}
	log.Fatal(server.ListenAndServe())
}
