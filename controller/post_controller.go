package controller

import (
	"fmt"
	"helloworld/models"
	"net/http"
)

type Post struct {}

func (p *Post)Index(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.NotFound(w, r)
		return
	}
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
	RenderTemplate(w, r, "post.html", &TemplateData{
		StringMap: stringMap,
	})
}