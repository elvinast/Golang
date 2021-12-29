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

type ProductResource struct {
	store store.Store
	cache *lru.TwoQueueCache
}

func NewProductResource(store store.Store, cache *lru.TwoQueueCache) *ProductResource {
	return &ProductResource{
		store: store,
		cache: cache,
	}
}

func (pr *ProductResource) Routes() chi.Router {
	r := chi.NewRouter()
	r.Post("/", pr.CreateProduct)
	r.Get("/", pr.AllProducts)
	r.Get("/{id}", pr.ByID)
	r.Put("/", pr.UpdateProduct)
	r.Delete("/{id}", pr.DeleteProduct)
	return r
}

func (pr *ProductResource) CreateProduct(w http.ResponseWriter, r *http.Request) {
	product := new(models.Product)
	if err := json.NewDecoder(r.Body).Decode(product); err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		fmt.Fprintf(w, "Post products error: %v", err)
		return
	}
	if err := pr.store.Products().Create(r.Context(), product); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "DB err: %v", err)
		return
	}
	// s.store.Products().Create(r.Context(), product)
	pr.cache.Purge()
	w.WriteHeader(http.StatusCreated)
}

func (pr *ProductResource) AllProducts(w http.ResponseWriter, r *http.Request) {
	queryValues := r.URL.Query()
	filter := &models.ProductFilter{}

	searchQuery := queryValues.Get("query")
	if searchQuery != "" {
		productsFromCache, ok := pr.cache.Get(searchQuery)
		if ok {
			render.JSON(w, r, productsFromCache)
			return
		}
		filter.Query = &searchQuery
	}

	products, err := pr.store.Products().All(r.Context(), filter)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "DB err: %v", err)
		return
	}

	if searchQuery != "" {
		pr.cache.Add(searchQuery, products)
	}
	render.JSON(w, r, products)
}

func (pr *ProductResource) ByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Unknown err: %v", err)
		return
	}
	productFromCache, ok := pr.cache.Get(id)
	if ok {
		render.JSON(w, r, productFromCache)
		return
	}
	product, err := pr.store.Products().ByID(r.Context(), id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "DB err: %v", err)
		return
	}
	pr.cache.Add(id, product)
	render.JSON(w, r, product)
}

func (pr *ProductResource) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	product := new(models.Product)
	if err := json.NewDecoder(r.Body).Decode(product); err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		fmt.Fprintf(w, "Unknown err: %v", err)
		return
	}
	err := validation.ValidateStruct(
		product,
		validation.Field(&product.ID, validation.Required),
		validation.Field(&product.Title, validation.Required),
	)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		fmt.Fprintf(w, "Unknown err: %v", err)
		return
	}
	if err := pr.store.Products().Update(r.Context(), product); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "DB err: %v", err)
		return
	}
	pr.cache.Remove(product.ID)
}

func (pr *ProductResource) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Unknown err: %v", err)
		return
	}
	if err := pr.store.Products().Delete(r.Context(), id); err != nil {
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