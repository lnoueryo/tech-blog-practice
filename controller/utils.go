package controller

import (
	"bytes"
	"helloworld/config"
	"helloworld/models"
	"html/template"
	"log"
	"net/http"
	"github.com/jinzhu/gorm"
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
	Posts	  []models.Post
	CSRFToken string
	// Form      *forms.Form
}

var infolog *log.Logger
var errorlog *log.Logger
var DB *gorm.DB


func init() {
	infolog = config.App.InfoLog
	errorlog = config.App.ErrorLog
	DB = models.DB
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
		errorlog.Print("could not get template")
	}

	buf := new(bytes.Buffer)

	_ = t.Execute(buf, td)

	_, err := buf.WriteTo(w)

	if err != nil {
		errorlog.Print("could not get template")
	}
}

