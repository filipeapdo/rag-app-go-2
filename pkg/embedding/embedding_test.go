package embedding

import (
	"testing"
)

func TestChunkText(t *testing.T) {
	// Test cases:
	tests := []struct {
		name        string
		input       string
		chunkSize   int
		overlapSize int
		expected    []string
	}{
		{
			name:        "Standard chunking with overlap",
			input:       "The quick brown fox jumps over the lazy dog near the riverbank.",
			chunkSize:   5,
			overlapSize: 2,
			expected: []string{
				"The quick brown fox jumps",
				"fox jumps over the lazy",
				"the lazy dog near the",
				"near the riverbank.",
			},
		},
		{
			name:        "Exact chunk size, no overlap",
			input:       "One two three four five six seven eight nine ten",
			chunkSize:   5,
			overlapSize: 0,
			expected: []string{
				"One two three four five",
				"six seven eight nine ten",
			},
		},
		{
			name:        "Short text (less than chunk size)",
			input:       "Hello world!",
			chunkSize:   5,
			overlapSize: 2,
			expected:    []string{"Hello world!"},
		},
		{
			name:        "Empty text",
			input:       "",
			chunkSize:   5,
			overlapSize: 2,
			expected:    []string{},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := GetChunks(tc.input, tc.chunkSize, tc.overlapSize)
			if len(result) != len(tc.expected) {
				t.Errorf("Expected %d chunks, got %d", len(tc.expected), len(result))
			}
			for i, chunk := range result {
				if chunk != tc.expected[i] {
					t.Errorf("Mismatch at index %d: expected %q, got %q", i, tc.expected[i], chunk)
				}
			}
		})
	}
}
