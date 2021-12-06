package main

import (
	"encoding/json"
	"fmt"
	"proj/internal/models"
	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
)

func main() {
	//ctx := context.Background()
	urlExample := "postgres://postgres:postgres@localhost:5432/products"
	conn, err := sqlx.Connect("pgx", urlExample)
	if err != nil {
		panic(err)
	}

	if err = conn.Ping(); err != nil {
		panic(err)
	}

	products := make([]*models.Coupon, 0)

	getProductQuery := `SELECT * FROM products`
	if err := conn.Select(&products, getProductQuery); err != nil {
		panic(err)
	}
	res, err := json.Marshal(products)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(res))
}