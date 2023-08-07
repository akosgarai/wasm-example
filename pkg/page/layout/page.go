package layout

import (
	"encoding/json"
	"fmt"
	"syscall/js"

	"github.com/akosgarai/wasm-example/pkg/client/dom"
	"github.com/akosgarai/wasm-example/pkg/page"
)

const (
	socketURL = "ws://localhost:9090/ws"
)

// Layout is the page for the layout.
type Layout struct {
	*page.Instance
	socket js.Value
}

type formItem struct {
	Tag        string
	Attributes map[string]interface{}
}

// map: form name -> map: form tag -> form item
var formItems = []formItem{
	{"checkbox", map[string]interface{}{
		"id":    "env-staging",
		"name":  "env-staging",
		"title": "Staging environment",
		"label": "Staging environment",
	}},
	{"checkbox", map[string]interface{}{
		"id":    "env-production",
		"name":  "env-production",
		"title": "Production environment",
		"label": "Production environment",
	}},
	{"input", map[string]interface{}{
		"id":          "project-client",
		"name":        "project-client",
		"type":        "text",
		"placeholder": "[a-z0-9-]",
		"title":       "Project client [a-z0-9-]",
		"label":       "Project client",
	}},
	{"select", map[string]interface{}{
		"id":          "project-name",
		"name":        "project-name",
		"type":        "api",
		"placeholder": "[a-z0-9-]",
		"title":       "Project name [a-z0-9-]",
		"label":       "Project name",
		"apiUrl":      "/options/projects/",
	}},
	{"input", map[string]interface{}{
		"id":          "project-owner-email",
		"name":        "project-owner-email",
		"type":        "email",
		"placeholder": "example@email.com",
		"label":       "Project owner email",
		"value":       "example@email.com",
	}},
	{"select", map[string]interface{}{
		"id":     "project-runtime",
		"name":   "project-runtime",
		"type":   "apisimple",
		"label":  "Project runtime",
		"title":  "Project runtime",
		"apiUrl": "/options/runtimes/",
	}},
	{"select", map[string]interface{}{
		"id":     "project-database",
		"name":   "project-database",
		"type":   "apisimple",
		"label":  "Project database",
		"title":  "Project database",
		"apiUrl": "/options/databases/",
	}},
}

// New returns a new Layout
func New(title string) *Layout {
	return &Layout{
		Instance: page.NewPage(title),
	}
}

// LoadPage loads the layout page.
func (l *Layout) LoadPage() {
	l.Instance.LoadPage()
	container := l.Document().Call("querySelector", "."+page.ContentClassName)
	inputContainer := dom.Div(l.Document(), map[string]interface{}{
		"className": "row",
	})
	form := dom.Form(l.Document(), map[string]interface{}{
		"id": "project-form",
	})
	container.Call("appendChild", form)
	form.Call("appendChild", inputContainer)
	for _, item := range formItems {
		inputContainer.Call("appendChild", l.buildFormItem(item.Tag, item.Attributes))
	}
	submitContainer := dom.Div(l.Document(), map[string]interface{}{
		"className": "row submit",
	})
	submit := l.buildFormItem("input", map[string]interface{}{
		"id":    "submit",
		"name":  "submit",
		"type":  "button",
		"title": "Submit",
		"value": "Submit",
	})
	submit.Set("onclick", l.submitForm().Call("bind", submit))
	submitContainer.Call("appendChild", submit)
	form.Call("appendChild", submitContainer)
	// Add the script execution result container.
	stagingErrorContainer := dom.Div(l.Document(), map[string]interface{}{
		"id":        "staging-error",
		"className": "row error",
	})
	stagingErrorContainer.Call("appendChild", dom.P(l.Document(), ""))
	productionErrorContainer := dom.Div(l.Document(), map[string]interface{}{
		"id":        "production-error",
		"className": "row error",
	})
	productionErrorContainer.Call("appendChild", dom.P(l.Document(), ""))
	stagingResultContainer := dom.Div(l.Document(), map[string]interface{}{
		"id":        "staging-path",
		"className": "row result",
	})
	stagingResultContainer.Call("appendChild", dom.P(l.Document(), ""))
	productionResultContainer := dom.Div(l.Document(), map[string]interface{}{
		"id":        "production-path",
		"className": "row result",
	})
	productionResultContainer.Call("appendChild", dom.P(l.Document(), ""))
	container.Call("appendChild", stagingErrorContainer)
	container.Call("appendChild", productionErrorContainer)
	container.Call("appendChild", stagingResultContainer)
	container.Call("appendChild", productionResultContainer)

	// create the socket
	l.socket = l.WebSocket().New(socketURL)
	l.socket.Set("onmessage", l.socketMessage())
}

