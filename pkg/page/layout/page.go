package layout

import (
	"syscall/js"

	"github.com/akosgarai/wasm-example/pkg/page"
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
		"id":        "form-container",
		"className": "row",
	})
	container.Call("appendChild", formContainer)
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
	}
	for _, item := range formItems {
		formContainer.Call("appendChild", l.buildFormItem(item.Tag, item.Attributes))
	}
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
