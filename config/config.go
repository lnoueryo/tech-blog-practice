package config

import (
	"html/template"
	"log"
	"path/filepath"
	_ "github.com/mattn/go-sqlite3"
	"github.com/joho/godotenv"
	"os"
	"helloworld/models"
)

type AppConfig struct {
	UseCache bool
	TemplateCache map[string]*template.Template
	InfoLog *log.Logger
	ErrorLog *log.Logger
	InProduction bool
	Addr string
	Static string
	Host string
}
type APIKey struct {
	GitHubClientId string
	GitHubSecretId string
}
var App AppConfig
var ApiKey APIKey
var infolog *log.Logger
var errorlog *log.Logger
func init() {
	infolog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorlog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	appEnv := readEnv()
	if appEnv == "local" {
		configureLocalSettings()
		// configureProdSettings()
	} else {
		configureProdSettings()
	}
	App.Addr = ":8080"
	App.Static = "public"
	App.InfoLog = infolog
	App.ErrorLog = errorlog
	ApiKey.GitHubClientId = os.Getenv("GITHUB_CLIENT_ID")
	ApiKey.GitHubSecretId = os.Getenv("GITHUB_SECRET_ID")
}

func readEnv() string {
    err := godotenv.Load(".env.dev")
    if err != nil {
		err = godotenv.Load(".env")
		if err != nil {
			// .env読めなかった場合の処理
			panic("Not found .env")
		}
    }
	return os.Getenv("APP_ENV")
}

func configureLocalSettings() {
	App.UseCache = false
	models.ConnectMysql()
	App.Host = os.Getenv("APP_HOST")
}

func configureProdSettings() {
	App.UseCache = true
	tc, err := CreateTemplateCache()
	if err != nil {
		errorlog.Fatal(err)
	}
	App.TemplateCache = tc
	models.ConnectSqlite3()
	App.Host = os.Getenv("APP_HOST")
}

func CreateTemplateCache() (map[string]*template.Template, error) {
	myCache := map[string]*template.Template{}

	pages, err := filepath.Glob("./templates/pages/*.html")
	if err != nil {
		return myCache, err
	}
	for _, page := range pages {
		name := filepath.Base(page)
		ts := template.Must(template.New(name).ParseFiles(page))
		matches, err := filepath.Glob("./templates/layouts/app.html")
		if err != nil {
			return myCache, err
		}

		if len(matches) > 0 {
			ts, err = ts.ParseGlob("./templates/layouts/app.html")
			if err != nil {
				return myCache, err
			}
		}

		myCache[name] = ts
	}
	return myCache, nil

}
