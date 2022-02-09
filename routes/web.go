package routes

import (
	"helloworld/config"
	"helloworld/controller"
	"log"
	"net/http"
	"golang.org/x/net/websocket"
)

var home controller.Home
var auth controller.Auth
var post controller.Post
var users controller.Users
var infolog *log.Logger

func init() {
	infolog = config.App.InfoLog
}

func Routes() http.Handler{
	mux := http.NewServeMux()

	// static
	staticFiles := http.FileServer(http.Dir(config.App.Static))
	uploadFiles := http.FileServer(http.Dir(config.App.Media))
	mux.Handle("/static/", http.StripPrefix("/static/", staticFiles))
	mux.Handle("/media/", http.StripPrefix("/media/", uploadFiles))

	// normal
	mux.HandleFunc("/login", auth.Login)
	mux.HandleFunc("/sign-up", auth.Register)
	mux.HandleFunc("/oauth/callback", auth.GitHubLogin)

	// Auth
	mux.Handle("/", Auth(http.HandlerFunc(home.Index)))
	mux.Handle("/logout", Auth(http.HandlerFunc(auth.Logout)))
	mux.Handle("/post", Auth(http.HandlerFunc(post.Index)))
	mux.Handle("/users/", Auth(http.HandlerFunc(users.Index)))

	// websocket
	mux.Handle("/chat", websocket.Handler(home.Chat))

	wrappedMux := NewLogger(mux)
	return wrappedMux
}
