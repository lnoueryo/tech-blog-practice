package models

import (
	"errors"
	"net/http"
	"regexp"
	"strings"
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

func NewUser(r *http.Request) (User, error) {
	var user User
	err := r.ParseForm()
	if err != nil {
		return user, err
	}

	name := r.Form.Get("name")
	email := r.Form.Get("email")
	image := r.Form.Get("image")
	password := Encrypt(r.Form.Get("password"))
	user = User{Name: name, Email: email, Image: image, Password: password}
	return user, nil
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

func (u *User) Validate(r *http.Request) error {
	err := u.CheckBlank(r); if err != nil {
		return err
	}

	err = u.ComparePassword(r); if err != nil {
		return err
	}

	err = u.CheckEmailFormat(r); if err != nil {
		return err
	}

	err = u.CheckLength(r); if err != nil {
		return err
	}
	return nil
}

func (u *User)TryToLogin(w http.ResponseWriter, r *http.Request) (error) {

	err := u.CheckLoginFormBlank(r); if err != nil {
		return err
	}

	err = u.CheckEmailFormat(r); if err != nil {
		return err
	}

	// Database
	result := DB.Where("email = ?", r.Form.Get("email")).First(&u)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		message := "your email address has not been registered"
		err := errors.New(message)
		return err
	}

	// Password check
	if u.Password != Encrypt(r.Form.Get("password")) {
		message := "password is wrong"
		err := errors.New(message)
		return err
	}
	return nil
}


func (u *User) CheckBlank(r *http.Request) error {
	name := r.FormValue("name")
	if name == "" {
		message := "name is blank"
		err := errors.New(message)
		return err
	}

	email := r.FormValue("email")
	if email == "" {
		message := "email address is blank"
		err := errors.New(message)
		return err
	}

	password := r.FormValue("password")
	if password == "" {
		message := "password is blank"
		err := errors.New(message)
		return err
	}

	confirmation := r.FormValue("confirmation")
	if confirmation == "" {
		message := "confirmation is blank"
		err := errors.New(message)
		return err
	}
	return nil
}

func (u *User) ComparePassword(r *http.Request) error {
	password := r.FormValue("password")
	confirmation := r.FormValue("confirmation")
	if password != confirmation {
		message := "password and confirmation must be the same"
		err := errors.New(message)
		return err
	}
	return nil
}

func (u *User) CheckEmailFormat(r *http.Request) error {
	email := r.FormValue("email")
	regex := `^[a-zA-Z0-9_.+-]+@([a-zA-Z0-9][a-zA-Z0-9-]*[a-zA-Z0-9]*\.)+[a-zA-Z]{2,}$`

	isEnabled := regexp.MustCompile(regex).Match([]byte(email))
	if !isEnabled {
		message := "invalid email address pattern"
		err := errors.New(message)
		return err
	}
	return nil
}

func (u *User) CheckLength(r *http.Request) error {
	name := r.FormValue("name")
	email := r.FormValue("email")
	password := r.FormValue("password")

	if len(name) > 50 {
		message := "name must be less than 50 characters"
		err := errors.New(message)
		return err
	}
	index := strings.Index(email, "@")
	localPart := email[:index]
	if len(localPart) > 64 {
		message := "invalid email address pattern"
		err := errors.New(message)
		return err
	}

	if 8 > len(password) {
		message := "password must be more than 8 characters"
		err := errors.New(message)
		return err
	}
	return nil
}

func (u *User)CheckLoginFormBlank(r *http.Request) error {
	email := r.FormValue("email")
	if email == "" {
		message := "email address is blank"
		err := errors.New(message)
		return err
	}

	password := r.FormValue("password")
	if password == "" {
		message := "password is blank"
		err := errors.New(message)
		return err
	}
	return nil
}