package utils

import (
	"fmt"
	"strconv"
)

// ParseFloatOrZero safely parses a string to float64, returning 0.0 on error.
func ParseFloatOrZero(s string) float64 {
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0.0
	}
	return f
}

func ConvertToIntSlice(strs []string) ([]int, error) {
	ints := make([]int, 0, len(strs))
	for _, s := range strs {
		n, err := strconv.Atoi(s)
		if err != nil {
			return nil, fmt.Errorf("invalid number %q: %w", s, err)
		}
		ints = append(ints, n)
	}
	return ints, nil
}

func ConvertToStrSlice(ints []int) []string {
	strs := make([]string, len(ints))
	for i, n := range ints {
		strs[i] = strconv.Itoa(n)
	}
	return strs
}
