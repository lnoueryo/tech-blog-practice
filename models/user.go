package models

import (
	"errors"
	"fmt"
	"helloworld/modules/crypto"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"gorm.io/gorm"
)

type User struct {
	Id        int `gorm:"AUTO_INCREMENT"json:"id"`
	Name      string `json:"name" sql:"CHARACTER SET utf8 COLLATE utf8_unicode_ci"`
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
	randStr, _ := crypto.MakeRandomStr(20)
	filename := randStr + ".png"
	password := crypto.Encrypt(r.Form.Get("password"))
	user = User{Name: r.Form.Get("name"), Email: r.Form.Get("email"), Image: filename, Password: password}
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

func UserLatest(limit int) ([]User, error) {
	var users []User
	result := DB.Order("created_at desc").Limit(limit).Preload("Posts").Find(&users)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return users, result.Error
	}
	return users, nil
}

func (u *User)Create() error {
	result := DB.Create(u)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return result.Error
	}
	return nil
}

func (u *User)Update() error {
	result := DB.Create(u)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return result.Error
	}
	return nil
}

// func User

func SearchUserLike(r *http.Request, column string) ([]User, int64, error) {
	var users []User
	query := r.URL.Query()
	page, _ := strconv.Atoi(query["page"][0])
	if query[column][0] == "" {
		users, count, _ := ChunkUser(page)
		return users, count, nil
	}
	split := 10
	offset := (page - 1) * split
	textSlice := strings.Split(query[column][0], " ")
	var tx *gorm.DB
	var count int64
	for _, text := range textSlice {
		likeText := "%" + text + "%"
		tx = DB.Model(&User{}).Limit(split).Offset(offset).Preload("Posts").Where(fmt.Sprintf("%s LIKE ? ", column), likeText)
	}
	result := tx.Find(&users).Limit(-1).Offset(-1).Count(&count)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return users, 0, result.Error
	}
	return users, count, nil
}

func ChunkUser(page int) ([]User, int64, error) {
	var users []User
	var count int64
	split := 10
	offset := (page -1) * split
	result := DB.Limit(split).Offset(offset).Preload("Posts").Find(&users).Limit(-1).Offset(-1).Count(&count)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return users, 0, result.Error
	}
	return users, count, nil
}

func (u *User)Validate(r *http.Request) error {
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

func (u *User) UpdateValidate(r *http.Request) error {
	err := u.CheckBlankForUpdate(r); if err != nil {
		return err
	}

	err = u.MatchPassword(r); if err != nil {
		return err
	}

	err = u.SearchSameEmail(r); if err != nil {
		return err
	}

	err = u.CheckImage(r); if err != nil {
		return err
	}

	err = u.CheckLengthForUpdate(r); if err != nil {
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
	if u.Password != crypto.Encrypt(r.Form.Get("password")) {
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

func (u *User) CheckBlankForUpdate(r *http.Request) error {
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

	password := r.FormValue("current-password")
	if password == "" {
		message := "password is blank"
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

func (u *User) CheckLengthForUpdate(r *http.Request) error {
	name := r.FormValue("name")
	email := r.FormValue("email")

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

func (u *User)MatchPassword(r *http.Request) error {
	currentPassword := crypto.Encrypt(r.FormValue("current-password"))
	if u.Password != currentPassword {
		message := "current password is wrong"
		err := errors.New(message)
		return err
	}
	return nil
}

func (u *User)CheckImage(r *http.Request) error {
	currentPassword := crypto.Encrypt(r.FormValue("current-password"))
	if u.Password != currentPassword {
		message := "current password is wrong"
		err := errors.New(message)
		return err
	}
	return nil
}

func (u *User)SearchSameEmail(r *http.Request) error {
	var user User
	result := DB.Where("email = ?", r.Form.Get("email")).First(&user)
	if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		if u.Email == user.Email {
			return nil
		}
		err := errors.New("email address is already registered")
		return err
	}
	return nil
}