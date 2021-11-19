package http

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	lru "github.com/hashicorp/golang-lru"
	"proj/internal/models"
	"proj/internal/store"
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

func (ac *ProductResource) Routes() chi.Router {
	r := chi.NewRouter()

	r.Post("/", ac.CreateProduct)
	r.Get("/", ac.AllProducts)
	r.Get("/{id}", ac.ByID)
	r.Put("/", ac.UpdateProduct)
	r.Delete("/{id}", ac.DeleteProduct)

	return r
}

func (cr *ProductResource) CreateProduct(w http.ResponseWriter, r *http.Request) {
	product := new(models.Product)
	if err := json.NewDecoder(r.Body).Decode(product); err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		fmt.Fprintf(w, "Unknown err: %v", err)
		return
	}

	if err := cr.store.Products().Create(r.Context(), product); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "DB err: %v", err)
		return
	}

	// Правильно пройтись по всем буквам и всем словам
	cr.cache.Purge() // в рамках учебного проекта полностью чистим кэш после создания новой категории

	w.WriteHeader(http.StatusCreated)
}

func (cr *ProductResource) AllProducts(w http.ResponseWriter, r *http.Request) {
	queryValues := r.URL.Query()
	filter := &models.ProductFilter{}

	searchQuery := queryValues.Get("query")
	if searchQuery != "" {
		productsFromCache, ok := cr.cache.Get(searchQuery)
		if ok {
			render.JSON(w, r, productsFromCache)
			return
		}

		filter.Query = &searchQuery
	}

	products, err := cr.store.Products().All(r.Context(), filter)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "DB err: %v", err)
		return
	}

	if searchQuery != "" {
		cr.cache.Add(searchQuery, products)
	}
	render.JSON(w, r, products)
}

func (cr *ProductResource) ByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Unknown err: %v", err)
		return
	}

	productFromCache, ok := cr.cache.Get(id)
	if ok {
		render.JSON(w, r, productFromCache)
		return
	}

	product, err := cr.store.Products().ByID(r.Context(), id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "DB err: %v", err)
		return
	}

	cr.cache.Add(id, product)
	render.JSON(w, r, product)
}

func (cr *ProductResource) UpdateProduct(w http.ResponseWriter, r *http.Request) {
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

	if err := cr.store.Products().Update(r.Context(), product); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "DB err: %v", err)
		return
	}

	cr.cache.Remove(product.ID)
}

func (cr *ProductResource) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Unknown err: %v", err)
		return
	}

	if err := cr.store.Products().Delete(r.Context(), id); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "DB err: %v", err)
		return
	}

	cr.cache.Remove(id)
}
