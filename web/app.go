package web

import (
	"github.com/SXUOJ/judge/env"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type App struct {
	Conf   *env.Conf
	Router *gin.Engine
}

func NewApp() *App {
	return &App{
		Conf:   env.LoadConf(),
		Router: loadRouter(),
	}
}

func (app *App) Run() {
	logrus.Print("Wechat-mall-backend runs on http://" + app.Conf.Listen)
	app.Router.Run(app.Conf.Listen)
}
