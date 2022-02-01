package auth

import (
	"errors"
	"fmt"
	"helloworld/config"
	"helloworld/controller"
	"helloworld/models"
	"log"
	"net/http"
	"regexp"

	"github.com/jinzhu/gorm"
)

var infolog *log.Logger
var errorlog *log.Logger
var DB *gorm.DB

func init() {
	infolog = config.App.InfoLog
	errorlog = config.App.ErrorLog
	DB = models.DB
}

func Login(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		session, err := models.CheckSession(r)
		if err == nil {
			infolog.Print(fmt.Sprintf("%v\t%v\t%v\t%v", r.URL, session.Name, session.Email, r.RemoteAddr))
			http.Redirect(w, r, "/", 302)
			return
		}
		fmt.Sprintf("%v\t%v", r.URL, r.RemoteAddr)
		stringMap := make(map[string]string)
		stringMap["email"] = ""
		stringMap["message"] = ""
		stringMap["github"] = "https://github.com/login/oauth/authorize?client_id=cfd4c11c88620861e0ad&redirect_uri=" + config.App.Host + "/oauth/callback"
		controller.RenderTemplate(w, r, "login.html", &controller.TemplateData{
			StringMap: stringMap,
		})
		return
	}
	if r.Method == "POST" {
		err := r.ParseForm()
		if err != nil {
			errorlog.Print(err, "Cannot find user")
		}

		email := r.FormValue("email")
		if email == "" {
			message := "email address is blank"
			errorlog.Print(message)
			validateLogin(w, r, message)
			return
		}
		password := r.FormValue("password")
		if password == "" {
			message := "password is blank"
			errorlog.Print(message)
			validateLogin(w, r, message)
			return
		}
		regex := `^[a-zA-Z0-9_.+-]+@([a-zA-Z0-9][a-zA-Z0-9-]*[a-zA-Z0-9]*\.)+[a-zA-Z]{2,}$`
		emailValidation := regexp.MustCompile(regex).Match([]byte(email))
		if !emailValidation {
			message := "invalid email address pattern"
			errorlog.Print(message)
			validateLogin(w, r, message)
			return
		}

		var user models.User
		// Database
		result := DB.Where("email = ?", r.Form.Get("email")).First(&user)
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			message := "your email address has not been registered"
			errorlog.Print(message)
			validateLogin(w, r, message)
			return
		}

		if user.Password != models.Encrypt(r.Form.Get("password")) {
			message := "password is wrong"
			errorlog.Print(message)
			validateLogin(w, r, message)
			return
		}

		cookieSession(w, r, user)
		http.Redirect(w, r, "/", 302)
		return
	}
	http.NotFound(w, r)
	return
}



func Logout(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.NotFound(w, r)
		return
	}
	ok, err := models.Auth(w, r)
	if !ok {
		errorlog.Print(err)
		http.Redirect(w, r, "/", 302)
		return
	}
	models.DeleteSession(w, r)
	http.Redirect(w, r, "/login", 302)
}

func Register(w http.ResponseWriter, r *http.Request) {
	session, err := models.CheckSession(r)
	if err == nil {
		infolog.Print(fmt.Sprintf("%v\t%v\t%v\t%v", r.URL, session.Name, session.Email, r.RemoteAddr))
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
	controller.RenderTemplate(w, r, "sign-up.html", &controller.TemplateData{
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
	controller.RenderTemplate(w, r, "login.html", &controller.TemplateData{
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
	controller.RenderTemplate(w, r, "sign-up.html", &controller.TemplateData{
		StringMap: stringMap,
	})
}

func GitHubLogin(w http.ResponseWriter, r *http.Request) {
	userInfo, err := GithubOAuth(w, r)
	if err != nil {
		errorlog.Print(err)
	}
	http.Redirect(w, r, "/home?access_token=" + userInfo.AccessToken, 302)
}

func cookieSession(w http.ResponseWriter, r *http.Request, user models.User) {
	infolog.Print(fmt.Sprintf("%v\t%v\t%v\t%v", r.URL, user.Name, user.Email, r.RemoteAddr))
	sessionId, err := models.CreateSession(user)
	if err != nil {
		// danger(err, "Cannot create session")
		errorlog.Print("Cannot create session")
	}
	// Path is needed for all path
	cookie := http.Cookie{
		Name:     "_cookie",
		Value:    sessionId,
		HttpOnly: true,
		Path:     "/",
	}
	http.SetCookie(w, &cookie)
}

