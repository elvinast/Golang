package http

import (
	// "Go/hw6/internal/models"
	"Go/hw6/internal/store"
	"context"
	// "encoding/json"
	// "fmt"
	"log"
	"net/http"
	// "strconv"
	"time"

	"github.com/go-chi/chi"
	// "github.com/go-chi/render"
	lru "github.com/hashicorp/golang-lru"
)

type Server struct {
	ctx         context.Context
	idleConnsCh chan struct{}
	store       store.Store
	cache 		*lru.TwoQueueCache

	Address string
}

func NewServer(ctx context.Context, address string, store store.Store, cache *lru.TwoQueueCache) *Server {
	return &Server{
		ctx:         ctx,
		idleConnsCh: make(chan struct{}),
		store:       store,
		cache: 		 cache,	
		Address: address,
	}
}

func (s *Server) basicHandler() chi.Router {
	r := chi.NewRouter()

	productsResource := NewProductResource(s.store, s.cache)
	r.Mount("/products", productsResource.Routes())

	usersResource := NewProductResource(s.store, s.cache)
	r.Mount("/users", usersResource.Routes())
	// r.Post("/products", func(w http.ResponseWriter, r *http.Request) {
	// 	product := new(models.Product)
	// 	if err := json.NewDecoder(r.Body).Decode(product); err != nil {
	// 		fmt.Fprintf(w, "Post products error: %v", err)
	// 		return
	// 	}
	// 	s.store.Products().Create(r.Context(), product)
	// })

	// r.Get("/products", func(w http.ResponseWriter, r *http.Request) {
	// 	queryValues := r.URL.Query()
	// 	filter := &models.ProductFilter{}
	// 	if searchQuery := queryValues.Get("query"); searchQuery != "" {
	// 		filter.Query = &searchQuery
	// 	}
	// 	product, err := s.store.Products().All(r.Context(), filter)
	// 	if err != nil {
	// 		fmt.Fprintf(w, "Get products error: %v", err)
	// 		return
	// 	}
	// 	render.JSON(w, r, product)
	// })

	// r.Get("/products/{id}", func(w http.ResponseWriter, r *http.Request) {
	// 	idStr := chi.URLParam(r, "id")
	// 	id, err := strconv.Atoi(idStr)
	// 	if err != nil {
	// 		fmt.Fprintf(w, "Get product error: %v", err)
	// 		return
	// 	}
	// 	product, err := s.store.Products().ByID(r.Context(), id)
	// 	if err != nil {
	// 		fmt.Fprintf(w, "Get product error: %v", err)
	// 		return
	// 	}
	// 	render.JSON(w, r, product)
	// })

	// r.Put("/products", func(w http.ResponseWriter, r *http.Request) {
	// 	product := new(models.Product)
	// 	if err := json.NewDecoder(r.Body).Decode(product); err != nil {
	// 		fmt.Fprintf(w, "Update product error: %v", err)
	// 		return
	// 	}
	// 	s.store.Products().Update(r.Context(), product)
	// })

	// r.Delete("/products/{id}", func(w http.ResponseWriter, r *http.Request) {
	// 	idStr := chi.URLParam(r, "id")
	// 	id, err := strconv.Atoi(idStr)
	// 	if err != nil {
	// 		fmt.Fprintf(w, "Delete product error: %v", err)
	// 		return
	// 	}
	// 	s.store.Products().Delete(r.Context(), id)
	// })


	// r.Post("/users", func(w http.ResponseWriter, r *http.Request) {
	// 	user := new(models.User)
	// 	if err := json.NewDecoder(r.Body).Decode(user); err != nil {
	// 		fmt.Fprintf(w, "Post user error: %v", err)
	// 		return
	// 	}
	// 	s.store.Users().Create(r.Context(), user)
	// })

	// r.Get("/users", func(w http.ResponseWriter, r *http.Request) {
	// 	queryValues := r.URL.Query()
	// 	filter := &models.UserFilter{}
	// 	if searchQuery := queryValues.Get("query"); searchQuery != "" {
	// 		filter.Query = &searchQuery
	// 	}
	// 	user, err := s.store.Users().All(r.Context(), filter)
	// 	if err != nil {
	// 		fmt.Fprintf(w, "Get users error: %v", err)
	// 		return
	// 	}
	// 	render.JSON(w, r, user)
	// })

	// r.Get("/users/{id}", func(w http.ResponseWriter, r *http.Request) {
	// 	idStr := chi.URLParam(r, "id")
	// 	id, err := strconv.Atoi(idStr)
	// 	if err != nil {
	// 		fmt.Fprintf(w, "Get user error: %v", err)
	// 		return
	// 	}
	// 	user, err := s.store.Users().ByID(r.Context(), id)
	// 	if err != nil {
	// 		fmt.Fprintf(w, "Get user error: %v", err)
	// 		return
	// 	}
	// 	render.JSON(w, r, user)
	// })

	// r.Put("/users", func(w http.ResponseWriter, r *http.Request) {
	// 	user := new(models.User)
	// 	if err := json.NewDecoder(r.Body).Decode(user); err != nil {
	// 		fmt.Fprintf(w, "Update user error: %v", err)
	// 		return
	// 	}
	// 	s.store.Users().Update(r.Context(), user)
	// })

	// r.Delete("/users/{id}", func(w http.ResponseWriter, r *http.Request) {
	// 	idStr := chi.URLParam(r, "id")
	// 	id, err := strconv.Atoi(idStr)
	// 	if err != nil {
	// 		fmt.Fprintf(w, "Delete user error: %v", err)
	// 		return
	// 	}
	// 	s.store.Users().Delete(r.Context(), id)
	// })

	return r
}

func (s *Server) Run() error {
	srv := &http.Server{
		Addr:         s.Address,
		Handler:      s.basicHandler(),
		ReadTimeout:  time.Second * 5,
		WriteTimeout: time.Second * 30,
	}
	go s.ListenCtxForGT(srv)
	log.Println("[HTTP] Server running on", s.Address)
	return srv.ListenAndServe()
}

func (s *Server) ListenCtxForGT(srv *http.Server) {
	<-s.ctx.Done() // блокируемся, пока контекст приложения не отменен
	if err := srv.Shutdown(context.Background()); err != nil {
		log.Printf("[HTTP] Got error while shutting down: %v", err)
	}
	log.Println("[HTTP] Proccessed all idle connections")
	close(s.idleConnsCh)
}

func (s *Server) WaitForGracefulTermination() {
	// блок до записи или закрытия канала
	<-s.idleConnsCh
}
