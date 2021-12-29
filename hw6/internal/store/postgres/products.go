package postgres

import (
	"Go/hw6/internal/models"
	"Go/hw6/internal/store"
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
)

func (db *DB) Products() store.ProductRepository {
	if db.products == nil {
		db.products = NewProductsRep(db.conn)
	}
	return db.products
}

type ProductRepository struct {
	conn *sqlx.DB
}

func NewProductsRep(conn *sqlx.DB) store.ProductRepository {
	return &ProductRepository{conn: conn}
}

func (c ProductRepository) Create(ctx context.Context, product *models.Product) error {
	_, err := c.conn.Exec("INSERT INTO products(id, title, description, price, isavailable, weight) VALUES ($1, $2, $3, $4, $5, $6)", product.ID, product.Title, product.Description, product.Price, product.IsAvailable, product.Weight)
	if err != nil {
		fmt.Println(err)
		return err
	} else {
		fmt.Printf("Success\n")
	}
	return nil
}

func (c ProductRepository) All(ctx context.Context, filter *models.ProductFilter) ([]*models.Product, error) {
	basicQuery := "SELECT * FROM products"
	if filter.Query != nil {
		basicQuery += " WHERE title ilike '%" + *filter.Query + "%'"
	}
	products := make([]*models.Product, 0)
	if err := c.conn.Select(&products, basicQuery); err != nil {
		return nil, err
	}
	return products, nil
}

func (c ProductRepository) ByID(ctx context.Context, id int) (*models.Product, error) {
	product := new(models.Product)
	if err := c.conn.Get(product, "SELECT id, title FROM products WHERE id=$1", id); err != nil {
		return nil, err
	}
	return product, nil
}

func (c ProductRepository) Update(ctx context.Context, product *models.Product) error {
	_, err := c.conn.Exec("UPDATE products SET title = $1 WHERE id = $2", product.Title, product.ID)
	if err != nil {
		return err
	}
	return nil
}

func (c ProductRepository) Delete(ctx context.Context, id int) error {
	_, err := c.conn.Exec("DELETE FROM products WHERE id = $1", id)
	if err != nil {
		return err
	}
	return nil
}