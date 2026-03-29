package indexer

import (
	"fmt"
	"testing"
)

func TestStem(t *testing.T) {
	tests := []struct {
		name     string
		word     string
		expected string
	}{
		{
			name: "basic",
			word: "Accumulations",
			expected: "ACCUMUL",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, err := stem(test.word)
			if err != nil {
				fmt.Printf("ster(%s) returns an err - %v\n", test.word, err)
			}
			if result != test.expected {
				t.Errorf("stem(%s) = %s, expected %s\n", test.word, result, test.expected)
			}
		})
	}
}