// Run runs the formatter page.
func (l *Layout) Run() {
	<-make(chan bool)
}

// buildFormItem returns a form item.
func (l *Layout) buildFormItem(tag string, attributes map[string]interface{}) js.Value {
	var formItem js.Value
	switch tag {
	case "input":
		formItem = l.buildInputFormItem(attributes["type"].(string), attributes)
		break
	case "select":
		formItem = l.buildSelectFormItem(attributes)
		break
	case "checkbox":
		formItem = l.buildCheckboxFormItem(attributes)
		break
	}
	return formItem
}

// submitForm submits the form.
// jsonWrapper returns the wrapper function for the JSON formatter.
func (l *Layout) submitForm() js.Func {
	jsonFunc := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		go func() {
			envStaging := "false"
			if l.Document().Call("querySelector", "#env-staging").Get("checked").Bool() {
				envStaging = "true"
			}
			envProduction := "false"
			if l.Document().Call("querySelector", "#env-production").Get("checked").Bool() {
				envProduction = "true"
			}
			// call the /ping endpoint with the form data
			// the response has to be logged to the console
			projectData := map[string]string{
				"command":             "create-project",
				"project-name":        l.Document().Call("querySelector", "#project-name").Get("value").String(),
				"project-client":      l.Document().Call("querySelector", "#project-client").Get("value").String(),
				"project-owner-email": l.Document().Call("querySelector", "#project-owner-email").Get("value").String(),
				"project-runtime":     l.Document().Call("querySelector", "#project-runtime").Get("value").String(),
				"project-database":    l.Document().Call("querySelector", "#project-database").Get("value").String(),
				"env-staging":         envStaging,
				"env-production":      envProduction,
			}
			jsonStr, err := json.Marshal(projectData)
			if err != nil {
				fmt.Println(err)
			}
			l.clearErrorMessages()
			l.socket.Call("send", string(jsonStr))
		}()
		return nil
	})
	return jsonFunc
}

// socketMessage returns the wrapper function for the socket message.
func (l *Layout) socketMessage() js.Func {
	jsonFunc := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		go func() {
			// the response is a JSON string
			response := args[0].Get("data").String()
			// parse the JSON string
			var responseMap map[string]interface{}
			err := json.Unmarshal([]byte(response), &responseMap)
			if err != nil {
				l.Alert().Invoke(fmt.Sprintf("Error unmarshall the response. %s", err.Error()))
			}
			// if the response has the "Error" key, the error has to be logged to the console
			if responseMap["Error"] != nil {
				// write the error message to their containers
				l.addErrors(responseMap["Error"].(map[string]interface{}))
				return
			}
			// if the response has the "Data" key,
			if responseMap["Data"] == nil {
				l.Alert().Invoke(fmt.Sprintf("No data in the response"))
				return
			}
			responseData := responseMap["Data"].(map[string]interface{})
			for key, value := range responseData {
				// write the response to the result containers
				resultContainer := l.Document().Call("querySelector", "#"+key+" p")
				resultContainer.Set("innerText", value)
			}
		}()
		return nil
	})
	return jsonFunc
}

