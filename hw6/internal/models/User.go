package models

type (
	User struct {
		ID            int          `gorm:"AUTO_INCREMENT;not null;PRIMARY_KEY" json:"id"  db:"id"`
		Name          string       `json:"name" db:"name"`
		Email  		  string   	   `json:"email" db:"email"`
		Address 	  string       `json:"address" db:"address"`
		ShoppingCart  []Product    `json:"cart"`
		MobilePhone   string       `json:"phone" db:"mobilePhone"`
	}

	UserFilter struct {
		Query *string `json:"query"`
	}
)