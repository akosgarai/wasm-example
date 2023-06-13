package dom

import (
	"encoding/json"
	"syscall/js"

	"github.com/akosgarai/wasm-example/pkg/client/request"
)

// Helper functions for the select component.

const (
	// MainWrapperClassName is the class name for the main wrapper.
	MainWrapperClassName = "select-wrapper"
	// OptionsWrapperClassName is the class name for the options wrapper.
	OptionsWrapperClassName = "select-options-wrapper hidden"
	// OptionClassName is the class name for an option.
	OptionClassName = "select-option"
)

type selectBuilder struct {
	document js.Value
	prefix   string
}

func newSelectBuilder(document js.Value, prefix string) *selectBuilder {
	return &selectBuilder{
		document: document,
		prefix:   prefix,
	}
}

func (b *selectBuilder) wrapper() js.Value {
	className := MainWrapperClassName
	if b.prefix != "" {
		className = b.prefix + "-" + MainWrapperClassName + " " + className
	}
	return Div(b.document, map[string]interface{}{
		"className": className,
	})
}
func (b *selectBuilder) optionsWrapper() js.Value {
	className := OptionsWrapperClassName
	if b.prefix != "" {
		className = b.prefix + "-" + OptionsWrapperClassName + " " + className
	}
	return Div(b.document, map[string]interface{}{
		"className": className,
	})
}

func (b *selectBuilder) hiddenInput(inputName string) js.Value {
	return Input(b.document, "hidden", map[string]interface{}{
		"name": inputName,
		"id":   inputName,
	})
}

func (b *selectBuilder) displayInput(value string, readonly bool) js.Value {
	attr := map[string]interface{}{
		"value": value,
	}
	if readonly {
		attr["readOnly"] = "readonly"
	}

	return Input(b.document, "text", attr)
}

func (b *selectBuilder) buildOptionsFromAPI(document, optionsWrapper, hiddenInput, displayInput js.Value, apiURL, selected string) {
	// gather the options from the API
	dataRaw, err := request.Get(apiURL)
	if err != nil {
		document.Get("alert").Invoke(err.Error())
		return
	}
	var resp = request.Response{}
	json.Unmarshal(dataRaw, &resp)
	b.buildOptionsFromMap(document, optionsWrapper, hiddenInput, displayInput, resp.Data, selected)
}
func (b *selectBuilder) buildOptionsFromMap(document, optionsWrapper, hiddenInput, displayInput js.Value, options map[string]string, selected string) {
	notSelectedClassName := OptionClassName
	if b.prefix != "" {
		notSelectedClassName = b.prefix + "-" + OptionClassName + " " + notSelectedClassName
	}
	if selected == "" {
		notSelectedClassName += " selected"
	}
	// Add the not selected option to the optionsWrapper
	notSelectedOption := Div(document, map[string]interface{}{
		"className": notSelectedClassName,
		"innerHTML": "-",
	})
	notSelectedOption.Get("dataset").Set("value", "")
	for _, value := range options {
		className := OptionClassName
		if b.prefix != "" {
			className = b.prefix + "-" + OptionClassName + " " + className
		}
		if value == selected {
			className += " selected"
		}
		optionElement := Div(document, map[string]interface{}{
			"className": className,
			"innerHTML": value,
		})
		optionElement.Get("dataset").Set("value", value)
		optionElement.Set("onclick", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			// Set the value of the hidden input element
			hiddenInput.Set("value", this.Get("dataset").Get("value"))
			// Set the value of the readonly text input element
			displayInput.Set("value", this.Get("innerHTML"))
			// Remove the selected class from all the options
			selectedOption := optionsWrapper.Call("querySelector", ".selected")
			if selectedOption.Truthy() {
				selectedOption.Get("classList").Call("remove", "selected")
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
}
