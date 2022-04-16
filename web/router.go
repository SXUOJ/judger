package web

import (
	"net/http"

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
	submit := Submit{}
	if err := c.ShouldBindJSON(&submit); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"msg": "bind model error",
		})
		return
	}
	worker, err := submit.Load()
	if err != nil {
		c.JSON(200, gin.H{
			"msg": "submit.Load() failed",
		})
	}
	worker.Run(c)
}
