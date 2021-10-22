package inmemory

import (
	"context"
	"fmt"
	"Go/hw6/internal/models"
	"Go/hw6/internal/store"
	"sync"
)

type DB struct {
	data map[int]*models.FoodItem
	mu *sync.RWMutex
}

func NewDB() store.Store {
	return &DB{
		data: make(map[int]*models.FoodItem),
		mu:   new(sync.RWMutex),
	}
}

func (db *DB) Create(ctx context.Context, foodItem *models.FoodItem) error {
	db.mu.Lock()
	defer db.mu.Unlock()
	db.data[foodItem.ID] = foodItem
	return nil
}

func (db *DB) All(ctx context.Context) ([]*models.FoodItem, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	foodItems := make([]*models.FoodItem, 0, len(db.data))
	for _, foodItem := range db.data {
		foodItems = append(foodItems, foodItem)
	}

	return foodItems, nil
}

func (db *DB) ByID(ctx context.Context, id int) (*models.FoodItem, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	foodItem, ok := db.data[id]
	if !ok {
		return nil, fmt.Errorf("No product with id %d", id)
	}

	return foodItem, nil
}

func (db *DB) Update(ctx context.Context, foodItem *models.FoodItem) error {
	db.mu.Lock()
	defer db.mu.Unlock()
	db.data[foodItem.ID] = foodItem
	return nil
}

func (db *DB) Delete(ctx context.Context, id int) error {
	db.mu.Lock()
	defer db.mu.Unlock()
	delete(db.data, id)
	return nil
}
