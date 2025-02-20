package embedding

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
)

const (
	OllamaURL      = "http://localhost:11434/api/embeddings"
	EmbeddingModel = "snowflake-arctic-embed:22m" // Specify the embedding model
)

type EmbedRequest struct {
	Model string `json:"model"`
	Text  string `json:"prompt"`
}

type EmbedResponse struct {
	Embedding []float32 `json:"embedding"`
}

// GetEmbeddings sends text to Ollama API for embedding, handling chunking
func GetEmbeddings(chunk string) ([]float32, error) {
	reqBody, err := json.Marshal(EmbedRequest{
		Model: EmbeddingModel,
		Text:  chunk,
	})
	if err != nil {
		return nil, err
	}

	resp, err := http.Post(OllamaURL, "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var embedResp EmbedResponse
	if err := json.NewDecoder(resp.Body).Decode(&embedResp); err != nil {
		return nil, err
	}

	if len(embedResp.Embedding) == 0 {
		return nil, errors.New("no embeddings received")
	}

	return embedResp.Embedding, nil
}

// GetChunks splits text into overlapping chunks
func GetChunks(text string, chunkSize, overlapSize int) []string {
	words := strings.Fields(text) // Tokenize input text into words
	var chunks []string

	// Iterate through the text, moving forward with overlap
	for i := 0; i < len(words); i += (chunkSize - overlapSize) {
		end := i + chunkSize

		// Ensure the last chunk captures remaining words properly
		if end > len(words) {
			end = len(words)
		}

		chunks = append(chunks, strings.Join(words[i:end], " "))

		// If we've reached the end of the text, stop
		if end == len(words) {
			break
		}
	}

	return chunks
}
