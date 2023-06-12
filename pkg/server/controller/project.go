package controller

import (
	"github.com/gin-gonic/gin"
)

// CreateProject is the handler function of the /project/create endpoint.
func CreateProject(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}
