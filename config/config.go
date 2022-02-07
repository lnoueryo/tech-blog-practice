package config

import (
	"html/template"
	"log"
	"path/filepath"
	_ "github.com/mattn/go-sqlite3"
	"github.com/joho/godotenv"
	"os"
)

type AppConfig struct {
	UseCache bool
	TemplateCache map[string]*template.Template
	InfoLog *log.Logger
	ErrorLog *log.Logger
	InProduction bool
	Addr string
	Static string
	Media string
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

	commonSettings()

	appEnv, err := readEnvFile(); if err!= nil {
		// if .env is not in local and production environment
		panic("Not found .env")
	}

	if appEnv == "local" {
		configureLocalSettings()
	} else {
		configureProdSettings()
	}
}

func readEnvFile() (string, error) {

	// local
    err := godotenv.Load(".env.dev"); if err == nil {
		return os.Getenv("APP_ENV"), nil
	}

	// production
	err = godotenv.Load(".env"); if err == nil {
		return os.Getenv("APP_ENV"), nil
	}

	return os.Getenv("APP_ENV"), err
}

func commonSettings() {
	// log
	infolog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorlog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	App.InfoLog = infolog
	App.ErrorLog = errorlog

	// file path
	App.Static = "public"
	App.Media = "upload"

	// APIKey
	ApiKey.GitHubClientId = os.Getenv("GITHUB_CLIENT_ID")
	ApiKey.GitHubSecretId = os.Getenv("GITHUB_SECRET_ID")
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
