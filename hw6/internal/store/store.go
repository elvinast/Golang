package store

import (
	"context"
	"Go/hw6/internal/models"
)

type Store interface {
	Connect(url string) error
	Close() error

	Products() ProductRepository
	Users() UserRepository
}

type ProductRepository interface {
	Create(ctx context.Context, product *models.Product) error
	All(ctx context.Context, filter *models.ProductFilter) ([]*models.Product, error)
	ByID(ctx context.Context, id int) (*models.Product, error)
	Update(ctx context.Context, product *models.Product) error
	Delete(ctx context.Context, id int) error
}

type UserRepository interface {
	Create(ctx context.Context, user *models.User) error
	All(ctx context.Context, filter *models.UserFilter) ([]*models.User, error)
	ByID(ctx context.Context, id int) (*models.User, error)
	Update(ctx context.Context, user *models.User) error
	Delete(ctx context.Context, id int) error
}