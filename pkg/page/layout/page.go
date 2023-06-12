package layout

import (
	"encoding/json"
	"fmt"
	"syscall/js"

	"github.com/akosgarai/wasm-example/pkg/page"
)

const (
	socketURL = "ws://localhost:9090/ws"
)

// Layout is the page for the layout.
type Layout struct {
	*page.Instance
}

type formItem struct {
	Tag        string
	Attributes map[string]interface{}
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
	formContainer := l.CreateElement("div", map[string]interface{}{
		"className": "row",
	})
	container.Call("appendChild", formContainer)
	form := l.CreateElement("form", map[string]interface{}{
		"id": "project-form",
	})
	formContainer.Call("appendChild", form)
	// map: form name -> map: form tag -> form item
	formItems := []formItem{
		{"input", map[string]interface{}{
			"id":          "project-client",
			"name":        "project-client",
			"type":        "text",
			"placeholder": "[a-z0-9-]",
			"title":       "Project client [a-z0-9-]",
			"label":       "Project client",
		}},
		{"input", map[string]interface{}{
			"id":          "project-name",
			"name":        "project-name",
			"type":        "text",
			"placeholder": "[a-z0-9-]",
			"title":       "Project name [a-z0-9-]",
			"label":       "Project name",
		}},
		{"input", map[string]interface{}{
			"id":          "project-owner-email",
			"name":        "project-owner-email",
			"type":        "email",
			"placeholder": "example@email.com",
			"label":       "Project owner email",
			"value":       "example@email.com",
		}},
		{"input", map[string]interface{}{
			"id":          "project-runtime",
			"name":        "project-runtime",
			"type":        "text",
			"placeholder": "NoPHP, PHP71FPM, PHP74FPM, PHP81FPM",
			"title":       "Options: NoPHP, PHP71FPM, PHP74FPM, PHP81FPM",
			"label":       "Project runtime",
		}},
		{"input", map[string]interface{}{
			"id":          "project-database",
			"name":        "project-database",
			"type":        "text",
			"placeholder": "no, mysql",
			"title":       "Options: no, mysql",
			"label":       "Project database",
		}},
		{"input", map[string]interface{}{
			"id":    "submit",
			"name":  "submit",
			"type":  "button",
			"title": "Submit",
			"value": "Submit",
		}},
	}
	for _, item := range formItems {
		form.Call("appendChild", l.buildFormItem(item.Tag, item.Attributes))
	}
	submit := form.Call("querySelector", "#submit")
	submit.Set("onclick", l.submitForm().Call("bind", submit))
}

// Run runs the formatter page.
func (l *Layout) Run() {
	<-make(chan bool)
}

// buildFormItem returns a form item.
func (l *Layout) buildFormItem(tag string, attributes map[string]interface{}) js.Value {
	element := l.CreateElement(tag, attributes)
	itemContainer := l.CreateElement("div", map[string]interface{}{"className": "form-item"})
	// if we have label, we have to create it and append it to the itemContainer
	if attributes["label"] != nil {
		label := l.CreateElement("label", map[string]interface{}{
			"htmlFor":   attributes["id"],
			"innerText": attributes["label"],
		})
		itemContainer.Call("appendChild", label)
	}
	itemContainer.Call("appendChild", element)
	return itemContainer
}

// submitForm submits the form.
// jsonWrapper returns the wrapper function for the JSON formatter.
func (l *Layout) submitForm() js.Func {
	jsonFunc := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		go func() {
			// call the /ping endpoint with the form data
			// the response has to be logged to the console
			projectData := map[string]string{
				"project-name":        l.Document().Call("querySelector", "#project-name").Get("value").String(),
				"project-client":      l.Document().Call("querySelector", "#project-client").Get("value").String(),
				"project-owner-email": l.Document().Call("querySelector", "#project-owner-email").Get("value").String(),
				"project-runtime":     l.Document().Call("querySelector", "#project-runtime").Get("value").String(),
				"project-database":    l.Document().Call("querySelector", "#project-database").Get("value").String(),
			}
			socket := l.WebSocket().New(socketURL)

			socket.Call("addEventListener", "open", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
				fmt.Println("open")
				jsonStr, err := json.Marshal(projectData)
				if err != nil {
					fmt.Println(err)
				}

				socket.Call("send", string(jsonStr))
				return nil
			}))
		}()
		return nil
	})
	return jsonFunc
}
