package inmemory

import (
	"context"
	"fmt"
	"Go/hw6/internal/models"
	"sync"
)

type ProductsRepo struct {
	data map[int]*models.Product
	mu *sync.RWMutex
}

func (db *ProductsRepo) Create(ctx context.Context, product *models.Product) error {
	db.mu.Lock()
	defer db.mu.Unlock()
	db.data[product.ID] = product
	return nil
}

func (db *ProductsRepo) All(ctx context.Context) ([]*models.Product, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()
	products := make([]*models.Product, 0, len(db.data))
	for _, product := range db.data {
		products = append(products, product)
	}
	return products, nil
}

func (db *ProductsRepo) ByID(ctx context.Context, id int) (*models.Product, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()
	product, ok := db.data[id]
	if !ok {
		return nil, fmt.Errorf("Error %d", id)
	}
	return product, nil
}

func (db *ProductsRepo) Update(ctx context.Context, product *models.Product) error {
	db.mu.Lock()
	defer db.mu.Unlock()
	db.data[product.ID] = product
	return nil
}

func (db *ProductsRepo) Delete(ctx context.Context, id int) error {
	db.mu.Lock()
	defer db.mu.Unlock()
	delete(db.data, id)
	return nil
}