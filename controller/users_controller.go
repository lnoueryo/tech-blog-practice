package controller

import (
	"errors"
	"fmt"
	"helloworld/models"
	"net/http"
	"strconv"
	"strings"
	"github.com/jinzhu/gorm"
)

type Users struct {}

func (u *Users)Index(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/users/") // URLを切り取ってなんとかする
	i := strings.Index(path, "/")
	// session
	s, err := models.CheckSession(r)
	if err != nil {
		fmt.Sprintf("%v\t%v", r.URL, r.RemoteAddr)
		http.Redirect(w, r, "/login", 302)
		return
	}
	infolog.Print(fmt.Sprintf("%v\t%v\t%v\t%v\t%v", r.Method, r.URL, s.Name, s.Email, r.RemoteAddr))
	stringMap := make(map[string]string)
	stringMap["csrf_token"] = s.CSRFToken
	if i == -1 {
		var users []models.User
		if path == "" { //usersのみ
			result := DB.Preload("Posts").Find(&users)
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				errorlog.Print(result)
			}
			RenderTemplate(w, r, "users.html", &TemplateData{
				StringMap: stringMap,
				Users: users,
			})
			return
		}
		 //usersのあとにキーがある場合
		u.Show(w, r, path, s)
		return
	}
	userIdStr, path := path[:i], path[i:]
	infolog.Print(userIdStr)
	infolog.Print(path)
	if path == "/" {
		u.Show(w, r, userIdStr, s)
		return
	}
	userId, _ := strconv.Atoi(userIdStr)
	var users []models.User
	result := DB.Preload("Posts").First(&users, userId)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		errorlog.Print(result)
	}
	if result.RowsAffected == 0 {
		http.NotFound(w, r)
		return
	}
	var posts []models.Post
	path = strings.TrimPrefix(path, "/")
	postId, _ := strconv.Atoi(path)
	result = DB.Preload("User").First(&posts, postId)
	if result.RowsAffected == 0 {
		http.NotFound(w, r)
		return
	}
	RenderTemplate(w, r, "diary.html", &TemplateData{
		StringMap: stringMap,
		Users: users,
		Posts: posts,
	})
}

func (u *Users)Show(w http.ResponseWriter, r *http.Request, id string, session models.Session) {
	var users []models.User
	userId, _ := strconv.Atoi(id)
	result := DB.Preload("Posts").First(&users, userId)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		errorlog.Print(result)
	}
	if result.RowsAffected == 0 {
		http.NotFound(w, r)
		return
	}
	stringMap := make(map[string]string)
	stringMap["csrf_token"] = session.CSRFToken
	RenderTemplate(w, r, "user.html", &TemplateData{
		StringMap: stringMap,
		Users: users,
	})
}