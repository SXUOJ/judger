package service

import (
	"github.com/SXUOJ/judge/main/web/env"
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
	logrus.Print("Wechat-mall-backend runs on http://" + app.Conf.Web.Listen)
	app.Router.Run(app.Conf.Web.Listen)
}
