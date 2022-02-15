package models

import (
	"errors"
	"net/http"
	"path/filepath"
	"time"
	"gorm.io/gorm"
)

type Post struct {
	Id        int `gorm:"AUTO_INCREMENT"json:"id"`
	Title     string `json:"title"`
	Content   string `json:"content"`
	Language  	  string `json:"type"`
	Image     string `json:"image"`
	UserID    int   `json:"user_id"`
	User User `gorm:"references:Id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func NewPost(r *http.Request) (Post, error) {
	var post Post
	err := r.ParseForm()
	if err != nil {
		return post, err
	}
    _, fileHeader, err := r.FormFile("image"); if (err != nil) {
		message := "image is required"
		err := errors.New(message)
		return post, err
    }
	randStr, _ := MakeRandomStr(20)
	newFileName := randStr + filepath.Ext(fileHeader.Filename)
	s := GetSession(r)
	userId := s.UserId
	post = Post{Title: r.Form.Get("title"), Content: r.Form.Get("content"), Language: r.Form.Get("language"), UserID: userId, Image: newFileName}
	return post, nil
}

func PostAll() ([]Post, error) {
	var posts []Post
	result := DB.Preload("User").Find(&posts)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return posts, result.Error
	}
	return posts, nil
}

func (p *Post) Create() error {
	result := DB.Create(p)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return result.Error
	}
	return nil
}

func (p *Post) Validate(r *http.Request) error {

	err := p.CheckBlank(r); if err != nil {
		return err
	}

	err = p.LimitSize(r); if err != nil {
		return err
	}

	return nil
}

func (p *Post) CheckBlank(r *http.Request) error {
	title := r.FormValue("title")
	if title == "" {
		message := "title is blank"
		err := errors.New(message)
		return err
	}

	content := r.FormValue("content")
	if content == "" {
		message := "content is blank"
		err := errors.New(message)
		return err
	}

	return nil
}

func (p *Post) LimitSize(r *http.Request) error {
    _, handler, _ := r.FormFile("image"); if handler.Size > 2000000 {
		message := "maximum file size is 2MB"
		err := errors.New(message)
		return err
    }

	return nil
}