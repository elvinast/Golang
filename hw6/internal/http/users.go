package http

import (
	"Go/hw6/internal/models"
	"Go/hw6/internal/store"
	"encoding/json"
	"fmt"
	"github.com/go-chi/render"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-chi/chi"
	lru "github.com/hashicorp/golang-lru"
	"net/http"
	"strconv"
)

type UserResource struct {
	store store.Store
	cache *lru.TwoQueueCache
}

func NewUserResource(store store.Store, cache *lru.TwoQueueCache) *UserResource {
	return &UserResource{
		store: store,
		cache: cache,
	}
}

func (pr *UserResource) Routes() chi.Router {
	r := chi.NewRouter()
	r.Post("/", pr.CreateUser)
	r.Get("/", pr.AllUsers)
	r.Get("/{id}", pr.ByID)
	r.Put("/", pr.UpdateUser)
	r.Delete("/{id}", pr.DeleteUser)
	return r
}

func (pr *UserResource) CreateUser(w http.ResponseWriter, r *http.Request) {
	user := new(models.User)
	if err := json.NewDecoder(r.Body).Decode(user); err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		fmt.Fprintf(w, "error: %v", err)
		return
	}
	if err := pr.store.Users().Create(r.Context(), user); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "DB err: %v", err)
		return
	}
	// s.store.Products().Create(r.Context(), product)
	pr.cache.Purge()
	w.WriteHeader(http.StatusCreated)
}

func (pr *UserResource) AllUsers(w http.ResponseWriter, r *http.Request) {
	queryValues := r.URL.Query()
	filter := &models.UserFilter{}

	searchQuery := queryValues.Get("query")
	if searchQuery != "" {
		usersFromCache, ok := pr.cache.Get(searchQuery)
		if ok {
			render.JSON(w, r, usersFromCache)
			return
		}
		filter.Query = &searchQuery
	}

	users, err := pr.store.Users().All(r.Context(), filter)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "DB err: %v", err)
		return
	}
	if searchQuery != "" {
		pr.cache.Add(searchQuery, users)
	}
	render.JSON(w, r, users)
}

func (pr *UserResource) ByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Unknown err: %v", err)
		return
	}
	userFromCache, ok := pr.cache.Get(id)
	if ok {
		render.JSON(w, r, userFromCache)
		return
	}
	user, err := pr.store.Users().ByID(r.Context(), id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "DB err: %v", err)
		return
	}
	pr.cache.Add(id, user)
	render.JSON(w, r, user)
}

func (pr *UserResource) UpdateUser(w http.ResponseWriter, r *http.Request) {
	user := new(models.User)
	if err := json.NewDecoder(r.Body).Decode(user); err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		fmt.Fprintf(w, "Unknown err: %v", err)
		return
	}
	err := validation.ValidateStruct(
		user,
		validation.Field(&user.ID, validation.Required),
		validation.Field(&user.Email, validation.Required),
	)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		fmt.Fprintf(w, "Unknown err: %v", err)
		return
	}
	if err := pr.store.Users().Update(r.Context(), user); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "DB err: %v", err)
		return
	}
	pr.cache.Remove(user.ID)
}

func (pr *UserResource) DeleteUser(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Unknown err: %v", err)
		return
	}
	if err := pr.store.Users().Delete(r.Context(), id); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "DB err: %v", err)
		return
	}
	pr.cache.Remove(id)
}


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