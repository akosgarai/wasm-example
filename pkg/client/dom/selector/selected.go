package selector

import (
	"errors"
	"fmt"
)

const (
	// StringValueType is the string value type.
	StringValueType = "string"
	// NumberValueType is the number value type.
	NumberValueType = "number"
)

var (
	// ErrInvalidValueType is returned when the value type is not string or number.
	ErrInvalidValueType = errors.New("invalid value type")
)

// Selected is a struct to store the selected value.
type Selected struct {
	// The selected value.
	// can be a string or a number
	value interface{}
	// the type of the value (string or number)
	valueType string
}

// Get returns the value
func (s *Selected) Get() interface{} {
	return s.value
}

// IsSelected returns true if the value is selected.
func (s *Selected) IsSelected(value interface{}) bool {
	switch s.valueType {
	case StringValueType:
		return s.value.(string) == value.(string)
	case NumberValueType:
		return s.value.(int) == value.(int)
	}
	return false
}

// IsEmpty returns true if the value is empty.
func (s *Selected) IsEmpty() bool {
	// Empty string or 0
	switch s.valueType {
	case StringValueType:
		return s.value.(string) == ""
	case NumberValueType:
		return s.value.(int) == 0
	}
	return false
}

// DisplayValue returns the display value.
func (s *Selected) DisplayValue() string {
	switch s.valueType {
	case StringValueType:
		return s.value.(string)
	case NumberValueType:
		return fmt.Sprintf("%d", (s.value.(int)))
	}
	return ""
}

// NewSelected creates a new Selected.
func NewSelected(value interface{}) (*Selected, error) {
	// Allow only string or number
	if value != nil {
		switch value.(type) {
		case string:
			return newSelectedString(value.(string)), nil
		case int:
			return newSelectedInt(value.(int)), nil
		case float64:
			// cast to int
			return newSelectedInt(int(value.(float64))), nil
		}
	}
	return nil, ErrInvalidValueType
}

func newSelectedString(value string) *Selected {
	return &Selected{
		value:     value,
		valueType: StringValueType,
	}
}

func newSelectedInt(value int) *Selected {
	return &Selected{
		value:     value,
		valueType: NumberValueType,
	}
}
