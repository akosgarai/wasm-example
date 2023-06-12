package dom

import "syscall/js"

// CreateElement returns a new element.
func CreateElement(document js.Value, tagName string, attrs map[string]interface{}) js.Value {
	element := document.Call("createElement", tagName)
	for key, value := range attrs {
		element.Set(key, value)
	}
	return element
}
