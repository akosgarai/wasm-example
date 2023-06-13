package dom

import (
	"syscall/js"
)

// APISelect returns a new select element.
// The options are gathered from the API.
// The input field on the select element is a search field.
// Once it changes, the options are filtered from the API.
// It allows the user to select an option from the filtered options.
// Also possible to add a new option.
func APISelect(document js.Value, apiURL, inputName, selected string) js.Value {
	builder := newSelectBuilder(document, "api")
	selectorWrapper := builder.wrapper()
	// Add the hidden input element to the selectorWrapper
	hiddenInput := builder.hiddenInput(inputName)
	selectorWrapper.Call("appendChild", hiddenInput)
	displayInput := builder.displayInput(selected, false)
	selectorWrapper.Call("appendChild", displayInput)
	optionsWrapper := builder.optionsWrapper()
	displayInput.Set("onclick", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
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
	builder.buildOptionsFromAPI(document, optionsWrapper, hiddenInput, displayInput, apiURL, selected)
	selectorWrapper.Call("appendChild", optionsWrapper)

	return selectorWrapper
}
