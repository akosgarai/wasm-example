package controller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/akosgarai/wasm-example/pkg/server/models"
	"github.com/akosgarai/wasm-example/pkg/server/request"
	"github.com/akosgarai/wasm-example/pkg/server/response"
	"github.com/gorilla/websocket"
	"golang.org/x/crypto/ssh"
)

var (
	mapIPToLocalhost = map[string]string{
		"staging":    "localhost:9091",
		"production": "localhost:9096",
	}
)

// WsHandler is the handler function of the /ws endpoint.
func (app *AppController) WsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := app.wsupgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Printf("Failed to set websocket upgrade: %+v", err)
		return
	}

	for {
		t, msg, err := conn.ReadMessage()
		if err != nil {
			break
		}
		response := app.processMessage(msg, conn)
		marshalled, _ := json.Marshal(&response)
		conn.WriteMessage(t, marshalled)
	}
}

func (app *AppController) processMessage(msg []byte, conn *websocket.Conn) *response.Socket {
	// msg is a json marshalled string, so we need to unmarshal it
	// and use the data to create the project
	resp := response.NewSocket()
	unmarshalled := &request.CreateProjectRequest{}
	err := json.Unmarshal(msg, unmarshalled)
	if err != nil {
		errMsg := fmt.Sprintf("Failed to unmarshal json: %+v", err)
		fmt.Printf(errMsg)
		responseError := map[string]interface{}{
			"general": errMsg,
		}
		resp.Error = responseError
		return resp
	}
	// Data validation
	validationErrors := unmarshalled.Validate(app.db)
	if len(validationErrors) > 0 {
		fmt.Printf("Validation errors: %+v", validationErrors)
		resp.Error = validationErrors
		return resp
	}
	// create project on each environment
	for _, env := range unmarshalled.Environment {
		// Create the Application table entry
		appEntry := &models.Application{
			ProjectID:     unmarshalled.ProjectID,
			ClientID:      unmarshalled.ClientID,
			OwnerEmail:    unmarshalled.OwnerEmail,
			RuntimeID:     unmarshalled.Runtime,
			DatabaseID:    unmarshalled.Database,
			EnvironmentID: env.ID,
		}
		app.db.Create(appEntry)
		app.db.Preload("Environment").Preload("Runtime").Preload("Database").Preload("Project").Preload("Client").First(appEntry)
		// Create the project on the environment
		// Get the entry from the host table by the environment id
		hostEntry := &models.Host{EnvironmentID: uint(env.ID)}
		// This solution assumes that we only have one host entry for each environment.
		app.db.Preload("Environment").Where("environment_id = ?", uint(env.ID)).First(hostEntry)
		logMessage := fmt.Sprintf("Host entry: %s , %s/%s, %d", hostEntry.IP, appEntry.Client.Name, appEntry.Project.Name, env.ID)
		resp.Data["temp-log"] = logMessage
		responseString := app.executeServerCommand(hostEntry, appEntry)
		if strings.Contains(responseString, "The project has been created.") {
			resp.Data[appEntry.Environment.Name+"-error"] = ""
			resp.Data[appEntry.Environment.Name+"-path"] = mapIPToLocalhost[hostEntry.IP] + "/" + unmarshalled.Client + "/" + unmarshalled.Name
		} else {
			resp.Data[appEntry.Environment.Name+"-error"] = appEntry.Environment.Name + ": " + responseString
			resp.Data[appEntry.Environment.Name+"-path"] = ""
			// if the creation of the project failed, we need to delete the project from the database
		}
	}
	return resp
}

func (app *AppController) executeServerCommand(host *models.Host, data *models.Application) string {
	key, err := ioutil.ReadFile(host.SSHKey)
	if err != nil {
		return fmt.Sprintf("unable to read private key: %v - %v", err, host)
	}
	// Create the Signer for this private key.
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		return fmt.Sprintf("unable to parse private key: %v", err)
	}
	config := &ssh.ClientConfig{
		User: host.SSHUser,
		Auth: []ssh.AuthMethod{
			// Add in password check here for moar security.
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	// Connect to the remote server and perform the SSH handshake.
	client, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", host.IP, host.SSHPort), config)
	if err != nil {
		return fmt.Sprintf("unable to connect to %s:%d with %s: %v", host.IP, host.SSHPort, host.SSHUser, err)
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
	cmdString := fmt.Sprintf("/usr/local/bin/setup-project.sh %s %s %s %s %s", data.Client.Name, data.Project.Name, data.Runtime.Name, data.Database.Name, data.OwnerEmail)
	ss.Run(cmdString)
	return stdoutBuf.String()
}
