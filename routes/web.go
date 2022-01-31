package routes

import (
	"net/http"
	"helloworld/config"
	"helloworld/controller/home"
	"helloworld/controller/auth"
	"golang.org/x/net/websocket"
)

func Routes() http.Handler{
	mux := http.NewServeMux()
	files := http.FileServer(http.Dir(config.App.Static))
	mux.Handle("/static/", http.StripPrefix("/static/", files))
	mux.HandleFunc("/", home.Index)
	mux.HandleFunc("/login", auth.Login)
	mux.HandleFunc("/sign-up", auth.Register)
	mux.HandleFunc("/logout", auth.Logout)
	mux.HandleFunc("/oauth/callback", auth.GitHubLogin)
	mux.Handle("/chat", websocket.Handler(home.Chat))
	return mux
}