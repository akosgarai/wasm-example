package controller

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/akosgarai/wasm-example/pkg/server/request"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type response struct {
	Error interface{}
	Data  interface{}
}

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
		response := processMessage(msg)
		marshalled, _ := json.Marshal(response)
		conn.WriteMessage(t, marshalled)
	}
}

func processMessage(msg []byte) response {
	// msg is a json marshalled string, so we need to unmarshal it
	// and use the data to create the project
	var resp response
	unmarshalled := &request.CreateProjectRequest{}
	err := json.Unmarshal(msg, unmarshalled)
	if err != nil {
		errMsg := fmt.Sprintf("Failed to unmarshal json: %+v", err)
		fmt.Printf(errMsg)
		resp.Error = errMsg
		return resp
	}
	// Data validation
	validationErrors := unmarshalled.Validate()
	if len(validationErrors) > 0 {
		fmt.Printf("Validation errors: %+v", validationErrors)
		resp.Error = validationErrors
		return resp
	}
	// Do the process stuff right here.
	// Now just return the message.
	resp.Data = "Project created successfully."
	return resp
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
