//go:build js && wasm
// +build js,wasm

package main

import (
	"fmt"

	"github.com/akosgarai/wasm-example/pkg/application/client"
	"github.com/akosgarai/wasm-example/pkg/page/formatter"
)

func main() {
	fmt.Println("WASM Go Initialized")
	clientApp := client.New(formatter.New("JSON Formatter"))
	clientApp.Run()
}
