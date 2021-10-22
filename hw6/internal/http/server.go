package http

import (
	"context"
	"encoding/json"
	"fmt"
	"Go/hw6/internal/models"
	"Go/hw6/internal/store"
	"log"
	"net/http"
	"strconv"
	"time"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
)

type Server struct {
	ctx         context.Context
	idleConnsCh chan struct{}
	store       store.Store
	Address string
}

func NewServer(ctx context.Context, address string, store store.Store) *Server {
	return &Server{
		ctx:         ctx,
		idleConnsCh: make(chan struct{}),
		store:       store,
		Address: address,
	}
}

func (s *Server) basicHandler() chi.Router {
	r := chi.NewRouter()

	// REST
	// сущность/идентификатор
	// /electronics/laptops
	// /electronics/phones
	r.Post("/products", func(w http.ResponseWriter, r *http.Request) {
		foodItem := new(models.FoodItem)
		if err := json.NewDecoder(r.Body).Decode(foodItem); err != nil {
			fmt.Fprintf(w, "Unknown err: %v", err)
			return
		}
		s.store.Create(r.Context(), foodItem)
	})
	r.Get("/products", func(w http.ResponseWriter, r *http.Request) {
		foodItems, err := s.store.All(r.Context())
		if err != nil {
			fmt.Fprintf(w, "Unknown err: %v", err)
			return
		}
		render.JSON(w, r, foodItems)
	})
	r.Get("/products/{id}", func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			fmt.Fprintf(w, "Unknown err: %v", err)
			return
		}
		foodItem, err := s.store.ByID(r.Context(), id)
		if err != nil {
			fmt.Fprintf(w, "Unknown err: %v", err)
			return
		}
		render.JSON(w, r, foodItem)
	})
	r.Put("/products", func(w http.ResponseWriter, r *http.Request) {
		foodItem := new(models.FoodItem)
		if err := json.NewDecoder(r.Body).Decode(foodItem); err != nil {
			fmt.Fprintf(w, "Unknown err: %v", err)
			return
		}
		s.store.Update(r.Context(), foodItem)
	})
	r.Delete("/products/{id}", func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			fmt.Fprintf(w, "Unknown err: %v", err)
			return
		}
		s.store.Delete(r.Context(), id)
	})
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
		log.Printf("[HTTP] Got err while shutting down^ %v", err)
	}
	log.Println("[HTTP] Proccessed all idle connections")
	close(s.idleConnsCh)
}

func (s *Server) WaitForGracefulTermination() {
	// блок до записи или закрытия канала
	<-s.idleConnsCh
}
