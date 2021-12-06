package main

import (
	"context"
	"proj/internal/cache/redis_cache"
	"proj/internal/http"
	"proj/internal/store/postgres"
	"log"
)

const (
	port = ":8081"
	cacheDB = 1
	cacheExpTime = 1800
	cachePort = "localhost:6379"
)

func main() {
	urlDB := "postgres://postgres:postgres@localhost:5432/products"
	store := postgres.NewDB()
	if err := store.Connect(urlDB); err != nil {
		panic(err)
	}
	defer store.Close()


	cache := redis_cache.NewRedisCache(cachePort, cacheDB, cacheExpTime)
	//defer cache.Close()

	srv := http.NewServer(context.Background(),
		http.WithAddress(port),
		http.WithStore(store),
		http.WithCache(cache),
	)
	if err := srv.Run(); err != nil {
		log.Println(err)
	}

	srv.WaitForGracefulTermination()
}