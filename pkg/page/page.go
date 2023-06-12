package page

import (
	"syscall/js"

	"github.com/akosgarai/wasm-example/pkg/client/dom"
)

const (
	// ContainerClassName is the class name of the container
	ContainerClassName = "container"
	// HeaderClassName is the class name of the header
	HeaderClassName = "header"
	// ContentClassName is the class name of the content
	ContentClassName = "content"
)

// Page is an interface that represents a page
type Page interface {
	LoadPage()
	Run()
}

// Instance is a struct that implements the Page interface
type Instance struct {
	Title     string
	document  js.Value
	webSocket js.Value
	alert     js.Value
	json      js.Value
}

// NewPage returns a new Instance
func NewPage(title string) *Instance {
	jsDoc := js.Global().Get("document")
	webSocket := js.Global().Get("WebSocket")
	alert := js.Global().Get("alert")
	json := js.Global().Get("JSON")
	if !jsDoc.Truthy() {
		panic("document is not available")
	}
	return &Instance{
		Title:     title,
		document:  jsDoc,
		webSocket: webSocket,
		alert:     alert,
		json:      json,
	}
}

// SetPageTitle sets the page title
func (p *Instance) setPageTitle() {
	title := p.document.Call("querySelector", "title")
	if !title.Truthy() {
		return
	}
	title.Set("innerHTML", p.Title)
}

// GetElementByID returns the element by the given id
func (p *Instance) GetElementByID(id string) js.Value {
	return p.document.Call("getElementById", id)
}

// LoadPage loads the page
func (p *Instance) LoadPage() {
	p.setPageTitle()
	p.buildLayout()
}

// Run runs the page
func (p *Instance) Run() {
}

// Document returns the document
func (p *Instance) Document() js.Value {
	return p.document
}

// WebSocket returns the websocket
func (p *Instance) WebSocket() js.Value {
	return p.webSocket
}

// Alert returns the alert
func (p *Instance) Alert() js.Value {
	return p.alert
}

// JSON returns the JSON
func (p *Instance) JSON() js.Value {
	return p.json
}

// CreateElement returns a new element.
func (p *Instance) CreateElement(tagName string, attrs map[string]interface{}) js.Value {
	return dom.CreateElement(p.document, tagName, attrs)
}

// CreateSelectElement returns a new dynamic select element.
// The first input parameter is the container element.
// The second input parameter is the map of the options.
// The third input parameter is the selected value.
func (p *Instance) CreateSelectElement(container js.Value, options map[string]string, selected string) js.Value {
	selectElement := p.CreateElement("select", map[string]interface{}{
		"className": "form-control",
	})
	for key, value := range options {
		optionElement := p.CreateElement("option", map[string]interface{}{
			"value": key,
		})
		if key == selected {
			optionElement.Set("selected", true)
		}
		optionElement.Set("innerHTML", value)
		selectElement.Call("appendChild", optionElement)
	}
	container.Call("appendChild", selectElement)
	return selectElement
}

// buildLayout builds the layout of the page.
func (p *Instance) buildLayout() {
	body := p.document.Call("querySelector", "body")
	// Contaner div with class container
	containerDiv := p.CreateElement("div", map[string]interface{}{
		"className": ContainerClassName,
	})
	body.Call("appendChild", containerDiv)
	// Header div with class header
	headerDiv := p.CreateElement("div", map[string]interface{}{
		"className": HeaderClassName,
	})
	containerDiv.Call("appendChild", headerDiv)
	// Header div constans a h1 with the title
	headerDiv.Call("appendChild", p.CreateElement("h1", map[string]interface{}{
		"innerHTML": p.Title,
	}))
	contentDiv := p.CreateElement("div", map[string]interface{}{
		"className": ContentClassName,
	})
	containerDiv.Call("appendChild", contentDiv)
}
