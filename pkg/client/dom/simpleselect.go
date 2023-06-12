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
		// Toggle the overlay
		overlay := document.Call("getElementById", "overlay")
		overlay.Get("classList").Call("toggle", "hidden")
		// add click event listener to the overlay to hide the optionsWrapper
		overlay.Set("onclick", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			// Hide the optionsWrapper
			optionsWrapper.Get("classList").Call("add", "hidden")
			// Hide the overlay
			overlay.Get("classList").Call("add", "hidden")
			return nil
		}))
		return nil
	}))
	selectorWrapper.Call("appendChild", button)
	// Add the select options to the selectorWrapper
	for key, value := range options {
		className := "simple-select-option"
		if key == selected {
			className += " selected"
		}
		optionElement := CreateElement(document, "div", map[string]interface{}{
			"className": className,
		})
		optionElement.Set("innerHTML", value)
		optionElement.Get("dataset").Set("value", key)
		optionElement.Set("onclick", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			// Set the value of the hidden input element
			hiddenInput.Set("value", this.Get("dataset").Get("value"))
			// Set the value of the readonly text input element
			readonlyInput.Set("value", this.Get("innerHTML"))
			// Remove the selected class from all the options
			selected := selectorWrapper.Call("querySelector", ".selected")
			if selected.Truthy() {
				selected.Get("classList").Call("remove", "selected")
			}
			// Add the selected class to the clicked option
			this.Get("classList").Call("add", "selected")
			// Hide the optionsWrapper
			optionsWrapper.Get("classList").Call("add", "hidden")
			// Hide the overlay
			overlay := document.Call("getElementById", "overlay")
			overlay.Get("classList").Call("add", "hidden")
			return nil
		}))
		optionsWrapper.Call("appendChild", optionElement)
	}
	selectorWrapper.Call("appendChild", optionsWrapper)
	return selectorWrapper
}
