import (
	"context"
	"Go/hw6/internal/http"
	"Go/hw6/internal/store/postgres"
)

func main() {
	urlExample := "postgres://localhost:5432/postgres"
	store := postgres.NewDB()
	if err := store.Connect(urlExample); err != nil {
		panic(err)
	}
	defer store.Close()

	srv := http.NewServer(context.Background(), ":8080", store)
	if err := srv.Run(); err != nil {
		panic(err)
	}

	srv.WaitForGracefulTermination()
}