package main

import (
	"errors"
	"fmt"
)

// Filter applies the filterFunction function to every element in a slice and returns a new slice with elements for which it returns true
func Filter[T any](tab *[]T, filterFunction func(todo T) bool) []T {
	res := []T{}
	for _, item := range *tab {
		if filterFunction(item) {
			res = append(res, item)
		}
	}
	return res
}

// updateFromMap gets the element corresponding to fieldName in the map parameter and updates field with if it exits and the types correspond
func updateFromMap[T any, MapKey comparable](field *T, m map[MapKey]any, fieldName MapKey) error {
	if value, ok := m[fieldName]; ok {
		if _, correctType := value.(T); !correctType {
			return errors.New(fmt.Sprintf("Parameter %v is of wrong type", fieldName))
		}
		*field = value.(T)
	}
	return nil
}
