package models

type FoodItem struct {
	ID           int     `json:"id"`
	Title        string  `json:"title"`
	Description  string  `json:"description"`
	Price 		 float64 `json:"price"`
	isAvailable  bool    `json:"isAvailable"`
}
