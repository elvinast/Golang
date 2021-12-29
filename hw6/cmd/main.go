package main

import (
	"Go/hw6/internal/http"
	"Go/hw6/internal/store/postgres"
	"context"

	lru "github.com/hashicorp/golang-lru"
	// "log"
)

func main() {
	urlExample := "postgres://localhost:5432/postgres"
	store := postgres.NewDB()
	if err := store.Connect(urlExample); err != nil {
		panic(err)
	}
	defer store.Close()
	cache, err := lru.New2Q(6)
	if err != nil {
		panic(err)
	}
	srv := http.NewServer(context.Background(), ":8080", store, cache)
	if err := srv.Run(); err != nil {
		panic(err)
	}

	srv.WaitForGracefulTermination()
}
