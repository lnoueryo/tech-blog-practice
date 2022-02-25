package main

import (
	_ "embed"
	"fmt"
	"helloworld/config"
	"helloworld/routes"
	"net/http"
	"os/exec"
	"time"
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

func SetAlarm() {
	for {
		now := time.Now()
		nowstr := fmt.Sprintf("%v:%v:%v", now.Hour(), now.Minute(), now.Second())
		setTime := time.Date(2022, 9, 12, 1, 53, 0, 0, time.Local)
		setTimeStr := fmt.Sprintf("%v:%v:%v", setTime.Hour(), setTime.Minute(), setTime.Second())
		fmt.Println(nowstr)
		time.Sleep(time.Second * 1)
		if nowstr ==  setTimeStr {
			c := exec.Command("go", "run", "main.go", "make", "controller", "hello")
			c.Start()
			c.Wait()
		}
	}
}
