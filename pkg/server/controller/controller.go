package controller

import (
	"github.com/akosgarai/wasm-example/pkg/server/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// AppController is the controller of the application.
// It contains the database connection.
type AppController struct {
	// use the gorm.DB type here
	db *gorm.DB
}

// NewAppController returns a new AppController instance.
func NewAppController(db *gorm.DB) *AppController {
	return &AppController{
		db: db,
	}
}

// ProjectDatabases is the handler function of the /databases endpoint.
// It returns the list of the databases for the select component.
func (app *AppController) ProjectDatabases(c *gin.Context) {
	databaseList := []models.Dbtype{}
	app.db.Find(&databaseList)
	c.JSON(200, gin.H{
		"data": databaseList,
	})
}

// ProjectRuntimes is the handler function of the /runtimes endpoint.
// It returns the list of the runtimes for the select component.
func (app *AppController) ProjectRuntimes(c *gin.Context) {
	runtimeList := []models.Runtime{}
	app.db.Find(&runtimeList)
	c.JSON(200, gin.H{
		"data": runtimeList,
	})
}

// ProjectNames is the handler function of the /projects endpoint.
// It returns the list of the projects for the select component.
func (app *AppController) ProjectNames(c *gin.Context) {
	projectList := []models.Project{}
	app.db.Find(&projectList)
	c.JSON(200, gin.H{
		"data": projectList,
	})
}
