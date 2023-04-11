// Package sliceutil provides utility functions for working with slices.
package sliceutil

// ContainsString returns true if the given string is in the given slice.
func ContainsString(slice []string, s string) bool {
	for _, v := range slice {
		if v == s {
			return true
		}
	}

	return false
}

// ContainsInt returns true if the given int is in the given slice.
func ContainsInt(slice []int, i int) bool {
	for _, v := range slice {
		if v == i {
			return true
		}
	}

	return false
}
