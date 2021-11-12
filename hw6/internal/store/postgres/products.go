package postgres

import (
	"context"
	"github.com/jmoiron/sqlx"
	"Go/hw6/internal/models"
	"Go/hw6/internal/store"
)

func (db *DB) Categories() store.ProductRepository {
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
	_, err := c.conn.Exec("INSERT INTO products(title) VALUES ($1)", product.Title)
	if err != nil {
		return err
	}

	return nil
}

func (c ProductRepository) All(ctx context.Context) ([]*models.Product, error) {
	products := make([]*models.Product, 0)
	if err := c.conn.Select(&products, "SELECT * FROM products"); err != nil {
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
	_, err := c.conn.Exec("UPDATE products SET name = $1 WHERE id = $2", product.Title, product.ID)
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