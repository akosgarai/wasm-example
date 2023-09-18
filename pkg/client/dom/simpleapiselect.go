package dom

import (
	"encoding/json"
	"syscall/js"

	"github.com/akosgarai/wasm-example/pkg/client/dom/selector"
	"github.com/akosgarai/wasm-example/pkg/client/request"
)

// SimpleAPISelect returns a simple select element.
// The document input parameter is the document.
// The next input parameter is the map of options.
func SimpleAPISelect(document js.Value, apiURL, inputName string, selected *selector.Selected) js.Value {
	// gather the options from the API
	dataRaw, err := request.Get(apiURL)
	if err != nil {
		js.Global().Get("alert").Invoke(err.Error())
		var emptyList selector.SelectOptions
		return SimpleSelect(document, emptyList, inputName, selected)
	}
	var resp = request.SelectOptionsResponse{}
	json.Unmarshal(dataRaw, &resp)
	return SimpleSelect(document, resp.Data, inputName, selected)
}
