package fileingestion

import (
	"os"
	"strings"
	"testing"
)

func TestExtractText(t *testing.T) {
	// Test case: Unsupported file type
	t.Run("Unsupported file type", func(t *testing.T) {
		// Create a temporary UNSUPPORTED file.
		content := "dummy"
		tmpFile, err := os.CreateTemp("", "test-*.md")
		if err != nil {
			t.Fatal(err)
		}
		defer os.Remove(tmpFile.Name())
		os.WriteFile(tmpFile.Name(), []byte(content), 0644)

		_, err = ExtractText(tmpFile.Name())
		if err == nil {
			t.Error("Expected error for UNSUPPORTED file type, got nil")
		}
	})

	// Test case: Valid ".txt" extraction
	t.Run("Valid \".txt\" text extraction", func(t *testing.T) {
		// Create a temporary TXT file.
		content := "Hello, this is a test file."
		tmpFile, err := os.CreateTemp("", "test-*.txt")
		if err != nil {
			t.Fatal(err)
		}
		defer os.Remove(tmpFile.Name())
		os.WriteFile(tmpFile.Name(), []byte(content), 0644)

		text, err := ExtractText(tmpFile.Name())
		if err != nil {
			t.Fatalf("ExtractText returned error: %v", err)
		}
		contentCheck := "Hello, this is a test file."
		if strings.TrimSpace(text) != contentCheck {
			t.Errorf("Expected |%s|, got |%s|", contentCheck, text)
		}
	})

	// Test case: Valid ".csv" extraction
	t.Run("Valid \".csv\" extraction", func(t *testing.T) {
		// Create a temporary CSV file
		csvContent := "name,age\nFilipe,37\nDani,32"
		tmpFile, err := os.CreateTemp("", "test-*.csv")
		if err != nil {
			t.Fatal(err)
		}
		defer os.Remove(tmpFile.Name())
		os.WriteFile(tmpFile.Name(), []byte(csvContent), 0644)

		text, err := ExtractText(tmpFile.Name())
		if err != nil {
			t.Fatalf("ExtractText returned error: %v", err)
		}
		contentCheck := "Dani"
		if !strings.Contains(text, contentCheck) {
			t.Errorf("Expected extracted text to contain |%s|, extracted text is |%s|", contentCheck, text)
		}
	})
}
