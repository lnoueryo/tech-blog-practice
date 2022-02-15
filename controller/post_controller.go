package controller

import (
	"errors"
	"fmt"
	"helloworld/models"
	"net/http"
	"os"
	"strconv"
	"strings"

	"gorm.io/gorm"
)

type Post struct {}

func (p *Post)Index(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		stringMap := make(map[string]string)
		stringMap["title"] = ""
		stringMap["content"] = ""
		stringMap["language"] = ""
		RenderTemplate(w, r, "post.html", &TemplateData{
			StringMap: stringMap,
			Session: models.DeliverSession(r),
		})
		return
	}
	if r.Method == "POST" {
		p, err := models.NewPost(r); if err != nil {
			errorlog.Print(err)
			redirectPost(w, r, err.Error())
			return
		}
		err = p.Validate(r); if err != nil {
			errorlog.Print(err)
			redirectPost(w, r, err.Error())
			return
		}
		dirName := "post"
		err = StoreImage(r, dirName, p.Image); if err != nil {
			errorlog.Print(err)
			redirectPost(w, r, err.Error())
			return
		}
		infolog.Print(p)
		// redirectPost(w, r, "123456789")

		// var post models.Post
		// Database
		// result := DB.Where("email = ?", r.Form.Get("email")).First(&post)
		// if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		// 	err = errors.New("email address is already registered")
		// 	redirectRegister(w, r, err.Error())
		// 	return
		// }

		err = p.Create(); if err != nil {
			err = errors.New("couldn't register your account")
			errorlog.Print(err)
			redirectPost(w, r, err.Error())
			return
		}
		s := models.GetSession(r)
		profilePage := "/users/" + strconv.Itoa(s.UserId)
		http.Redirect(w, r, profilePage, http.StatusFound)
		return
	}

	http.NotFound(w, r)
}
func (p *Post)Delete(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.NotFound(w, r)
		return
	}
	path := strings.TrimPrefix(r.URL.Path, "/post/delete/") // URLを切り取ってなんとかする
	i := strings.Index(path, "/")
	fmt.Println(path, i)
	if i == -1 {
		var post models.Post
		if path != "" { //usersのみ
			result := DB.First(&post, path)
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				errorlog.Print(result)
			}
			s := models.GetSession(r)
			if post.UserID == s.UserId {
				os.Remove("./upload/post/" + post.Image)
				DB.Delete(&post)
				fmt.Println(post.UserID, s.UserId)
				profilePage := "/users/" + strconv.Itoa(s.UserId)
				http.Redirect(w, r, profilePage, http.StatusFound)
				return
			}
		}
	}
	http.NotFound(w, r)
}

func redirectPost(w http.ResponseWriter, r *http.Request, message string) {
    stringMap := make(map[string]string)
	title := r.FormValue("title")
	language := r.FormValue("language")
	content := r.FormValue("content")
	stringMap["title"] = title
	stringMap["language"] = language
	stringMap["content"] = content
	stringMap["message"] = message
	RenderTemplate(w, r, "post.html", &TemplateData{
		StringMap: stringMap,
	})
}
