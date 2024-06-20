package cron

import (
	"fmt"
	"testing"
)

func TestParsePartValid(t *testing.T) {
	tests := []struct {
		part   string
		min    int
		max    int
		expect []string
		err    error
	}{
		{"*", 0, 59, genStringRangeWithStep(0, 59, 1), nil},
		{"10", 0, 59, []string{"10"}, nil},
		{"10,20", 0, 59, []string{"10", "20"}, nil},
		{"10/2", 0, 59, genStringRangeWithStep(10, 58, 2), nil},
		{"1-5", 1, 31, []string{"1", "2", "3", "4", "5"}, nil},
		{"15-20/2", 1, 31, []string{"15", "17", "19"}, nil},
	}

	for _, tc := range tests {
		t.Run(tc.part, func(t *testing.T) {
			result, err := parseField(tc.part, tc.min, tc.max)
			if err != nil {
				if tc.err == nil {
					t.Errorf("Unexpected error: %v", err)
				}
			} else if tc.err != nil {
				t.Errorf("Expected error, got none")
			} else {
				if !compareSlices(result, tc.expect) {
					t.Errorf("Expected %v, got %v", tc.expect, result)
				}
			}
		})
	}
}

func TestParsePartInvalid(t *testing.T) {
	tests := []struct {
		part string
		min  int
		max  int
		err  string
	}{
		{"abc", 0, 59, "invalid expression format"},
		{"10,abc", 0, 59, "invalid expression format"},
		{"10/abc", 0, 59, "invalid expression format"},
		{"60", 0, 59, "invalid expression format"},
		{"32", 1, 31, "invalid expression format"},
		{"10-15/0", 0, 59, "invalid expression format"},
		{"15-10/2", 1, 31, "invalid expression format"},
		{"*-*", 1, 31, "invalid expression format"},
	}

	for _, tc := range tests {
		t.Run(tc.part, func(t *testing.T) {
			_, err := parseField(tc.part, tc.min, tc.max)
			if err == nil {
				t.Errorf("Expected error, got none")
			}
		})
		fmt.Printf("test")
	}
}

func compareSlices(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
