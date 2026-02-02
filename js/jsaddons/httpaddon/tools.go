package httpaddon

import "slices"

func AddUnique(slice []string, value string) []string {
	if slices.Contains(slice, value) {
		return slice
	}
	return append(slice, value)
}

func AddUniqueValues(slice []string, values ...string) []string {
	for _, value := range values {
		slice = AddUnique(slice, value)
	}
	return slice
}
