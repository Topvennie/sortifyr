package utils

// SliceMap maps a slice of type T to a slice of type U
func SliceMap[T any, U any](input []T, mapFunc func(T) U) []U {
	result := make([]U, len(input))

	for i, item := range input {
		result[i] = mapFunc(item)
	}

	return result
}

// SliceFilter returns a new slice consisting of elements that passed the filter function
func SliceFilter[T any](input []T, filter func(T) bool) []T {
	var result []T

	for _, item := range input {
		if filter(item) {
			result = append(result, item)
		}
	}

	return result
}

// SliceFind returns a pointer to the first item determined by the equal function, nil if not found
// The second argument returns true if found, false otherwise
func SliceFind[T any](input []T, equal func(T) bool) (*T, bool) {
	for _, item := range input {
		if equal(item) {
			return &item, true
		}
	}

	return nil, false
}

// SliceUnique filters out all duplicate elements
func SliceUnique[T comparable](input []T) []T {
	result := make([]T, 0, len(input))
	seen := make(map[T]struct{}, len(input))

	for _, item := range input {
		if _, ok := seen[item]; !ok {
			seen[item] = struct{}{}
			result = append(result, item)
		}
	}

	return result
}

// SliceReference converts a []T slice to []*T
func SliceReference[T any](input []T) []*T {
	result := make([]*T, len(input))

	for i, item := range input {
		result[i] = &item
	}

	return result
}

// SliceDereference converts a []*T slice to []T
func SliceDereference[T any](input []*T) []T {
	result := make([]T, len(input))

	for i, ptr := range input {
		result[i] = *ptr
	}

	return result
}

// SliceToMap maps a slice to a map with
// key -> result of toKey function
// value -> slice entry
func SliceToMap[T any, U comparable](input []T, toKey func(T) U) map[U]T {
	result := make(map[U]T, len(input))

	for _, item := range input {
		result[toKey(item)] = item
	}

	return result
}

// SliceRepeat constructs a slice of length `count` of element `value`
func SliceRepeat[T any](value T, count int) []T {
	result := make([]T, count)
	for i := range count {
		result[i] = value
	}

	return result
}

// SliceFlatten flattens 2D slices to 1D
func SliceFlatten[T any](slice [][]T) []T {
	var result []T
	for _, s := range slice {
		result = append(result, s...)
	}

	return result
}

// SliceMerge merges multiple slices together
func SliceMerge[T any](slices ...[]T) []T {
	var result []T
	for _, slice := range slices {
		result = append(result, slice...)
	}

	return result
}

// SliceSanitize removes all null values according to T{}
func SliceSanitize[T comparable](slice []T) []T {
	sanitized := []T{}
	var zero T

	for _, item := range slice {
		if item != zero {
			sanitized = append(sanitized, item)
		}
	}

	return sanitized
}
