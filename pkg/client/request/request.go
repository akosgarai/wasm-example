package request

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/akosgarai/wasm-example/pkg/client/dom/selector"
)

// Response represents the response of the request
type Response struct {
	Data interface{} `json:"data"`
}

// SelectOptionsResponse represents the response of the request
// that returns options for select element
type SelectOptionsResponse struct {
	Data selector.SelectOptions `json:"data"`
}

// Get returns the body of the response and error if any
func Get(path string) ([]byte, error) {
	resp, err := http.Get(path)
	if err != nil {
		return []byte(""), err
	}
	// if the status code is not 2xx, we return an error
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return []byte(""), fmt.Errorf("Status code is not 2xx. Status code: %d", resp.StatusCode)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	return body, err
}

// Post returns the body of the response and error if any
func Post(path string, body interface{}) ([]byte, error) {
	marshal, err := json.Marshal(body)
	if err != nil {
		return []byte(""), err
	}
	resp, err := http.Post(path, "application/json", bytes.NewReader(marshal))
	if err != nil {
		return []byte(""), err
	}
	// if the status code is not 2xx, we return an error
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return []byte(""), fmt.Errorf("Status code is not 2xx. Status code: %d", resp.StatusCode)
	}
	defer resp.Body.Close()
	bodyResp, err := ioutil.ReadAll(resp.Body)
	return bodyResp, err
}
