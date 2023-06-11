package formatter

import (
	"fmt"
	"syscall/js"

	"github.com/akosgarai/wasm-example/pkg/formatter"
	"github.com/akosgarai/wasm-example/pkg/page"
)

// Formatter is the page for the json formatter.
type Formatter struct {
	*page.Instance

	inputTextArea  js.Value
	outputTextArea js.Value
}

// New returns a new Formatter
func New(title string) *Formatter {
	return &Formatter{
		Instance: page.NewPage(title),
	}
}

// jsonWrapper returns the wrapper function for the JSON formatter.
func (f *Formatter) jsonWrapper() js.Func {
	jsonFunc := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		inputJSON := f.inputTextArea.Get("value").String()
		fmt.Printf("input %s\n", inputJSON)
		pretty, err := formatter.PrettyJSON(inputJSON)
		if err != nil {
			js.Global().Get("alert").Invoke(fmt.Sprintf("unable to parse JSON. Error %s occurred\n", err))
			return nil
		}
		f.outputTextArea.Set("value", pretty)
		return nil
	})
	return jsonFunc
}

// Run runs the formatter page.
func (f *Formatter) Run() {
	<-make(chan bool)
}

// LoadPage loads the formatter page.
func (f *Formatter) LoadPage() {
	f.Instance.LoadPage()
	// Append the input text area to the body
	// <textarea id="jsoninput" name="jsoninput" cols="80" rows="20"></textarea>
	inputAttrs := map[string]interface{}{
		"id":   "jsoninput",
		"name": "jsoninput",
		"cols": 80,
		"rows": 20,
	}
	f.inputTextArea = f.CreateElement("textarea", inputAttrs)
	body := f.Document().Call("querySelector", "body")
	body.Call("appendChild", f.inputTextArea)
	// Append the submit button to the body
	// <input id="button" type="submit" name="button" value="pretty json" onclick="json(jsoninput.value)"/>
	buttonAttrs := map[string]interface{}{
		"id":    "button",
		"type":  "submit",
		"name":  "button",
		"value": "pretty json",
	}
	jsonSubmitButton := f.CreateElement("input", buttonAttrs)
	jsonSubmitButton.Set("onclick", f.jsonWrapper().Call("bind", jsonSubmitButton))
	body.Call("appendChild", jsonSubmitButton)
	// Append the output text area to the body
	// <textarea id="jsonoutput" name="jsonoutput" cols="80" rows="20"></textarea>
	outputAttrs := map[string]interface{}{
		"id":   "jsonoutput",
		"name": "jsonoutput",
		"cols": 80,
		"rows": 20,
	}
	f.outputTextArea = f.CreateElement("textarea", outputAttrs)
	body.Call("appendChild", f.outputTextArea)
}
