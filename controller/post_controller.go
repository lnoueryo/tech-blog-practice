package controller

import (
	"helloworld/models"
	"net/http"
)

type Post struct {}

func (p *Post)Index(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.NotFound(w, r)
		return
	}

	stringMap := make(map[string]string)
	RenderTemplate(w, r, "post.html", &TemplateData{
		StringMap: stringMap,
		CSRFToken: models.GenerateCSRFToken(r),
	})
}