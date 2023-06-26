package controller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/akosgarai/wasm-example/pkg/server/request"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"golang.org/x/crypto/ssh"
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
	// if the unmarshalled.Staging is true, we need to execute the staging command
	responseString := ""
	if unmarshalled.Staging != "false" {
		responseString = executeStagingCommand(unmarshalled)
	}
	// if the unmarshalled.Production is true, we need to execute the production command
	if unmarshalled.Production != "false" {
		responseString += executeProductionCommand(unmarshalled)
	}
	resp.Data = responseString

	return resp
}

func executeStagingCommand(data *request.CreateProjectRequest) string {
	sshUser := "scriptexecutor"
	sshHost := "staging"
	sshPort := "2222"
	sshKey := "/root/.ssh/id_rsa_shared"
	return executeServerCommand(sshUser, sshHost, sshPort, sshKey, data)
}
func executeProductionCommand(data *request.CreateProjectRequest) string {
	sshUser := "scriptexecutor"
	sshHost := "production"
	sshPort := "2222"
	sshKey := "/root/.ssh/id_rsa_shared"
	return executeServerCommand(sshUser, sshHost, sshPort, sshKey, data)
}

func executeServerCommand(sshUser, sshHost, sshPort, sshKey string, data *request.CreateProjectRequest) string {
	key, err := ioutil.ReadFile(sshKey)
	if err != nil {
		return fmt.Sprintf("unable to read private key: %v", err)
	}
	// Create the Signer for this private key.
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		return fmt.Sprintf("unable to parse private key: %v", err)
	}
	config := &ssh.ClientConfig{
		User: sshUser,
		Auth: []ssh.AuthMethod{
			// Add in password check here for moar security.
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	// Connect to the remote server and perform the SSH handshake.
	client, err := ssh.Dial("tcp", sshHost+":"+sshPort, config)
	if err != nil {
		return fmt.Sprintf("unable to connect to %s:%s with %s: %v", sshHost, sshPort, sshUser, err)
	}
	defer client.Close()
	ss, err := client.NewSession()
	if err != nil {
		return fmt.Sprintf("unable to create SSH session: %v", err)
	}
	defer ss.Close()
	// Creating the buffer which will hold the remotly executed command's output.
	var stdoutBuf bytes.Buffer
	ss.Stdout = &stdoutBuf
	cmdString := fmt.Sprintf("/usr/local/bin/setup-project.sh %s %s %s %s %s", data.Client, data.Name, data.Runtime, data.Database, data.OwnerEmail)
	ss.Run(cmdString)
	return stdoutBuf.String()
}

// ProjectNames is the handler function of the /projects endpoint.
// It returns the list of the projects for the select component.
func ProjectNames(c *gin.Context) {
	c.JSON(200, gin.H{
		"data": orderedOptionListFromSlice(ProjectList),
	})
}

// ProjectNamesWithQuery is the handler function of the /projects endpoint.
// It returns the list of the projects for the select component.
func ProjectNamesWithQuery(c *gin.Context) {
	queryRequest := request.QueryRequest{}
	if err := c.ShouldBindJSON(&queryRequest); err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}
	list := ProjectList
	// if we have a query, we need to filter the list
	if queryRequest.Query != "" {
		list = []string{}
		for _, v := range ProjectList {
			if strings.Contains(v, queryRequest.Query) {
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
