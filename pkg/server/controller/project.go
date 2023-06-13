package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

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
var (
	// ProjectList is the list of the projects.
	ProjectList = []string{"asd", "asf", "ass", "bobo", "cat", "mkp"}
	// RuntimeList is the list of the possible runtimes
	RuntimeList = []string{"NoPHP", "PHP71FPM", "PHP74FPM", "PHP81FPM"}
)

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
		response := processMessage(msg, conn)
		marshalled, _ := json.Marshal(response)
		conn.WriteMessage(t, marshalled)
	}
}

func processMessage(msg []byte, conn *websocket.Conn) response {
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

// ProjectNames is the handler function of the /projects endpoint.
// It returns the list of the projects for the select component.
func ProjectNames(c *gin.Context) {
	list := ProjectList
	// if we have a query, we need to filter the list
	query := c.Query("query")
	if query != "" {
		list = []string{}
		for _, v := range ProjectList {
			if strings.Contains(v, query) {
				list = append(list, v)
			}
		}
	}
	c.JSON(200, gin.H{
		"data": orderedOptionListFromSlice(list),
	})
}

// ProjectRuntimes is the handler function of the /runtimes endpoint.
// It returns the list of the runtimes for the select component.
func ProjectRuntimes(c *gin.Context) {
	c.JSON(200, gin.H{
		"data": orderedOptionListFromSlice(RuntimeList),
	})
}

// ProjectDatabases is the handler function of the /databases endpoint.
// It returns the list of the databases for the select component.
func ProjectDatabases(c *gin.Context) {
	c.JSON(200, gin.H{
		"data": orderedOptionListFromSlice([]string{"no", "mysql"}),
	})
}

func orderedOptionListFromSlice(list []string) map[string]string {
	options := make(map[string]string)
	for _, v := range list {
		options[v] = v
	}
	return options
}
