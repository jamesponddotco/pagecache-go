// Package sliceutil provides utility functions for working with slices.
package sliceutil

// MatchString returns true if the given string is in the given slice.
func MatchString(slice []string, s string) bool {
	for _, v := range slice {
		if v == s {
			return true
		}
	}

	return false
}

// MatchInt returns true if the given int is in the given slice.
func MatchInt(slice []int, i int) bool {
	for _, v := range slice {
		if v == i {
			return true
		}
	}

	return false
}
