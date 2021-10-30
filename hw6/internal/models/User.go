package models

type User struct {
	ID            int          `gorm:"AUTO_INCREMENT;not null;PRIMARY_KEY" json:"id"`
	Name          string       `json:"name"`
	Email  		  string   	   `json:"email"`
	Address 	  string       `json:"address"`
	ShoppingCart  []Product    `json:"cart"`
	MobilePhone   string       `json:"phone"`
}