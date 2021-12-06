package cache

import (
	"context"
	"proj/internal/models"
)

type Cache interface {
	Close() error

	Products() ProductsCacheRepo
	Users() UsersCacheRepo

	DeleteAll(ctx context.Context) error
}

type ProductsCacheRepo interface {
	Set(ctx context.Context, key string, value []*models.Product) error
	Get(ctx context.Context, key string) ([]*models.Product, error)
}

type UsersCacheRepo interface {
	Set(ctx context.Context, key string, value []*models.User) error
	Get(ctx context.Context, key string) ([]*models.User, error)
}