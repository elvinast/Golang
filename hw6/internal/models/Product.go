package models

type ( 
	Product struct {
		ID           int     `gorm:"AUTO_INCREMENT;not null;PRIMARY_KEY" json:"id" db:"id"`
		Title        string  `json:"title" db:"title"`
		Description  string  `json:"description" db:"description"`
		Price 		 float64 `json:"price" db:"price"`
		IsAvailable  bool    `json:"isAvailable" db:"isavailable"`
		Weight       float64 `json:"weight" db:"weight"`
	}

	ProductFilter struct {
		Query *string `json:"query"`
	}
)