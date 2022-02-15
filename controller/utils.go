package controller

import (
	"bytes"
	"errors"
	"helloworld/config"
	"helloworld/models"
	"html/template"
	"log"
	"net/http"
	"io"
	"os"
	"gorm.io/gorm"
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
	Session	  models.Session
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

func StoreImage(r *http.Request, dir string, name string) error {

    // POSTされたファイルデータを取得する
    file, _, err := r.FormFile("image"); if(err != nil) {
        errorlog.Print(err);
        errorlog.Print("couldn't upload the file");
		message := "name is blank"
		err := errors.New(message)
		return err
    }
    // // サーバー側に保存するために空ファイルを作成
    saveImage, err := os.Create("./upload/" + dir + "/" + name); if (err != nil) {
		message := "could't upload the file"
		err := errors.New(message)
		return err
    }
    defer saveImage.Close();
    defer file.Close();
    _, err = io.Copy(saveImage, file); if (err != nil) {
        errorlog.Print(err);
        errorlog.Print("failed to write the file");
		message := "could't upload the file"
		err := errors.New(message)
		return err
    }
	return nil
}


