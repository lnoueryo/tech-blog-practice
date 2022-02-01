package post

import (
	"fmt"
	"helloworld/config"
	"helloworld/controller"
	"helloworld/models"
	"log"
	"net/http"
)

var infolog *log.Logger

func init() {
	infolog = config.App.InfoLog
}

func Index(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.NotFound(w, r)
		return
	}
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
	controller.RenderTemplate(w, r, "post.html", &controller.TemplateData{
		StringMap: stringMap,
	})
}