package controller

import (
	"bytes"
	"html/template"
	"log"
	"net/http"
	"helloworld/config"
	"helloworld/models"
)

type TemplateData struct {
	StringMap map[string]string
	IntMap    map[string]string
	FloatMap  map[string]float32
	Data      map[string]interface{}
	Flash     string
	Warning   string
	Error     string
	JSON      []byte
	Users     []models.User
	// Form      *forms.Form
}

func RenderTemplate(w http.ResponseWriter, r *http.Request, tmpl string, td *TemplateData) {
	var tc map[string]*template.Template
	if config.App.UseCache {
		tc = config.App.TemplateCache
	} else {
		tc, _ = config.CreateTemplateCache()
	}
	t, ok := tc[tmpl]
	if !ok {
		log.Fatal("could not get template")
	}

	buf := new(bytes.Buffer)

	//csrftoken
	// td = AddDefaultData(td, r)

	_ = t.Execute(buf, td)

	_, err := buf.WriteTo(w)

	if err != nil {
		log.Fatal(err)
	}
}


