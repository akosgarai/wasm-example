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

// buildLayout builds the layout of the page.
func (p *Instance) buildLayout() {
	body := p.document.Call("querySelector", "body")
	// Contaner div with class container
	containerDiv := dom.Div(p.document, map[string]interface{}{
		"className": ContainerClassName,
	})
	body.Call("appendChild", containerDiv)
	// Header div with class header
	headerDiv := dom.Div(p.document, map[string]interface{}{
		"className": HeaderClassName,
	})
	containerDiv.Call("appendChild", headerDiv)
	// Header div constans a h1 with the title
	headerDiv.Call("appendChild", dom.H1(p.document, p.Title))
	contentDiv := dom.Div(p.document, map[string]interface{}{
		"className": ContentClassName,
	})
	containerDiv.Call("appendChild", contentDiv)
	// overlay div
	overlayDiv := dom.Div(p.document, map[string]interface{}{
		"id":        "overlay",
		"className": "overlay hidden",
	})
	body.Call("appendChild", overlayDiv)
}
