package indexer

import (
	"testing"
)

func TestLevenshteinDist(t *testing.T) {
	tests := []struct {
		name     string
		a, b     string
		expected int
	}{
		{
			name:     "basic",
			a:        "sitting",
			b:        "kitten",
			expected: 3,
		},
		{
			name:     "equal",
			a:        "kitten",
			b:        "kitten",
			expected: 0,
		},
		{
			name:     "empty a",
			a:        "",
			b:        "kitten",
			expected: 6,
		},
		{
			name:     "empty b",
			a:        "sitting",
			b:        "",
			expected: 7,
		},
		{
			name:     "empty a and b",
			a:        "",
			b:        "",
			expected: 0,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := levenshteinDist(test.a, test.b)
			if result != test.expected {
				t.Errorf("levenshteinDist(%s, %s) = %d, expected %d", test.a, test.b, result, test.expected)
			}
		})
	}
}
