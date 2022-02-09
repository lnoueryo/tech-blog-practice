package controller

import (
	"errors"
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
	stringMap := make(map[string]string)
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
				Session: models.DeliverSession(r),
			})
			return
		}
		 //usersのあとにキーがある場合
		u.Show(w, r, path)
		return
	}
	userIdStr, path := path[:i], path[i:]
	if path == "/" {
		u.Show(w, r, userIdStr)
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
		Session: models.DeliverSession(r),
	})
}

func (u *Users)Show(w http.ResponseWriter, r *http.Request, id string) {
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
	RenderTemplate(w, r, "user.html", &TemplateData{
		StringMap: stringMap,
		Users: users,
		Session: models.DeliverSession(r),
	})
}