// Package utils contains various util functions
package utils

import (
	"fmt"
)

// MapGetKeyAsType retrieves a key from a map and converts it to a given type
func MapGetKeyAsType[T any](key string, m map[string]any) (T, error) {
	if valueRaw, found := m[key]; found {
		if value, ok := valueRaw.(T); ok {
			return value, nil
		}
	}

	var zero T
	return zero, fmt.Errorf("unable to find %s key in %v", key, m)
}

// MapValues returns a slice of all values
func MapValues[T comparable, U any](input map[T]U) []U {
	result := make([]U, 0, len(input))

	for _, v := range input {
		result = append(result, v)
	}

	return result
}
