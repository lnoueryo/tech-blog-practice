package config

import (
	"helloworld/models"
	"os"
)

func configureProdSettings() {
	App.UseCache = true
	tc, err := CreateTemplateCache()
	if err != nil {
		errorlog.Fatal(err)
	}
	App.TemplateCache = tc
	App.Host = os.Getenv("APP_HOST")
	App.Addr = ":8080"
	// DBSet := models.Database{
	// 	Name: os.Getenv("DB_NAME"),
	// 	Host: os.Getenv("DB_HOST"),
	// 	User: os.Getenv("DB_USER"),
	// 	Password: os.Getenv("DB_PASSWORD"),
	// 	Port: os.Getenv("DB_PORT"),
	// 	Query: os.Getenv("DB_QUERY"),
	// }
	models.ConnectSqlite3()
}