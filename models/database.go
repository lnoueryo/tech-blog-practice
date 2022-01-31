package models

import (
	"os"
	"github.com/jinzhu/gorm"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/mattn/go-sqlite3"
)

var DB *gorm.DB

func ConnectMysql() (*gorm.DB, error) {
    // // [ユーザ名]:[パスワード]@tcp([ホスト名]:[ポート番号])/[データベース名]?charset=[文字コード]
    dbconf := os.Getenv("DB_CONF")

    db, err := gorm.Open("mysql", dbconf + "&charset=utf8&loc=Local")
	DB = db
    if err != nil {
        panic(err.Error())
    }
    return DB, err
}

func ConnectSqlite3() (*gorm.DB, error) {

	db, err := gorm.Open("sqlite3", "gorm.db")
	if err != nil {
		panic("failed to connect database")
	}
	DB = db

	return DB, err
}