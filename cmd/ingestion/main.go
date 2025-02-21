package main

import (
	"context"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/ollama/ollama/api"

	"github.com/filipeapdo/rag-app-go/internal/qdrant"
	"github.com/filipeapdo/rag-app-go/pkg/logger"
)

const (
	ollamaRawUrl = "http://localhost:11434"
)

type VectorRecord struct {
	Prompt    string    `json:"prompt"`
	Embedding []float32 `json:"embedding"`
}

// TO-DO handle error!!!
func ChunkText(text string, chunkSize, overlap int) []string {
	chunks := []string{}
	for start := 0; start < len(text); start += chunkSize - overlap {
		end := start + chunkSize
		if end > len(text) {
			end = len(text)
		}
		chunks = append(chunks, text[start:end])
	}
	return chunks
}

func GetEmbeddingFromChunk(ctx context.Context, client *api.Client, doc string) ([]float32, error) {
	embeddingsModel := "snowflake-arctic-embed:22m"

	req := &api.EmbeddingRequest{
		Model:  embeddingsModel,
		Prompt: doc,
	}
	// get embeddings
	resp, err := client.Embeddings(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.Embedding, nil
}

func main() {
	logger.InitLogger("json")

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	qdrantClient, err := qdrant.NewClient("localhost", 6334)
	if err != nil {
		slog.Error("Application startup failed", "error", err)
		os.Exit(1)
	}

	// Read multiple text files from a directory
	// slog.Info("Reading document from Knowledge Base...")
	files, err := filepath.Glob("./cmd/ingestion/knowledge_base/*.txt")
	if err != nil {
		slog.Error("Error reading document directory", "error", err)
		os.Exit(1)
	}

	for _, filePath := range files {
		fileName := filepath.Base(filePath)
		parts := strings.Split(fileName, "_")
		if len(parts) < 3 {
			slog.Warn("Skipping file due to invalid naming convention", "fileName", fileName)
			continue
		}

		department := parts[0]
		documentType := parts[1]
		referenceID := strings.TrimSuffix(parts[2], filepath.Ext(parts[2]))

		content, err := os.ReadFile(filePath)
		if err != nil {
			slog.Error("Error reading file", "filePath", filePath, "error", err)
			continue
		}

		// slog.Info("Processing document", "filePath", filePath, "department", department)
		//
		// slog.Info("Chunking document into smaller pieces...")
		chunks := ChunkText(string(content), 512, 128)

		// Create embeddings from chuncks and save them in the VectorStore store
		// slog.Info("Create embedding from each chunck...")
		url, _ := url.Parse(ollamaRawUrl)
		ollamaClient := api.NewClient(url, http.DefaultClient)
		for idx, chunk := range chunks {
			// slog.Info("Creating embedding for chunck: ", "chuckNb", idx)

			embedding, err := GetEmbeddingFromChunk(ctx, ollamaClient, chunk)
			if err != nil {
				slog.Error("Error creating embedding", "chunkNb", idx)
			}

			metadata := qdrant.Metadata{
				Department:   department,
				DocumentType: documentType,
				ReferenceID:  referenceID,
				CreatedAt:    time.Now().Format(time.RFC3339),
			}

			// fmt.Println(embedding)

			err = qdrant.StoreVectorWithMetadata(qdrantClient, "test", embedding, chunk, metadata)
			if err != nil {
				slog.Error("Failed to store vector in Qdrant", "chunkIndex", idx, "error", err)
			}
		}
	}

	// slog.Info("Ingestion process completed successfully!")
}