// clear error messages removes the current error messages.
func (l *Layout) clearErrorMessages() {
	for _, item := range formItems {
		errorMessageContainer := l.Document().Call("querySelector", "#"+item.Attributes["id"].(string)+"-error-message")
		errorMessageContainer.Set("innerText", "")
	}
	// clear the script execution messages also.
	selectors := []string{"staging-error", "staging-path", "production-error", "production-path"}
	for _, item := range selectors {
		executionMessageContainer := l.Document().Call("querySelector", "#"+item+" p")
		executionMessageContainer.Set("innerText", "")
	}
}

// add errors adds the errors to the form.
func (l *Layout) addErrors(errorMapInterface map[string]interface{}) {
	for _, item := range formItems {
		attrID := item.Attributes["id"].(string)
		if errorMapInterface[attrID] != nil {
			errorMap := errorMapInterface[attrID].([]interface{})
			errorMessageContainer := l.Document().Call("querySelector", "#"+attrID+"-error-message")
			for _, message := range errorMap {
				msgParagraph := dom.P(l.Document(), message.(string))
				errorMessageContainer.Call("appendChild", msgParagraph)
			}
		}
	}
}

// buildSelectFormItem returns a select form item.
func (l *Layout) buildSelectFormItem(attributes map[string]interface{}) js.Value {
	id := attributes["id"].(string)
	var selector js.Value
	if attributes["type"] == "api" {
		selector = dom.APISelect(l.Document(), attributes["apiUrl"].(string), id, "")
	} else if attributes["type"] == "apisimple" {
		// Get the option values from the given API url.
		selector = dom.SimpleAPISelect(l.Document(), attributes["apiUrl"].(string), id, "")
	}
	itemContainer := dom.Div(l.Document(), map[string]interface{}{"className": "form-item", "id": attributes["id"].(string) + "-container"})
	// if we have label, we have to create it and append it to the itemContainer
	if attributes["label"] != nil {
		label := dom.Label(l.Document(), attributes["label"].(string), attributes["id"].(string))
		itemContainer.Call("appendChild", label)
	}
	itemContainer.Call("appendChild", selector)
	// add the error message container
	errorMessageContainer := dom.Div(l.Document(), map[string]interface{}{"className": "error-message", "id": attributes["id"].(string) + "-error-message"})
	itemContainer.Call("appendChild", errorMessageContainer)
	return itemContainer
}

// buildInputFormItem returns an input form item.
func (l *Layout) buildInputFormItem(inputType string, attributes map[string]interface{}) js.Value {
	element := dom.Input(l.Document(), inputType, attributes)
	itemContainer := dom.Div(l.Document(), map[string]interface{}{"className": "form-item", "id": attributes["id"].(string) + "-container"})
	// if we have label, we have to create it and append it to the itemContainer
	if attributes["label"] != nil {
		label := dom.Label(l.Document(), attributes["label"].(string), attributes["id"].(string))
		itemContainer.Call("appendChild", label)
	}
	itemContainer.Call("appendChild", element)
	// add the error message container
	errorMessageContainer := dom.Div(l.Document(), map[string]interface{}{"className": "error-message", "id": attributes["id"].(string) + "-error-message"})
	itemContainer.Call("appendChild", errorMessageContainer)
	return itemContainer
}

// buildCheckboxFormItem returns a checkbox form item.
func (l *Layout) buildCheckboxFormItem(attributes map[string]interface{}) js.Value {
	itemContainer := dom.Div(l.Document(), map[string]interface{}{"className": "form-item center", "id": attributes["id"].(string) + "-container"})
	// append the checkbox
	checkbox := dom.CheckBox(l.Document(), attributes["id"].(string), attributes["label"].(string), false)
	itemContainer.Call("appendChild", checkbox)
	// add the error message container
	errorMessageContainer := dom.Div(l.Document(), map[string]interface{}{"className": "error-message", "id": attributes["id"].(string) + "-error-message"})
	itemContainer.Call("appendChild", errorMessageContainer)
	return itemContainer
}
