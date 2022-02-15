package controller

import (
	"errors"
	"helloworld/models"
	"net/http"
	"os"
	"strconv"
	"strings"
	"gorm.io/gorm"
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

func (u *Users)Delete(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.NotFound(w, r)
		return
	}
	path := strings.TrimPrefix(r.URL.Path, "/users/delete/") // URLを切り取ってなんとかする
	i := strings.Index(path, "/")
	if i == -1 {
		var user models.User
		if path != "" { //usersのみ
			result := DB.Preload("Posts.User").First(&user, path)
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				errorlog.Print(result)
			}
			s := models.GetSession(r)
			if user.Id == s.UserId {
				result = DB.Select("Posts").Delete(&user)
				if DB.Error != nil {
					profilePage := "/users/" + strconv.Itoa(s.UserId)
					http.Redirect(w, r, profilePage, http.StatusFound)
					return
				} else if result.RowsAffected < 1 {
					profilePage := "/users/" + strconv.Itoa(s.UserId)
					http.Redirect(w, r, profilePage, http.StatusFound)
					return
				}
				os.Remove("./upload/user/" + user.Image)
				for _, p := range user.Posts {
					os.Remove("./upload/post/" + p.Image)
				}

				s := models.GetSession(r)

				err := s.DeleteSession(w, r); if err !=nil {
					errorlog.Print(err)
					http.Redirect(w, r, "/login", http.StatusFound)
					return
				}
				http.Redirect(w, r, "/login", http.StatusFound)
				return
			}
		}
	}
	http.NotFound(w, r)
}
