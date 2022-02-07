package controller

import (
	"errors"
	"fmt"
	"helloworld/config"
	"helloworld/models"
	"net/http"
	"regexp"
	"helloworld/modules/oauth"
	"github.com/jinzhu/gorm"
)

type Auth struct {}


func (au *Auth)LoginIndex(w http.ResponseWriter, r *http.Request) {
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
		user, ok := tryToLogin(w, r); if !ok {
			return
		}
		cookieSession(w, r, user)
		http.Redirect(w, r, "/", 302)
		return
	}
	http.NotFound(w, r)
}

func tryToLogin(w http.ResponseWriter, r *http.Request) (userInfo models.User, status bool) {
	var user models.User
	err := r.ParseForm(); if err != nil {
		errorlog.Print(err, "Cannot find user")
	}

	email := r.FormValue("email")
	if email == "" {
		message := "email address is blank"
		errorlog.Print(message)
		validateLogin(w, r, message)
		return user, false
	}
	password := r.FormValue("password")
	if password == "" {
		message := "password is blank"
		errorlog.Print(message)
		validateLogin(w, r, message)
		return user, false
	}
	regex := `^[a-zA-Z0-9_.+-]+@([a-zA-Z0-9][a-zA-Z0-9-]*[a-zA-Z0-9]*\.)+[a-zA-Z]{2,}$`
	emailValidation := regexp.MustCompile(regex).Match([]byte(email))
	if !emailValidation {
		message := "invalid email address pattern"
		errorlog.Print(message)
		validateLogin(w, r, message)
		return user, false
	}

	// Database
	result := DB.Where("email = ?", r.Form.Get("email")).First(&user)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		message := "your email address has not been registered"
		errorlog.Print(message)
		validateLogin(w, r, message)
		return user, false
	}

	if user.Password != models.Encrypt(r.Form.Get("password")) {
		message := "password is wrong"
		errorlog.Print(message)
		validateLogin(w, r, message)
		return user, false
	}
	return user, true
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
	s, err := models.CheckSession(r)
	if err == nil {
		infolog.Print(fmt.Sprintf("%v\t%v\t%v\t%v\t%v", r.Method, r.URL, s.Name, s.Email, r.RemoteAddr))
		http.Redirect(w, r, "/", 302)
		return
	}
    // Is form submitted?
    if r.Method == "POST" {
		err := r.ParseForm()
		if err != nil {
			errorlog.Print(err, "Cannot find user")
		}

		name := r.FormValue("name")
		if name == "" {
			message := "name is blank"
			errorlog.Print(message)
			validateRegistration(w, r, message)
			return
		}

		email := r.FormValue("email")
		if email == "" {
			message := "email address is blank"
			errorlog.Print(message)
			validateRegistration(w, r, message)
			return
		}

		password := r.FormValue("password")
		if password == "" {
			message := "password is blank"
			errorlog.Print(message)
			validateRegistration(w, r, message)
			return
		}

		confirmation := r.FormValue("confirmation")
		if password == "" {
			message := "confirmation is blank"
			errorlog.Print(message)
			validateRegistration(w, r, message)
			return
		}

		if password != confirmation {
			message := "password and confirmation must be the same"
			validateRegistration(w, r, message)
			return
		}

		regex := `^[a-zA-Z0-9_.+-]+@([a-zA-Z0-9][a-zA-Z0-9-]*[a-zA-Z0-9]*\.)+[a-zA-Z]{2,}$`
		emailValidation := regexp.MustCompile(regex).Match([]byte(email))
		if !emailValidation {
			message := "invalid email address pattern"
			errorlog.Print(message)
			validateRegistration(w, r, message)
			return
		}

		if len(name) > 50 {
			message := "name must be less than 50 characters"
			errorlog.Print(message)
			validateRegistration(w, r, message)
			return
		}

		if len(password) < 8 {
			message := "password must be more than 8 characters"
			errorlog.Print(message)
			validateRegistration(w, r, message)
			return
		}

		var user models.User
		// Database
		result  := DB.Where("email = ?", r.Form.Get("email")).First(&user)
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			user = models.NewUser(r)
			err = user.Create()
			if err != nil {
				errorlog.Print(err)
			}
			cookieSession(w, r, user)
			http.Redirect(w, r, "/", 302)
		}
		return // AND return!
    }
	stringMap := make(map[string]string)
	stringMap["name"] = ""
	stringMap["email"] = ""
	stringMap["message"] = ""
	RenderTemplate(w, r, "sign-up.html", &TemplateData{
		StringMap: stringMap,
	})
}

func validateLogin(w http.ResponseWriter, r *http.Request, message string) {
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

func validateRegistration(w http.ResponseWriter, r *http.Request, message string) {
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

func (au *Auth)GitHubLogin(w http.ResponseWriter, r *http.Request) {
	userInfo, err := oauth.GithubOAuth(w, r)
	if err != nil {
		errorlog.Print(err)
	}
	// databaseの処理Createを記載する↓↓
	http.Redirect(w, r, "/?access_token=" + userInfo.AccessToken, 302)
}

func cookieSession(w http.ResponseWriter, r *http.Request, user models.User) {
	infolog.Print(fmt.Sprintf("%v\t%v\t%v\t%v", r.URL, user.Name, user.Email, r.RemoteAddr))
	sessionId, err := models.CreateSession(user)
	if err != nil {
		errorlog.Print(err)
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
}

