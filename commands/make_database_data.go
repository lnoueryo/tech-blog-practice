package commands

import (
	"encoding/json"
	"fmt"
	"helloworld/models"
	"helloworld/modules/image"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
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
	allDatabase()
	// if name == "" {
	// 	fmt.Print("sqlite3 or mysql? ")
	// 	name = AskRequiredThing()
	// }

	// if name == "sqlite3" {
	// 	sqlite3()
	// 	return nil
	// }

	// if name == "mysql" {
	// 	connectMysql()
	// 	return nil
	// }
	err := fmt.Errorf(fmt.Sprintf(`no command "go run main.go %v"`, CreateArgsText()))
	return err
}

// func sqlite3() {
// 	// DB, err := gorm.Open("sqlite3", "gorm.db")
// 	DB, err := gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{})
// 	if err != nil {
// 		panic("failed to connect database")
// 	}
// 	seriesOfCreation(DB)
// }

// func connectMysql() {
// 	dbconf := `root:popo0908@/practices?parseTime=true&charset=utf8&loc=Local`
// 	// dbconf := `root:admin@/tech-blog?parseTime=true&charset=utf8&loc=Local`
// 	// DB, err := gorm.Open("mysql", dbconf)
// 	DB, err := gorm.Open(mysql.Open(dbconf), &gorm.Config{})
// 	if err != nil {
// 		panic("failed to connect database")
// 	}
// 	seriesOfCreation(DB)
// }

func allDatabase() {
	SQDB, err := gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	dbconf := `root:popo0908@/practices?parseTime=true&charset=utf8&loc=Local`
	// dbconf := `root:admin@/tech-blog?parseTime=true&charset=utf8&loc=Local`
	// DB, err := gorm.Open("mysql", dbconf)
	MQDB, err := gorm.Open(mysql.Open(dbconf), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	seriesOfCreation(SQDB, MQDB)
}

func seriesOfCreation(SQDB *gorm.DB, MQDB *gorm.DB) {
	// isUsersTable := db.Migrator().HasTable("a")
	// fmt.Println(isUsersTable)
	migrateUserTable(SQDB, MQDB)
	migratePostTable(SQDB, MQDB)
	createUser(SQDB, MQDB)
	createUsers(SQDB, MQDB)
	createPost(SQDB, MQDB)
	createPosts(SQDB, MQDB)
	readUser(SQDB)
	readPost(SQDB)
}

func migrateUserTable(SQDB *gorm.DB, MQDB *gorm.DB) {
	SQDB.AutoMigrate(&User{})
	SQDB.Migrator().DropTable(&User{})
	SQDB.Migrator().CreateTable(&User{})
	MQDB.AutoMigrate(&User{})
	MQDB.Migrator().DropTable(&User{})
	MQDB.Migrator().CreateTable(&User{})
}

func migratePostTable(SQDB *gorm.DB, MQDB *gorm.DB) {
	SQDB.AutoMigrate(&Post{})
	SQDB.Migrator().DropTable(&Post{})
	SQDB.Migrator().CreateTable(&Post{})
	MQDB.AutoMigrate(&Post{})
	MQDB.Migrator().DropTable(&Post{})
	MQDB.Migrator().CreateTable(&Post{})
}

func createUser(SQDB *gorm.DB, MQDB *gorm.DB) {
	var newUser = User{
		Id: 1,
		Name: "RIO",
		Email: "popo62520908@gmail.com",
		Password: "52f96e51831c8229413e28b0e58fa3b992f7571e4ff5bf5ccfc1a21f391e4f05",
		Image: "CWCM67iUYAAZ1Kp.png",
		CreatedAt: time.Now(),
	}
	SQDB.Create(&newUser)
	MQDB.Create(&newUser)
}

func createUsers(SQDB *gorm.DB, MQDB *gorm.DB) {
	deleteUsersImages()
	byteArray, err := ioutil.ReadFile("data/users.json"); if err != nil {
		fmt.Println(err)
	}
	var users []User
	err = json.Unmarshal(byteArray, &users); if err != nil {
		fmt.Println(err)
	}
	result := make([]User, 0)
	now := time.Now()
	fmt.Println(now)
	for i, u := range users {
		u.Id = i + 2
		u.Password = "15e2b0d3c33891ebb0f1ef609ec419420c20e320ce94c65fbc8c3312448eb225"
		u.CreatedAt = time.Now()
		// u.getImage()
		randStr, _ := models.MakeRandomStr(20)
		u.Image = randStr + ".png"
		image.CreateImage(u.Name, u.Image)
		result = append(result, u)
	}
	SQDB.Create(&result)
	MQDB.Create(&result)
	now = time.Now()
	fmt.Println(now)
}

func createPost(SQDB *gorm.DB, MQDB *gorm.DB) {
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
	SQDB.Create(&post1)
	SQDB.Create(&post2)
	MQDB.Create(&post1)
	MQDB.Create(&post2)
}

func createPosts(SQDB *gorm.DB, MQDB *gorm.DB) {
	deletePostsImages()
	byteArray, err := ioutil.ReadFile("data/posts.json"); if err != nil {
		fmt.Println(err)
	}
	var posts []Post
	err = json.Unmarshal(byteArray, &posts); if err != nil {
		fmt.Println(err)
	}
	result := make([]Post, 0)
	now := time.Now()
	fmt.Println(now)
	for i, p := range posts {
		p.Id = i + 3
		p.Language = "English"
		p.UserID = randomInt(2, 102)
		p.CreatedAt = time.Now()
		p.getImage()
		result = append(result, p)
	}
	MQDB.Create(&result)
	SQDB.Create(&result)
	now = time.Now()
	fmt.Println(now)
}

func readUser(db *gorm.DB) {
	var user User
	db.Preload("Posts").Where("email = ?", "popo62520908@gmail.com").First(&user)
	fmt.Println(user)
}

func readPost(db *gorm.DB) {
	var post Post
	db.Preload("User").Find(&post)
	fmt.Println(post)
}


func (u *User)getImage() {
	url := "https://loremflickr.com/320/240?random=1"

	response, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()

	randStr, _ := models.MakeRandomStr(15)
	extension := ".png"
	filename := randStr + extension
	path := "upload/user/"

	file, err := os.Create(path + filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	io.Copy(file, response.Body)
	u.Image = filename
}

func (p *Post)getImage() {
	url := "https://loremflickr.com/320/240?random=1"

	response, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()

	randStr, _ := models.MakeRandomStr(15)
	extension := ".png"
	filename := randStr + extension
	path := "upload/post/"

	file, err := os.Create(path + filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	io.Copy(file, response.Body)
	p.Image = filename
}

func deleteUsersImages() {
	files, _ := ioutil.ReadDir("./upload/user")
	if len(files) != 0 {
		for _, v := range files {
			if v.Name() != "CWCM67iUYAAZ1Kp.png" {
				os.Remove("./upload/user/" + v.Name())
			}
		}
	}
}

func deletePostsImages() {
	files, _ := ioutil.ReadDir("./upload/post")
	if len(files) != 0 {
		for _, v := range files {
			if v.Name() != "abc.png" && v.Name() != "abcd.png"{
				os.Remove("./upload/post/" + v.Name())
			}
		}
	}
}

func randomInt(min, max int) int {
	return min + rand.Intn(max-min)
}