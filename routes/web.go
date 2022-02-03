package routes

import (
	"net/http"
	"helloworld/config"
	"helloworld/controller/home"
	"helloworld/controller/auth"
	"helloworld/controller/post"
	"helloworld/controller/users"
	"golang.org/x/net/websocket"
)

func Routes() http.Handler{
	mux := http.NewServeMux()
	staticFiles := http.FileServer(http.Dir(config.App.Static))
	uploadFiles := http.FileServer(http.Dir("upload"))
	mux.Handle("/static/", http.StripPrefix("/static/", staticFiles))
	mux.Handle("/media/", http.StripPrefix("/media/", uploadFiles))
	mux.HandleFunc("/", home.Index)
	mux.HandleFunc("/login", auth.Login)
	mux.HandleFunc("/post", post.Index)
	mux.HandleFunc("/sign-up", auth.Register)
	mux.HandleFunc("/logout", auth.Logout)
	mux.HandleFunc("/oauth/callback", auth.GitHubLogin)
	mux.Handle("/chat", websocket.Handler(home.Chat))
	mux.HandleFunc("/users/", users.Index)
	return mux
}