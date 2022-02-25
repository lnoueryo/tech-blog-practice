package models

import (
	"helloworld/config"
	"time"
	"gorm.io/gorm"
)

var DB *gorm.DB

func init() {
	DB = config.DB
}

func timeToString(t time.Time) string {
	str := t.Format("20060102150405")
	return str
}
