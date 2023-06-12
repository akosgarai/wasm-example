package controller

import (
	"fmt"
	"net/http"

	"github.com/akosgarai/wasm-example/pkg/server/request"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var wsupgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// WsHandler is the handler function of the /ws endpoint.
func WsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := wsupgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Printf("Failed to set websocket upgrade: %+v", err)
		return
	}

	for {
		t, msg, err := conn.ReadMessage()
		if err != nil {
			break
		}
		fmt.Printf("Message received: %+v\n", string(msg))
		conn.WriteMessage(t, msg)
	}
}

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
