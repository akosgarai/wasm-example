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
	r.GET("/projects", controller.ProjectNames)
	r.GET("/projects/:query", controller.ProjectNames)
	r.GET("/runtimes", controller.ProjectRuntimes)
	r.GET("/databases", controller.ProjectDatabases)
	r.GET("/ws", func(c *gin.Context) {
		controller.WsHandler(c.Writer, c.Request)
	})
	r.Run(":9090")
}
