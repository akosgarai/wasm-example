package main

import (
	"fmt"
	"os"

	"github.com/akosgarai/wasm-example/pkg/server/controller"
	"github.com/akosgarai/wasm-example/pkg/server/models"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func initDatabaseConnection() *gorm.DB {
	user := os.Getenv("MYSQL_USER")
	passwd := os.Getenv("MYSQL_PASSWORD")
	host := os.Getenv("MYSQL_HOST")
	dbname := os.Getenv("MYSQL_DATABASE")
	connectionString := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local", user, passwd, host, dbname)

	db, err := gorm.Open(mysql.Open(connectionString), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to database!")
	}
	// Execute the migrations
	db.AutoMigrate(&models.Project{}, &models.Runtime{}, &models.Dbtype{}, &models.Client{}, &models.Environment{})
	return db
}

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
	// Add the databse to the controller package as package variable.
	appController := controller.NewAppController(initDatabaseConnection())
	options := r.Group("/options")
	{
		options.GET("/projects", appController.ProjectNames)
		options.POST("/projects", appController.ProjectNamesWithQuery)
		options.GET("/runtimes", appController.ProjectRuntimes)
		options.GET("/databases", appController.ProjectDatabases)
	}
	r.GET("/ws", func(c *gin.Context) {
		appController.WsHandler(c.Writer, c.Request)
	})
	r.Run(":9090")
}
