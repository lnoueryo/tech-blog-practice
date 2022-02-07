package config

import (
	"helloworld/models"
	"os"
)

func configureLocalSettings() {
	App.UseCache = false
	App.Host = os.Getenv("APP_HOST")
	App.Addr = "127.0.0.1:8080"

	// DB接続
	DBSet := models.Database{
		Name: os.Getenv("DB_NAME"),
		Host: os.Getenv("DB_HOST"),
		User: os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		Port: os.Getenv("DB_PORT"),
		Query: os.Getenv("DB_QUERY"),
	}
	models.ConnectMysql(DBSet)
}