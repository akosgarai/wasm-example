package main

import (
	"github.com/akosgarai/wasm-example/pkg/server/controller"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.Static("/assets", "../../assets")
	r.LoadHTMLGlob("../../assets/*.html")
	r.GET("/", func(c *gin.Context) {
		c.HTML(200, "index.html", gin.H{})
	})
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.POST("/project/create", controller.CreateProject)
	r.Run(":9090")
}
