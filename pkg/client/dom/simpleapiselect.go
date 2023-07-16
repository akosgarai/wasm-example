package dom

import (
	"encoding/json"
	"syscall/js"

	"github.com/akosgarai/wasm-example/pkg/client/request"
)

// SimpleAPISelect returns a simple select element.
// The document input parameter is the document.
// The next input parameter is the map of options.
func SimpleAPISelect(document js.Value, apiURL, inputName, selected string) js.Value {
	// gather the options from the API
	dataRaw, err := request.Get(apiURL)
	if err != nil {
		document.Get("alert").Invoke(err.Error())
		emptyMap := make(map[string]string)
		return SimpleSelect(document, emptyMap, inputName, selected)
	}
	var resp = request.Response{}
	json.Unmarshal(dataRaw, &resp)
	return SimpleSelect(document, resp.Data, inputName, selected)
}
