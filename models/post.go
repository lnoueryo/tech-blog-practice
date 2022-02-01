package models

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/jinzhu/gorm"
)

type Post struct {
	Id        int `gorm:"AUTO_INCREMENT"json:"id"`
	Title     string `json:"title"`
	Content   string `json:"content"`
	Type  	  string `json:"type"`
	UserID    int   `json:"user_id"`
	User User `gorm:"foreignKey:Id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func NewPost(r *http.Request) Post {
	title := r.Form.Get("title")
	content := r.Form.Get("content")
	tag := r.Form.Get("tag")
	session, err := CheckSession(r)
	if err != nil {
		fmt.Println(err)
	}
	userId := session.UserId
	post := Post{Title: title, Content: content, Type: tag, UserID: userId}
	return post
}

func (post *Post) Create() error {
	result := DB.Create(post)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return result.Error
	}
	return nil
}