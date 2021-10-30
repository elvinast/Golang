package models

type Product struct {
	ID           int     `gorm:"AUTO_INCREMENT;not null;PRIMARY_KEY" json:"id"`
	Title        string  `json:"title"`
	Description  string  `json:"description"`
	Price 		 float64 `json:"price"`
	IsAvailable  bool    `json:"isAvailable"`
	Weight       float64 `json:"weight"`
}
