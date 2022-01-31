// Sample run-helloworld is a minimal Cloud Run service.
package main

import (
	_ "embed"
	"helloworld/config"
	"helloworld/routes"
	"net/http"
)

func main() {
	infolog := config.App.InfoLog
	infolog.Print("starting server...")
	server := http.Server{
		Addr: config.App.Addr,
		Handler: routes.Routes(),
	}
	infolog.Print("run server!!")
	server.ListenAndServe()
}
