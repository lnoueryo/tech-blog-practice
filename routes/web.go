package routes

import (
	"helloworld/config"
	"helloworld/controller"
	"net/http"
	"golang.org/x/net/websocket"
)

var home controller.Home
var auth controller.Auth
var post controller.Post
var users controller.Users

func Routes() http.Handler{
	mux := http.NewServeMux()
	staticFiles := http.FileServer(http.Dir(config.App.Static))
	uploadFiles := http.FileServer(http.Dir("upload"))
	mux.Handle("/static/", http.StripPrefix("/static/", staticFiles))
	mux.Handle("/media/", http.StripPrefix("/media/", uploadFiles))

	mux.HandleFunc("/", home.Index)
	mux.HandleFunc("/login", auth.LoginIndex)
	mux.HandleFunc("/sign-up", auth.Register)
	mux.HandleFunc("/logout", auth.Logout)
	mux.HandleFunc("/oauth/callback", auth.GitHubLogin)
	mux.HandleFunc("/post", post.Index)
	mux.HandleFunc("/users/", users.Index)

	mux.Handle("/chat", websocket.Handler(home.Chat))
	return mux
}