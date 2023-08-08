package controller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/akosgarai/wasm-example/pkg/server/request"
	"github.com/akosgarai/wasm-example/pkg/server/response"
	"github.com/gorilla/websocket"
	"golang.org/x/crypto/ssh"
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
		// insert the project into the database
		responseString = app.executeStagingCommand(unmarshalled)
		// if the output is the success string of the script, then we have to set the path to the project
		if strings.Contains(responseString, "The project has been created.") {
			resp.Data["staging-error"] = ""
			resp.Data["staging-path"] = "localhost:9091/" + unmarshalled.Client + "/" + unmarshalled.Name
		} else {
			resp.Data["staging-error"] = "Staging: " + responseString
			resp.Data["staging-path"] = ""
			// if the creation of the project failed, we need to delete the project from the database
		}
	}
	responseString = ""
	// if the unmarshalled.Production is true, we need to execute the production command
	if unmarshalled.Production != "false" {
		// insert the project into the database
		responseString += app.executeProductionCommand(unmarshalled)
		// if the output is the success string of the script, then we have to set the path to the project
		if strings.Contains(responseString, "The project has been created.") {
			resp.Data["production-error"] = ""
			resp.Data["production-path"] = "localhost:9096/" + unmarshalled.Client + "/" + unmarshalled.Name
		} else {
			resp.Data["production-error"] = "Production: " + responseString
			resp.Data["production-path"] = ""
			// if the creation of the project failed, we need to delete the project from the database
		}
	}
	return resp
}

func (app *AppController) executeStagingCommand(data *request.CreateProjectRequest) string {
	sshUser := "scriptexecutor"
	sshHost := "staging"
	sshPort := "2222"
	sshKey := "/root/.ssh/id_rsa_shared"
	return app.executeServerCommand(sshUser, sshHost, sshPort, sshKey, data)
}
func (app *AppController) executeProductionCommand(data *request.CreateProjectRequest) string {
	sshUser := "scriptexecutor"
	sshHost := "production"
	sshPort := "2222"
	sshKey := "/root/.ssh/id_rsa_shared"
	return app.executeServerCommand(sshUser, sshHost, sshPort, sshKey, data)
}

func (app *AppController) executeServerCommand(sshUser, sshHost, sshPort, sshKey string, data *request.CreateProjectRequest) string {
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
