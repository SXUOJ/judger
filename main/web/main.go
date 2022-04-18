package main

import "github.com/SXUOJ/judge/main/web/service"

func main() {
	app := service.NewApp()
	app.Run()
}
