package terminal

import (
	"testing"
)

func TestHumanizeSize(t *testing.T) {
	tests := []struct {
		size     *int64
		expected string
	}{
		{nil, ""},
		{new(int64), "0 B"},
		{int64Ptr(1023), "1023 B"},
		{int64Ptr(1024), "1 KiB"},
		{int64Ptr(1048576), "1 MiB"},
	}

	for _, test := range tests {
		result := humanizeSize(test.size)
		if result != test.expected {
			t.Errorf("Expected %s, got %s", test.expected, result)
		}
	}
}

func int64Ptr(i int64) *int64 {
	return &i
}
