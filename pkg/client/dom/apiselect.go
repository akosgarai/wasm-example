package dom

import (
	"syscall/js"

	"github.com/akosgarai/wasm-example/pkg/client/dom/selector"
)

// APISelect returns a new select element.
// The options are gathered from the API.
// The input field on the select element is a search field.
// Once it changes, the options are filtered from the API.
// It allows the user to select an option from the filtered options.
// Also possible to add a new option.
func APISelect(document js.Value, apiURL, inputName string, selected *selector.Selected) js.Value {
	builder := newSelectBuilder(document, "api")
	selectorWrapper := builder.wrapper()
	// Add the hidden input element to the selectorWrapper
	hiddenInput := builder.hiddenInput(inputName)
	selectorWrapper.Call("appendChild", hiddenInput)
	displayInput := builder.displayInput(selected.DisplayValue(), false)
	selectorWrapper.Call("appendChild", displayInput)
	optionsWrapper := builder.optionsWrapper()
	displayInput.Set("onclick", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
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
	displayInput.Set("oninput", displayInputChanged(document, optionsWrapper, hiddenInput, displayInput, apiURL, builder))
	builder.buildOptionsFromAPI(document, optionsWrapper, hiddenInput, displayInput, apiURL, selected)
	selectorWrapper.Call("appendChild", optionsWrapper)

	return selectorWrapper
}

// displayInputChanged is the event handler for the displayInput.
// It rebuilds the optionsWrapper.
func displayInputChanged(document, optionsWrapper, hiddenInput, displayInput js.Value, apiURL string, builder *selectBuilder) js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		// Get the value of the displayInput
		value := displayInput.Get("value").String()
		selected, _ := selector.NewSelected(value)
		optionsWrapper.Set("innerHTML", "")
		// set the value of the hidden input
		hiddenInput.Set("value", value)
		// Build the options from the API
		go builder.buildOptionsFromSearchAPI(document, optionsWrapper, hiddenInput, displayInput, apiURL, selected)
		return nil
	})
}
