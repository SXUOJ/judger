package main

import "github.com/SXUOJ/judge/web"

var (
	showPrint   = true
	showDetails = true
)

func main() {
	app := web.NewApp()
	app.Run()
}
