package request

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

// Response represents the response of the request
type Response struct {
	Data map[string]string
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
