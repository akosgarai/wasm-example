package dom

import (
	"encoding/json"
	"syscall/js"

	"github.com/akosgarai/wasm-example/pkg/client/request"
)

// CheckBox returns a new checkbox element.
func CheckBox(document js.Value, inputName, inputID, inputLabel string, inputValue interface{}, selected bool) js.Value {
	switchWrapper := Div(document, map[string]interface{}{
		"className": "switch-wrapper",
	})
	// add the label to the switchWrapper. It has to be empty.
	label := Label(document, "", inputID)
	switchWrapper.Call("appendChild", label)
	// add the checkbox to the label
	checkbox := Input(document, "checkbox", map[string]interface{}{
		"name":    inputName,
		"id":      inputID,
		"checked": selected,
		"value":   inputValue,
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

// CheckBoxList returns a new list of checkbox list elements.
func CheckBoxList(document js.Value, listName, listLabel, listTitle, apiURL string) js.Value {
	listContainer := Div(document, map[string]interface{}{"className": "form-item center", "id": listName + "-container", "title": listTitle})
	// gather the options from the API
	dataRaw, err := request.Get(apiURL)
	if err != nil {
		js.Global().Get("alert").Invoke(err.Error())
		return listContainer
	}
	var resp = request.SelectOptionsResponse{}
	json.Unmarshal(dataRaw, &resp)
	for _, item := range resp.Data {
		itemContainer := Div(document, map[string]interface{}{"className": "form-item center"})
		// append the checkbox
		itemLabel, itemValue := item.Get()
		checkbox := CheckBox(document, listName, listName+"-"+itemLabel, itemLabel, itemValue.Get(), false)
		itemContainer.Call("appendChild", checkbox)
		// add the error message container
		errorMessageContainer := Div(document, map[string]interface{}{"className": "error-message", "id": listName + "-" + itemLabel + "-error-message"})
		itemContainer.Call("appendChild", errorMessageContainer)
		listContainer.Call("appendChild", itemContainer)
	}
	return listContainer
}
