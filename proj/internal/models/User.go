package models

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"golang.org/x/crypto/bcrypt"
)

var userNextID = 0

type (
	User struct {
		ID            int          `gorm:"AUTO_INCREMENT;not null;PRIMARY_KEY" json:"id"  db:"id"`
		Name          string       `json:"name" db:"name"`
		Email  		  string   	   `json:"email" db:"email"`
		Password          string `json:"password" db:"-"`
		Address 	  string       `json:"address" db:"address"`
		ShoppingCart  []Product    `json:"cart"`
		MobilePhone   string       `json:"phone" db:"mobilePhone"`
	}

	UsersFilter struct {
		Query *string `json:"query"`
	}
)

func (u *User) NextId() {
	u.ID = userNextID
	userNextID++
}

func (u *User) Validate() error {
	return validation.ValidateStruct(
		u,
		validation.Field(&u.Email, validation.Required, is.Email),
		validation.Field(&u.Password, validation.Required, validation.Length(8, 30)),
	)
}

func (u *User) Sanitize() {
	u.Password = ""
}