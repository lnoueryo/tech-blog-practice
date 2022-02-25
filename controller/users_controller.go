package controller

import (
	"errors"
	"helloworld/models"
	"helloworld/modules/crypto"
	"math"
	"net/http"
	"os"
	"path/filepath"
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
		if path == "" { //usersのみ
			var users []models.User
			var split int64 = 10
			searchColumn := "name"
			query := r.URL.Query()
			if len(query["page"]) == 0 || len(query["name"]) == 0 {
				http.Redirect(w, r, "/users/?name=&page=1", http.StatusFound)
				return
			}
			users, count, _ := models.SearchUserLike(r, searchColumn)
			pageNum, _ := strconv.Atoi(query["page"][0])
			keyword := query["name"][0]
			next := pageNum + 1
			previous := pageNum - 1
			pagination := math.Floor(float64(count / split))
			var s []int
			for i := 0; i <= int(pagination); i++ {
				s = append(s, i + 1)
			}
			if len(s) == 0 {
				s = append(s, 1)
			}
			infolog.Print(s)
			page := Page{keyword, pageNum, next, previous, int(pagination), s}
			RenderTemplate(w, r, "users.html", &TemplateData{
				StringMap: stringMap,
				Users: users,
				Session: models.DeliverSession(r),
				Page: page,
			})
			return
		}
		 // if id exists after path "users"
		id, err := strconv.Atoi(path); if err != nil {
			http.NotFound(w, r)
			return
		}
		u.Show(w, r, id)
		return
	}
	// edit page
	if path == "edit/" {
		u.Edit(w, r)
		return
	}
	userIdStr, path := path[:i], path[i:]
	if path == "/" {
		id, err := strconv.Atoi(userIdStr); if err != nil {
			http.NotFound(w, r)
			return
		}
		u.Show(w, r, id)
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

func (u *Users)Show(w http.ResponseWriter, r *http.Request, id int) {
	var users []models.User
	result := DB.Preload("Posts").First(&users, id)
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

func (u *Users)Edit(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		stringMap := make(map[string]string)
		stringMap["message"] = ""
		RenderTemplate(w, r, "user-edit.html", &TemplateData{
			StringMap: stringMap,
			Session: models.DeliverSession(r),
		})
		return
	}
	if r.Method == "POST" {
		u.Update(w, r)
		return
	}

	http.NotFound(w, r)
	return
}

func (us *Users)Update(w http.ResponseWriter, r *http.Request) {
	s := models.GetSession(r)

	var u models.User
	result := DB.First(&u, s.UserId)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		err := errors.New("couldn't update")
		redirectUserEdit(w, r, err.Error())
		return
	}

	err := u.UpdateValidate(r); if err != nil {
		errorlog.Print(err)
		redirectUserEdit(w, r, err.Error())
		return
	}

    _, fileHeader, err := r.FormFile("image"); if err == nil {
		randStr, _ := crypto.MakeRandomStr(20)
		u.Image = randStr + filepath.Ext(fileHeader.Filename)
		dirName := "user"
		err = StoreImage(r, dirName, u.Image); if err != nil {
			errorlog.Print(err)
			redirectUserEdit(w, r, err.Error())
			return
		}
		// ファイルサイズ変更処理　↓↓
		// if fileHeader.Size > 2000000 {
		// 	image.ResizeImage(u.Image)
		// }
		imagePath := "./upload/user/" + s.Image
		os.Remove(imagePath)
    }
	u.Name = r.FormValue("name")
	u.Email = r.FormValue("email")
	DB.Save(&u)
	// セッション変更
	err = s.DeleteSession(w, r); if err !=nil {
		errorlog.Print(err)
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}
	err = setCookieSession(w, r, u); if err != nil {
		http.Redirect(w, r, "/login", http.StatusFound)
	}
	profilePage := "/users/" + strconv.Itoa(u.Id)
	http.Redirect(w, r, profilePage, http.StatusFound)
	return
	// dirName := "user"
	// err = StoreImage(r, dirName, u.Image); if err != nil {
	// 	errorlog.Print(err)
	// 	redirectUserEdit(w, r, err.Error())
	// 	return
	// }
	// infolog.Print(p)
	// err = u.Create()
	// if err != nil {
	// 	os.Remove("/upload/user/" + u.Image)
	// 	err = errors.New("couldn't register your account")
	// 	errorlog.Print(err)
	// 	redirectRegister(w, r, err.Error())
	// 	return
	// }
	// err = setCookieSession(w, r, u); if err != nil {
	// 	errorlog.Print(err)
	// }
	http.Redirect(w, r, "/", http.StatusFound)
	return
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


func redirectUserEdit(w http.ResponseWriter, r *http.Request, message string) {
    stringMap := make(map[string]string)
	stringMap["message"] = message
	RenderTemplate(w, r, "user-edit.html", &TemplateData{
		StringMap: stringMap,
		Session: models.DeliverSession(r),
	})
}