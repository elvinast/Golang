package inmemory

import (
	"Go/hw6/internal/models"
	"Go/hw6/internal/store"
	"sync"
)

type DB struct {
	products store.ProductRepository
	users store.UserRepository
	mu *sync.RWMutex
}

func NewDB() store.Store {
	return &DB{
		mu:   new(sync.RWMutex),
	}
}

func (db *DB) Products() store.ProductRepository {
	if db.products == nil {
		db.products = &ProductsRepo{
			data: make(map[int]*models.Product),
			mu:   new(sync.RWMutex),
		}
	}
	return db.products
}

func (db *DB) Users() store.UserRepository {
	if db.users == nil {
		db.users = &UsersRepo{
			data: make(map[int]*models.User),
			mu:   new(sync.RWMutex),
		}
	}
	return db.users
}