package main

import (
	"os"

	"github.com/akosgarai/wasm-example/pkg/server/controller"
	"github.com/gin-gonic/gin"
)

func main() {
	// Get the assets directory from the environment variable.
	// If it is not set, use the default value.
	assetsDir := os.Getenv("ASSETS_DIR")
	if assetsDir == "" {
		// if the env is not set, we have to use the directory relative to this file.
		assetsDir = "../../assets"
	}
	r := gin.Default()
	r.Static("/assets", assetsDir)
	r.LoadHTMLGlob(assetsDir + "/*.html")
	r.GET("/", func(c *gin.Context) {
		c.HTML(200, "index.html", gin.H{})
	})
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	options := r.Group("/options")
	{
		options.GET("/projects", controller.ProjectNames)
		options.POST("/projects", controller.ProjectNamesWithQuery)
		options.GET("/runtimes", controller.ProjectRuntimes)
		options.GET("/databases", controller.ProjectDatabases)
	}
	r.GET("/ws", func(c *gin.Context) {
		controller.WsHandler(c.Writer, c.Request)
	})
	r.Run(":9090")
}
