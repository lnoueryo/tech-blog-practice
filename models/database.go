package models

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
)

type Database struct {
	Name 	 string
	Host 	 string
	User 	 string
	Password string
	Port	 string
	Query	 string
}

var DB *gorm.DB

func ConnectMysql(DBSettings Database) (*gorm.DB, error) {
    // // [ユーザ名]:[パスワード]@tcp([ホスト名]:[ポート番号])/[データベース名]?charset=[文字コード]
    dbconf := createMysqlPath(DBSettings)

    db, err := gorm.Open("mysql", dbconf)
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

func createMysqlPath(DBSettings Database) string {
	path := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?%v", DBSettings.User, DBSettings.Password, DBSettings.Host, DBSettings.Port, DBSettings.Name, DBSettings.Query)
	return path
}