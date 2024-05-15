package main

import (
	"errors"
)

// Filter returns applies the filterFunction function to every element in a slice and returns a new slice with elements for which it returns true
func Filter[T any](tab []T, filterFunction func(todo T) bool) (res []T) {
	res = []T{}
	for _, item := range tab {
		if filterFunction(item) {
			res = append(res, item)
		}
	}
	return
}

// updateFromMap gets the element corresponding to fieldName in the map parameter and updates field with if it exits and the types correspond
func updateFromMap[T any](field *T, m map[string]any, fieldName string) error {
	if value, ok := m[fieldName]; ok {
		if _, correctType := value.(T); !correctType {
			return errors.New("Parameter " + fieldName + " is of wrong type")
		}
		*field = value.(T)
	}
	return nil
}
