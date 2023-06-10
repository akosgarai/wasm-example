//go:build js && wasm
// +build js,wasm

package main

import (
	"fmt"
	"syscall/js"

	"github.com/akosgarai/wasm-example/pkg/formatter"
)

func jsonWrapper() js.Func {
	jsonFunc := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		if len(args) != 1 {
			return map[string]interface{}{
				"error": "Invalid no of arguments passed",
			}
		}
		jsDoc := js.Global().Get("document")
		if !jsDoc.Truthy() {
			return map[string]interface{}{
				"error": "Unable to get document object",
			}
		}
		jsonOuputTextArea := jsDoc.Call("getElementById", "jsonoutput")
		if !jsonOuputTextArea.Truthy() {
			return map[string]interface{}{
				"error": "Unable to get output text area",
			}
		}
		inputJSON := args[0].String()
		fmt.Printf("input %s\n", inputJSON)
		pretty, err := formatter.PrettyJSON(inputJSON)
		if err != nil {
			return map[string]interface{}{
				"error": fmt.Sprintf("unable to parse JSON. Error %s occurred\n", err),
			}
		}
		jsonOuputTextArea.Set("value", pretty)
		return nil
	})
	return jsonFunc
}

func main() {
	fmt.Println("WASM Go Initialized")
	js.Global().Set("formatJSON", jsonWrapper())
	<-make(chan bool)
}
