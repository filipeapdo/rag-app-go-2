package pipeline

import (
	"fmt"
	"log"
	"path/filepath"
	"sync"

	"github.com/google/uuid"

	"github.com/filipeapdo/rag-app-go/pkg/embedding"
	"github.com/filipeapdo/rag-app-go/pkg/fileingestion"
	"github.com/filipeapdo/rag-app-go/pkg/storage"
)

// Document represents a document to be processed
type Document struct {
	ID     string   // Unique document ID
	Path   string   // File path
	Chunks []string // Text chunks
}

// ProcessDocuments handles chunking, embedding, and storing documents in Qdrant
func ProcessDocuments(docPaths []string) error {
	var wg sync.WaitGroup
	documents := []Document{}

	// Step 1: Read and chunk each document
	for _, path := range docPaths {
		wg.Add(1)

		go func(path string) {
			defer wg.Done()
			docID := uuid.New().String()

			// Extract text from file
			text, err := fileingestion.ExtractText(path)
			if err != nil {
				log.Printf("Failed to extract text from %s: %v", path, err)
				return
			}

			// Chunk the document text
			chunks := embedding.GetChunks(text, 512, 50) // 512 tokens per chunk with overlap
			if len(chunks) == 0 {
				log.Printf("No chunks generated for %s", path)
				return
			}

			documents = append(documents, Document{
				ID:     docID,
				Path:   path,
				Chunks: chunks,
			})

			fmt.Printf("Processed %d chunks for document: %s\n", len(chunks), filepath.Base(path))
		}(path)
	}

	wg.Wait() // Ensure all documents are processed

	// Step 2: Generate embeddings and store in Qdrant
	for _, doc := range documents {
		wg.Add(1)

		go func(doc Document) {
			defer wg.Done()
			for _, chunk := range doc.Chunks {
				// Generate embedding for chunk
				embeddingVector, err := embedding.GetEmbeddings(chunk)
				if err != nil {
					log.Printf("Error generating embedding for document %s: %v", doc.ID, err)
					continue
				}

				// Save to Qdrant
				chunkID := uuid.New().String()
				err = storage.UpsertVector(doc.ID, chunkID, embeddingVector, chunk)
				if err != nil {
					log.Printf("Failed to store chunk in Qdrant for doc %s: %v", doc.ID, err)
				}
			}
			fmt.Printf("Stored embeddings for document: %s\n", filepath.Base(doc.Path))
		}(doc)
	}

	wg.Wait() // Ensure all embeddings are stored
	return nil
}
