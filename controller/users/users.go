package users

import (
	"errors"
	"fmt"
	"helloworld/config"
	"helloworld/controller"
	"helloworld/models"
	"log"
	"net/http"
	"strconv"
	"strings"
	"github.com/jinzhu/gorm"
)

var infolog *log.Logger
var errorlog *log.Logger
var DB *gorm.DB

func init() {
	infolog = config.App.InfoLog
	DB = models.DB
}

func Index(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/users/") // URLを切り取ってなんとかする
	i := strings.Index(path, "/")
	// session
	session, err := models.CheckSession(r)
	if err != nil {
		fmt.Sprintf("%v\t%v", r.URL, r.RemoteAddr)
		http.Redirect(w, r, "/login", 302)
		return
	}
	infolog.Print(fmt.Sprintf("%v\t%v\t%v\t%v", r.URL, session.Name, session.Email, r.RemoteAddr))
	stringMap := make(map[string]string)
	stringMap["csrf_token"] = session.CSRFToken
	if i == -1 {
		var users []models.User
		if path == "" { //usersのみ
			result := DB.Preload("Posts").Find(&users)
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				errorlog.Print(result)
			}
			controller.RenderTemplate(w, r, "users.html", &controller.TemplateData{
				StringMap: stringMap,
				Users: users,
			})
			return
		}
		 //usersのあとにキーがある場合
		Show(w, r, path, session)
		return
	}
	userIdStr, path := path[:i], path[i:]
	infolog.Print(userIdStr)
	infolog.Print(path)
	if path == "/" {
		Show(w, r, userIdStr, session)
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
	controller.RenderTemplate(w, r, "diary.html", &controller.TemplateData{
		StringMap: stringMap,
		Users: users,
		Posts: posts,
	})
}

func Show(w http.ResponseWriter, r *http.Request, id string, session models.Session) {
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
	controller.RenderTemplate(w, r, "user.html", &controller.TemplateData{
		StringMap: stringMap,
		Users: users,
	})
}