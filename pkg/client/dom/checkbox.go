package dom

import "syscall/js"

// CheckBox returns a new checkbox element.
func CheckBox(document js.Value, inputName, inputLabel string, selected bool) js.Value {
	switchWrapper := Div(document, map[string]interface{}{
		"className": "switch-wrapper",
	})
	// add the label to the switchWrapper. It has to be empty.
	label := Label(document, "", inputName)
	switchWrapper.Call("appendChild", label)
	// add the checkbox to the label
	checkbox := Input(document, "checkbox", map[string]interface{}{
		"name":    inputName,
		"id":      inputName,
		"checked": selected,
	})
	label.Call("appendChild", checkbox)
	// add the slider to the label
	slider := Span(document, "", "switch")
	label.Call("appendChild", slider)
	// Add the text node to the switchWrapper
	textNode := Label(document, inputLabel, "")
	switchWrapper.Call("appendChild", textNode)
	return switchWrapper
}
