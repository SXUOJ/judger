package service

import (
	"net/http"
	"os"

	"github.com/SXUOJ/judge/main/model"
	"github.com/gin-gonic/gin"
)

func loadRouter() *gin.Engine {
	r := gin.Default()

	r.GET("/ping", ping)
	r.POST("/submit", submit)

	return r
}

func ping(c *gin.Context) {
	c.JSON(200, gin.H{
		"msg": "pong",
	})
}

func submit(c *gin.Context) {
	submit := model.Submit{}
	if err := c.ShouldBindJSON(&submit); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"msg": "bind model error",
		})
		return
	}

	judger, err := submit.Load()
	defer remove(judger.WorkDir)
	if err != nil {
		c.JSON(200, gin.H{
			"msg": "submit.Load() failed",
		})
	}
	judger.Run(c)
	return
}

func remove(path string) error {
	if path != "" {
		return os.Remove(path)
	}
	return nil
}
