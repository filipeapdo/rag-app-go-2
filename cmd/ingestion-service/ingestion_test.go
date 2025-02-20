package main

import (
	"os"
	"testing"

	"github.com/filipeapdo/rag-app-go/pkg/embedding"
	"github.com/filipeapdo/rag-app-go/pkg/fileingestion"
)

// Test the full ingestion pipeline, ensuring correct text processing and embedding
func TestIngestionPipeline(t *testing.T) {
	// Step 1: Create a temporary text file
	tmpFile, err := os.CreateTemp("", "test-*.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())

	content := "The quick brown fox jumps over the lazy dog near the riverbank."
	os.WriteFile(tmpFile.Name(), []byte(content), 0644)

	// Step 2: Extract text
	extractedText, err := fileingestion.ExtractText(tmpFile.Name())
	if err != nil {
		t.Fatalf("ExtractText failed: %v", err)
	}
	if extractedText != content {
		t.Errorf("Expected extracted text to match input, got: %s", extractedText)
	}

	// Step 3: Chunk text (only once, avoiding duplication)
	chunks := embedding.GetChunks(extractedText, 5, 2) // Now we explicitly call it only here

	expectedChunks := []string{
		"The quick brown fox jumps",
		"fox jumps over the lazy",
		"the lazy dog near the",
		"near the riverbank.",
	}

	if len(chunks) != len(expectedChunks) {
		t.Fatalf("Expected %d chunks, got %d", len(expectedChunks), len(chunks))
	}

	for i, chunk := range chunks {
		if chunk != expectedChunks[i] {
			t.Errorf("Chunk mismatch at %d: expected %q, got %q", i, expectedChunks[i], chunk)
		}
	}

	// Step 4: Get embeddings from Ollama for each chunk
	for _, chunk := range chunks {
		embeddingVector, err := embedding.GetEmbeddings(chunk) // No chunking inside this function now!
		if err != nil {
			t.Fatalf("GetEmbeddings failed: %v", err)
		}
		if len(embeddingVector) == 0 {
			t.Errorf("Expected non-empty embedding, got empty vector")
		}
	}
}
