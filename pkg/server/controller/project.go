package controller

import (
	"github.com/akosgarai/wasm-example/pkg/server/request"
	"github.com/gin-gonic/gin"
)

// CreateProject is the handler function of the /project/create endpoint.
func CreateProject(c *gin.Context) {
	req := &request.CreateProjectRequest{}
	err := c.BindJSON(req)
	if err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}
	validationErrors := req.Validate()
	if len(validationErrors) > 0 {
		c.JSON(400, gin.H{
			"error": validationErrors,
		})
		return
	}
	// start the project creation in the background
	// then return the connection to the websocket
	// that will be used to send the progress
	c.JSON(200, gin.H{
		"message": "pong",
		"req":     req,
	})
}
