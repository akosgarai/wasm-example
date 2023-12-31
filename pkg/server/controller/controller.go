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

// ProjectEnvironments is the handler function of the /environments endpoint.
// It returns the list of the environments for the select component.
func (app *AppController) ProjectEnvironments(c *gin.Context) {
	environmentList := []models.Environment{}
	app.db.Find(&environmentList)
	c.JSON(200, gin.H{
		"data": environmentList,
	})
}

// ProjectNames is the handler function of the /projects endpoint.
// It returns the list of the projects for the select component.
func (app *AppController) ProjectNames(c *gin.Context) {
	projectList := []models.Project{}
	app.db.Find(&projectList)
	c.JSON(200, gin.H{
		"data": app.projectTransform(projectList),
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
		"data": app.projectTransform(projectList),
	})
}

// ClientNames is the handler function of the /clients endpoint.
// It returns the list of the clients for the select component.
func (app *AppController) ClientNames(c *gin.Context) {
	clientList := []models.Client{}
	app.db.Find(&clientList)
	c.JSON(200, gin.H{
		"data": app.clientTransform(clientList),
	})
}

// ClientNamesWithQuery is the handler function of the /clients endpoint.
// It returns the list of the projects for the select component.
func (app *AppController) ClientNamesWithQuery(c *gin.Context) {
	queryRequest := request.QueryRequest{}
	if err := c.ShouldBindJSON(&queryRequest); err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}
	if queryRequest.Query == "" {
		app.ClientNames(c)
		return
	}
	clientList := []models.Client{}
	app.db.Where("name LIKE ?", "%"+queryRequest.Query+"%").Find(&clientList)
	c.JSON(200, gin.H{
		"data": app.clientTransform(clientList),
	})
}

// projectTransform changes the output structure to make it fit to the frontend.
func (app *AppController) projectTransform(projects []models.Project) []map[string]interface{} {
	var result []map[string]interface{}
	for _, project := range projects {
		result = append(result, map[string]interface{}{
			"id":   project.Name,
			"name": project.Name,
		})
	}
	return result
}

// clientTransform changes the output structure to make it fit to the frontend.
func (app *AppController) clientTransform(clients []models.Client) []map[string]interface{} {
	var result []map[string]interface{}
	for _, client := range clients {
		result = append(result, map[string]interface{}{
			"id":   client.Name,
			"name": client.Name,
		})
	}
	return result
}
