package page

import (
	"syscall/js"
)

// Page is an interface that represents a page
type Page interface {
	LoadPage()
	Run()
}

// Instance is a struct that implements the Page interface
type Instance struct {
	Title    string
	document js.Value
}

// NewPage returns a new Instance
func NewPage(title string) *Instance {
	jsDoc := js.Global().Get("document")
	if !jsDoc.Truthy() {
		panic("document is not available")
	}
	return &Instance{
		Title:    title,
		document: jsDoc,
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
}

// Run runs the page
func (p *Instance) Run() {
}

// Document returns the document
func (p *Instance) Document() js.Value {
	return p.document
}

// CreateElement returns a new element.
func (p *Instance) CreateElement(tagName string, attrs map[string]interface{}) js.Value {
	element := p.Document().Call("createElement", tagName)
	for key, value := range attrs {
		element.Set(key, value)
	}
	return element
}
