package dom

import "syscall/js"

// SimpleSelect returns a simple select element.
// The document input parameter is the document.
// The next input parameter is the map of options.
func SimpleSelect(document js.Value, options map[string]string, inputName, selected string) js.Value {
	selectorWrapper := CreateElement(document, "div", map[string]interface{}{
		"className": "simple-select-wrapper",
	})
	// Add the hidden input element to the selectorWrapper
	hiddenInput := CreateElement(document, "input", map[string]interface{}{
		"type": "hidden",
		"name": inputName,
		"id":   inputName,
	})
	selectorWrapper.Call("appendChild", hiddenInput)
	// The selected is displayed in a readonly text input and next to it is a button to open the select options.
	// Add the readonly text input to the selectorWrapper
	readonlyInput := CreateElement(document, "input", map[string]interface{}{
		"type":     "text",
		"readOnly": "readonly",
		"value":    options[selected],
	})
	selectorWrapper.Call("appendChild", readonlyInput)
	// Add the button to the selectorWrapper
	button := CreateElement(document, "button", map[string]interface{}{
		"type":      "button",
		"innerHTML": "Select",
	})
	optionsWrapper := CreateElement(document, "div", map[string]interface{}{
		"className": "simple-select-options-wrapper hidden",
	})
	button.Set("onclick", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		// Toggle the optionsWrapper
		optionsWrapper.Get("classList").Call("toggle", "hidden")
		return nil
	}))
	selectorWrapper.Call("appendChild", button)
	// Add the select options to the selectorWrapper
	for key, value := range options {
		optionElement := CreateElement(document, "div", map[string]interface{}{
			"className": "simple-select-option",
		})
		optionElement.Set("innerHTML", value)
		optionElement.Set("dataset", map[string]interface{}{
			"value": key,
		})
		optionsWrapper.Call("appendChild", optionElement)
	}
	selectorWrapper.Call("appendChild", optionsWrapper)
	return selectorWrapper
}
