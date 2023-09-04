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
	db.AutoMigrate(&models.Project{}, &models.Runtime{}, &models.Dbtype{}, &models.Client{}, &models.Environment{}, &models.Host{})
	// Seed the database with the default values.
	dbNames := []string{"mysql", "no"}
	for _, name := range dbNames {
		var dbtype models.Dbtype
		dbtype.Name = name
		if err := db.Where("name = ?", name).FirstOrCreate(&dbtype).Error; err != nil {
			panic("Failed to seed database!")
		}
	}
	runtimeNames := []string{"noPHP", "PHP71FPM", "PHP74FPM", "PHP81FPM"}
	for _, name := range runtimeNames {
		var runtime models.Runtime
		runtime.Name = name
		if err := db.Where("name = ?", name).FirstOrCreate(&runtime).Error; err != nil {
			panic("Failed to seed database!")
		}
	}
	var stagingEnvironment models.Environment
	stagingHostName := "staging" // the service name in the docker-compose file
	stagingEnvironment.Name = stagingHostName
	var productionEnvironment models.Environment
	productionHostName := "production" // the service name in the docker-compose file
	productionEnvironment.Name = productionHostName
	if err := db.Where("name = ?", stagingHostName).FirstOrCreate(&stagingEnvironment).Error; err != nil {
		panic("Failed to seed database!")
	}
	if err := db.Where("name = ?", productionHostName).FirstOrCreate(&productionEnvironment).Error; err != nil {
		panic("Failed to seed database!")
	}
	var stagingHost models.Host
	stagingHost.IP = stagingHostName
	stagingHost.EnvironmentID = stagingEnvironment.ID
	stagingHost.Name = "server-staging"
	var productionHost models.Host
	productionHost.IP = productionHostName
	productionHost.EnvironmentID = productionEnvironment.ID
	productionHost.Name = "server-production"
	if err := db.Where("ip = ?", stagingHostName).Where("environment_id = ?", stagingEnvironment.ID).Attrs(stagingHost).FirstOrCreate(&stagingHost).Error; err != nil {
		panic("Failed to seed database!")
	}
	if err := db.Where("ip = ?", productionHostName).Where("environment_id = ?", productionEnvironment.ID).Attrs(productionHost).FirstOrCreate(&productionHost).Error; err != nil {
		panic("Failed to seed database!")
	}
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
		options.GET("/environments", appController.ProjectEnvironments)
	}
	r.GET("/ws", func(c *gin.Context) {
		appController.WsHandler(c.Writer, c.Request)
	})
	r.Run(":9090")
}
