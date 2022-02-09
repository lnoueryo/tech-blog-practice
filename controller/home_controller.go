package controller

import (
	"helloworld/models"
	"log"
	"net/http"
	"golang.org/x/net/websocket"
)

var ws_array []*websocket.Conn // *websocket.Connを入れる配列

type Message struct {
	Name    string
	Message string
}

type Home struct {}


func (h *Home)Index(w http.ResponseWriter, r *http.Request) {
	// return Not found when the path is not "/" like "/123"
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	if r.Method != "GET" {
		http.NotFound(w, r)
		return
	}

	//users
	users, err := models.UserAll(); if err != nil {
		errorlog.Print(err)
	}

	//posts
	posts, err := models.PostAll(); if err != nil {
		errorlog.Print(err)
	}

	RenderTemplate(w, r, "index.html", &TemplateData{
		CSRFToken: models.GenerateCSRFToken(r),
		Posts: posts,
		Users: users,
	})
}

func (h *Home)Chat(ws *websocket.Conn) {
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
