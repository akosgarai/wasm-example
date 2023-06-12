package dom

import "syscall/js"

// CreateElement returns a new element.
func CreateElement(document js.Value, tagName string, attrs map[string]interface{}) js.Value {
	element := document.Call("createElement", tagName)
	for key, value := range attrs {
		element.Set(key, value)
	}
	return element
}

// Div returns a new div element.
func Div(document js.Value, attributes map[string]interface{}) js.Value {
	return CreateElement(document, "div", attributes)
}

// Input returns a new input element.
func Input(document js.Value, inputType string, attributes map[string]interface{}) js.Value {
	attributes["type"] = inputType
	return CreateElement(document, "input", attributes)
}

// Button returns a new button element.
func Button(document js.Value, label string) js.Value {
	return CreateElement(document, "button", map[string]interface{}{
		"type":      "button",
		"innerHTML": label,
	})
}

// Label returns a new label element.
func Label(document js.Value, label, labelFor string) js.Value {
	return CreateElement(document, "label", map[string]interface{}{
		"innerHTML": label,
		"htmlFor":   labelFor,
	})
}

// P returns a new p element.
func P(document js.Value, text string) js.Value {
	return CreateElement(document, "p", map[string]interface{}{
		"innerHTML": text,
	})
}

// Form returns a new form element.
func Form(document js.Value, attributes map[string]interface{}) js.Value {
	return CreateElement(document, "form", attributes)
}

// H1 returns a new h1 element.
func H1(document js.Value, text string) js.Value {
	return CreateElement(document, "h1", map[string]interface{}{
		"innerHTML": text,
	})
}
