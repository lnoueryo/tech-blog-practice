package controller

import (
	"errors"
	"fmt"
	"helloworld/config"
	"helloworld/models"
	"net/http"
	"helloworld/modules/oauth"
	"github.com/jinzhu/gorm"
)

type Auth struct {}


func (au *Auth)Login(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		s, err := models.CheckSession(r); if err == nil {
			infolog.Print(fmt.Sprintf("%v\t%v\t%v\t%v\t%v", r.Method, r.URL, s.Name, s.Email, r.RemoteAddr))
			http.Redirect(w, r, "/", 302)
			return
		}
		infolog.Print(fmt.Sprintf("%v\t%v\t%v", r.Method, r.URL, r.RemoteAddr))
		stringMap := make(map[string]string)
		stringMap["email"] = ""
		stringMap["message"] = ""
		stringMap["github"] = "https://github.com/login/oauth/authorize?client_id=cfd4c11c88620861e0ad&redirect_uri=" + config.App.Host + "/oauth/callback"
		RenderTemplate(w, r, "login.html", &TemplateData{
			StringMap: stringMap,
		})
		return
	}
	if r.Method == "POST" {
		infolog.Print(fmt.Sprintf("%v\t%v\t%v", r.Method, r.URL, r.RemoteAddr))
		u, err := models.NewUser(r); if err != nil {
			errorlog.Print(err)
		}
		err = u.TryToLogin(w, r); if err != nil {
			errorlog.Print(err)
			redirectLogin(w, r, err.Error())
			return
		}
		err = setCookieSession(w, r, u)
		http.Redirect(w, r, "/", 302)
		return
	}
	http.NotFound(w, r)
}

func (au *Auth)Logout(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		infolog.Print(fmt.Sprintf("%v\t%v", r.URL, r.RemoteAddr))
		http.NotFound(w, r)
		return
	}

	s, err := models.CheckSession(r); if err != nil {
		errorlog.Print(err)
		err = models.DeleteCookie(w, r); if err != nil {
			errorlog.Print(err)
		}
		http.Redirect(w, r, "/login", 302)
		return
	}
	infolog.Print(fmt.Sprintf("%v\t%v\t%v\t%v\t%v", r.Method, r.URL, s.Name, s.Email, r.RemoteAddr))

	if !s.CheckCSRFToken(r) {
		err = errors.New("invalid csrf_token")
		errorlog.Print(err)
		http.Redirect(w, r, "/", 302)
		return
	}

	err = s.DeleteSession(w, r); if err !=nil {
		errorlog.Print(err)
		http.Redirect(w, r, "/login", 302)
		return
	}
	http.Redirect(w, r, "/login", 302)
}

func (au *Auth)Register(w http.ResponseWriter, r *http.Request) {

	if r.Method == "GET" {
		s, err := models.CheckSession(r)
		if err == nil {
			infolog.Print(fmt.Sprintf("%v\t%v\t%v\t%v\t%v", r.Method, r.URL, s.Name, s.Email, r.RemoteAddr))
			http.Redirect(w, r, "/", 302)
			return
		}
		infolog.Print(fmt.Sprintf("%v\t%v\t%v", r.Method, r.URL, r.RemoteAddr))
		stringMap := make(map[string]string)
		stringMap["name"] = ""
		stringMap["email"] = ""
		stringMap["message"] = ""
		RenderTemplate(w, r, "sign-up.html", &TemplateData{
			StringMap: stringMap,
		})
		return
	}

    if r.Method == "POST" {
		u, err := models.NewUser(r); if err != nil {
			errorlog.Print(err)
		}

		err = u.Validate(r); if err != nil {
			errorlog.Print(err)
			redirectRegister(w, r, err.Error())
			return
		}

		var user models.User
		// Database
		result := DB.Where("email = ?", r.Form.Get("email")).First(&user)
		if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
			err = errors.New("email address is already registered")
			redirectRegister(w, r, err.Error())
			return
		}

		err = u.Create()
		if err != nil {
			err = errors.New("couldn't register your account")
			errorlog.Print(err)
			redirectRegister(w, r, err.Error())
			return
		}
		err = setCookieSession(w, r, user)
		http.Redirect(w, r, "/", 302)
		return
    }
	http.NotFound(w, r)
}

func (au *Auth)GitHubLogin(w http.ResponseWriter, r *http.Request) {
	userInfo, err := oauth.GithubOAuth(w, r)
	if err != nil {
		errorlog.Print(err)
	}
	// databaseの処理Createを記載する↓↓
	http.Redirect(w, r, "/?access_token=" + userInfo.AccessToken, 302)
}

func redirectLogin(w http.ResponseWriter, r *http.Request, message string) {
    // Params for rendering the page
    stringMap := make(map[string]string)
	email := r.FormValue("email")
	stringMap["email"] = email
	stringMap["message"] = message
	stringMap["github"] = "https://github.com/login/oauth/authorize?client_id=cfd4c11c88620861e0ad&redirect_uri=" + config.App.Host + "/oauth/callback"
	RenderTemplate(w, r, "login.html", &TemplateData{
		StringMap: stringMap,
	})
}

func redirectRegister(w http.ResponseWriter, r *http.Request, message string) {
    // Params for rendering the page
    stringMap := make(map[string]string)
	name := r.FormValue("name")
	email := r.FormValue("email")
	stringMap["name"] = name
	stringMap["email"] = email
	stringMap["message"] = message
	RenderTemplate(w, r, "sign-up.html", &TemplateData{
		StringMap: stringMap,
	})
}

func setCookieSession(w http.ResponseWriter, r *http.Request, u models.User) (error) {
	sessionId, err := models.CreateSession(u)
	if err != nil {
		return err
	}
	// Path is needed for all path
	cookie := http.Cookie{
		Name:     "_cookie",
		Value:    sessionId,
		HttpOnly: true,
		Secure:   true,
		Path:     "/",
	}
	http.SetCookie(w, &cookie)
	return nil
}