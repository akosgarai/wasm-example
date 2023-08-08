package controller

import (
	"github.com/akosgarai/wasm-example/pkg/server/models"
	"github.com/akosgarai/wasm-example/pkg/server/request"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"gorm.io/gorm"
)

// AppController is the controller of the application.
// It contains the database connection.
type AppController struct {
	// use the gorm.DB type here
	db *gorm.DB
	// websocket upgrader configuration
	wsupgrader websocket.Upgrader
}

// NewAppController returns a new AppController instance.
func NewAppController(db *gorm.DB) *AppController {
	return &AppController{
		db: db,
		wsupgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		},
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

// ProjectNamesWithQuery is the handler function of the /projects endpoint.
// It returns the list of the projects for the select component.
func (app *AppController) ProjectNamesWithQuery(c *gin.Context) {
	queryRequest := request.QueryRequest{}
	if err := c.ShouldBindJSON(&queryRequest); err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}
	if queryRequest.Query == "" {
		app.ProjectNames(c)
		return
	}
	projectList := []models.Project{}
	app.db.Where("name LIKE ?", "%"+queryRequest.Query+"%").Find(&projectList)
	c.JSON(200, gin.H{
		"data": projectList,
	})
}
