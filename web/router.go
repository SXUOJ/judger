package web

import (
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
	c.JSON(200, gin.H{
		"msg": "pong",
	})
}
