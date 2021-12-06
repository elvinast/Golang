package main

import (
	"context"
	"proj/internal/cache/redis_cache"
	"proj/internal/http"
	"proj/internal/message_broker/kafka"
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
	brokers := []string{"localhost:29092"}
	broker := kafka.NewBroker(brokers, cache, "peer3")
	if err = broker.Connect(ctx); err != nil {
		panic(err)
	}
	defer broker.Close()
	srv := http.NewServer(context.Background(),
		http.WithAddress(port),
		http.WithStore(store),
		http.WithCache(cache),
		http.WithBroker(broker),
	)
	if err := srv.Run(); err != nil {
		log.Println(err)
	}

	srv.WaitForGracefulTermination()
}