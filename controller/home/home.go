package home

import (
	"errors"
	"fmt"
	"helloworld/config"
	"helloworld/controller"
	"helloworld/models"
	"log"
	"net/http"
	"github.com/jinzhu/gorm"
	"golang.org/x/net/websocket"
)

var infolog *log.Logger
var ws_array []*websocket.Conn // *websocket.Connを入れる配列

type Message struct {
	Name    string
	Message string
}

var errorlog *log.Logger
var DB *gorm.DB

func init() {
	infolog = config.App.InfoLog
	DB = models.DB
}

func Index(w http.ResponseWriter, r *http.Request) {
	// return Not found when the path is not "/" like "/123"
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	if r.Method != "GET" {
		http.NotFound(w, r)
		return
	}
	// session
	session, err := models.CheckSession(r)
	if err != nil {
		fmt.Sprintf("%v\t%v", r.URL, r.RemoteAddr)
		http.Redirect(w, r, "/login", 302)
		return
	}
	infolog.Print(fmt.Sprintf("%v\t%v\t%v\t%v", r.URL, session.Name, session.Email, r.RemoteAddr))
	var posts []models.Post
	result := DB.Preload("User").Find(&posts)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		errorlog.Print(result)
	}
	var users []models.User
	result = DB.Preload("Posts").Find(&users)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		errorlog.Print(result)
	}
	stringMap := make(map[string]string)
	stringMap["csrf_token"] = session.CSRFToken
	controller.RenderTemplate(w, r, "index.html", &controller.TemplateData{
		StringMap: stringMap,
		Posts: posts,
		Users: users,
	})
}

func Chat(ws *websocket.Conn) {
	ws_array = append(ws_array, ws)
    data_receive(ws)
}

func data_receive(ws *websocket.Conn) {
    for {
        var message Message
        if err := websocket.JSON.Receive(ws, &message); err != nil {
            log.Println("Receive error:", err)
			break
        } else {
            for _, con := range ws_array {
				con := con
                go func() {
                    err = websocket.JSON.Send(con, message)
                    log.Println("con:", con)
                    if err != nil {
                        log.Println("Send error:", err)
                    }
                }()
            }
        }
    }
}
