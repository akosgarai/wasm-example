package selector

import "fmt"

// SelectOption represents a select option
type SelectOption struct {
	Name string `json:"name"`
	// The id could be number or string, so we use interface{} here.
	ID interface{} `json:"id"`
}

// Get returns the name of the option as string and the is as Selected
func (s *SelectOption) Get() (string, *Selected) {
	selected, err := NewSelected(s.ID)
	if err != nil {
		panic(fmt.Sprintf("invalid id: %v, %v", s.ID, err))
	}
	return s.Name, selected
}

// SelectOptions is a slice of SelectOption
type SelectOptions []SelectOption
