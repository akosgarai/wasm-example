package dom

import "syscall/js"

// SimpleSelect returns a simple select element.
// The document input parameter is the document.
// The next input parameter is the map of options.
func SimpleSelect(document js.Value, options map[string]string, inputName, selected string) js.Value {
	builder := newSelectBuilder(document, "simple")
	selectorWrapper := builder.wrapper()
	// Add the hidden input element to the selectorWrapper
	hiddenInput := builder.hiddenInput(inputName)
	selectorWrapper.Call("appendChild", hiddenInput)
	// The selected is displayed in a readonly text input and next to it is a button to open the select options.
	// Add the readonly text input to the selectorWrapper
	readonlyInput := builder.displayInput(options[selected], true)
	selectorWrapper.Call("appendChild", readonlyInput)
	// Add the button to the selectorWrapper
	button := Button(document, "Select")
	optionsWrapper := builder.optionsWrapper()
	button.Set("onclick", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		// Toggle the optionsWrapper
		selectorWrapper.Get("classList").Call("toggle", OpenClassName)
		// Toggle the overlay
		overlay := document.Call("getElementById", "overlay")
		overlay.Get("classList").Call("toggle", "hidden")
		// add click event listener to the overlay to hide the optionsWrapper
		overlay.Set("onclick", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			// Hide the optionsWrapper
			selectorWrapper.Get("classList").Call("remove", OpenClassName)
			// Hide the overlay
			overlay.Get("classList").Call("add", "hidden")
			return nil
		}))
		return nil
	}))
	selectorWrapper.Call("appendChild", button)
	// Add the select options to the selectorWrapper
	builder.buildOptionsFromMap(document, optionsWrapper, hiddenInput, readonlyInput, options, selected)
	selectorWrapper.Call("appendChild", optionsWrapper)
	return selectorWrapper
}
