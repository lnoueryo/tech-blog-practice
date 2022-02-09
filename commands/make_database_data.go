package commands

import (
	"errors"
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
	_ "github.com/go-sql-driver/mysql"
)

type User struct {
	Id        int `gorm:"AUTO_INCREMENT"json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	Image     string `json:"image"`
	Posts 	  []Post `json:"posts"`
	CreatedAt time.Time `json:"created_at"`
}

type Post struct {
	Id        int `gorm:"AUTO_INCREMENT"json:"id"`
	Title     string `json:"title"`
	Content   string `json:"content"`
	Language  	  string `json:"type"`
	Image     string `json:"image"`
	UserID    int   `json:"user_id"`
	User User `gorm:"references:Id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func MakeDBData(name string, arg1 string) error {
	if arg1 != "" {
        err := fmt.Errorf(fmt.Sprintf(`no command "go run main.go %v"`, CreateArgsText()))
		return err
	}

	if name == "sqlite3" {
		sqlite()
		return nil
	}

	if name == "mysql" {
		mysql()
		return nil
	}
	err := fmt.Errorf(fmt.Sprintf(`no command "go run main.go %v"`, CreateArgsText()))
	return err
}

func sqlite() {
	DB, err := gorm.Open("sqlite3", "gorm.db")
	if err != nil {
		panic("failed to connect database")
	}
	defer DB.Close()
	seriesOfCreation(DB)
}

func mysql() {
	dbconf := `root:admin@/tech-blog?parseTime=true`
	DB, err := gorm.Open("mysql", dbconf + "&charset=utf8&loc=Local")
	if err != nil {
		panic("failed to connect database")
	}
	defer DB.Close()
	seriesOfCreation(DB)
}

func seriesOfCreation(DB *gorm.DB) {
	// isUsersTable := db.Migrator().HasTable("a")
	// fmt.Println(isUsersTable)
	migrateUserTable(DB)
	migratePostTable(DB)
	createUser(DB)
	createPost(DB)
	// readUser(DB)
	// readPost(DB)
}

func migrateUserTable(db *gorm.DB) {
	db.AutoMigrate(&User{})
	db.DropTable(&User{})
	db.CreateTable(&User{})
}

func migratePostTable(db *gorm.DB) {
	db.AutoMigrate(&Post{})
	db.DropTable(&Post{})
	db.CreateTable(&Post{})
}

func createUser(db *gorm.DB) {
	var newUser = User{
		Id: 1,
		Name: "RIO",
		Email: "popo62520908@gmail.com",
		Password: "52f96e51831c8229413e28b0e58fa3b992f7571e4ff5bf5ccfc1a21f391e4f05",
		Image: "CWCM67iUYAAZ1Kp.png",
		CreatedAt: time.Now(),
	}
	db.Create(&newUser)
}

func createPost(db *gorm.DB) {
	var post1 = Post{
		Id: 1,
		Title: "I started to keep diaries",
		Content: "helohelo",
		Language: "English",
		Image: "abc.png",
		UserID: 1,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	var post2 = Post{
		Id: 2,
		Title: "How to learn Go",
		Content: "hogehoge",
		Language: "English",
		Image: "abcd.png",
		UserID: 1,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	db.Create(&post1)
	db.Create(&post2)
}

func readUser(db *gorm.DB) {
	var user User
	result := db.Preload("Posts").Where("email = ?", "popo62520908@gmail.com").First(&user)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		
	}
	fmt.Println(user)
}

func readPost(db *gorm.DB) {
	var post Post
	result := db.Preload("User").Find(&post)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		
	}
	fmt.Println(post)
}