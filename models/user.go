package models

import (
	"errors"
	"net/http"
	"time"
	"github.com/jinzhu/gorm"
)

type User struct {
	Id        int `gorm:"AUTO_INCREMENT"json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	Image     string `json:"image"`
	Posts 	  []Post `json:"posts"`
	CreatedAt time.Time `json:"created_at"`
}

func NewUser(r *http.Request) User {
	name := r.Form.Get("name")
	email := r.Form.Get("email")
	image := r.Form.Get("image")
	password := Encrypt(r.Form.Get("password"))
	user := User{Name: name, Email: email, Image: image, Password: password}
	return user
}

func UserAll() ([]User, error) {
	var users []User
	result := DB.Preload("Posts").Find(&users)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return users, result.Error
	}
	return users, nil
}

func (u *User) Create() error {
	result := DB.Create(u)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return result.Error
	}
	return nil
}