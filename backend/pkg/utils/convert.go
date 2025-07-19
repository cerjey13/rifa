package utils

import "strconv"

// ParseFloatOrZero safely parses a string to float64, returning 0.0 on error.
func ParseFloatOrZero(s string) float64 {
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0.0
	}
	return f
}
