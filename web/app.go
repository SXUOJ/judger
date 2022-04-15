package web

import (
	"github.com/SXUOJ/judge/config"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type App struct {
	Conf   *config.WebConf
	Router *gin.Engine
}

func NewApp() *App {
	return &App{
		Conf:   config.LoadConf(),
		Router: loadRouter(),
	}
}

func (app *App) Run() {
	logrus.Print("Wechat-mall-backend runs on http://" + app.Conf.Listen)
	app.Router.Run(app.Conf.Listen)
}
