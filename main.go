package main

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/qdrant/go-client/qdrant"
)

const (
	QdrantURL          = "http://localhost:6333"
	CollectionName     = "test"
	OllamaSummaryURL   = "http://localhost:11434/api/generate"
	OllamaEmbeddingURL = "http://localhost:11434/api/embeddings"
)

// ----- Qdrant section -----
// TO-DO: use qdrant API (go)

// QdrantSearchRequest is the payload for searching.
type QdrantSearchRequest struct {
	Query []float32 `json:"query"`
	Limit int       `json:"limit"`
}

// QdrantSearchResult represents one search result.
type QdrantSearchResult struct {
	ID      int                    `json:"id"`
	Score   float32                `json:"score"`
	Payload map[string]interface{} `json:"payload"`
}

// Qdrant Query Result items
// Define the struct for the "id" field
type ID struct {
	UUID string `json:"uuid"`
}

// Define the struct for the "value" field inside the payload
type PayloadValue struct {
	StringValue string `json:"string_value"`
}

// Define the struct for the "payload" field
type Payload struct {
	Key   string       `json:"key"`
	Value PayloadValue `json:"value"`
}

// Define the struct for each item in the response
type QdrantItem struct {
	ID      ID      `json:"id"`
	Payload Payload `json:"payload"`
	Score   float32 `json:"score"`
	Version int     `json:"version,omitempty"` // Use omitempty if the field is optional
}

// Define the struct for the entire response (an array of QdrantItem)
type QdrantResponse []QdrantItem

// queryVectorDB searches the Qdrant collection for nearest points.
func queryVectorDB(queryEmbedding []float32, top uint64) ([]QdrantSearchResult, error) {
	client, err := qdrant.NewClient(&qdrant.Config{
		Host: "localhost",
		Port: 6334,
	})
	if err != nil {
		return nil, err
	}

	queryResult, err := client.Query(context.Background(), &qdrant.QueryPoints{
		CollectionName: "test",
		Query:          qdrant.NewQuery(queryEmbedding...),
		Limit:          &top,
		WithPayload:    qdrant.NewWithPayloadInclude("text"),
	})
	if err != nil {
		return nil, err
	}

	// fmt.Println(queryResult)

	// queryResultJson, err := json.Marshal(queryResult)
	// if err != nil {
	// 	return nil, err
	// }

	// Create a variable to hold the decoded data
	var response QdrantResponse

	// Unmarshal the JSON into the struct
	err = json.Unmarshal([]byte(string(queryResult)), &response)
	if err != nil {
		fmt.Println("Error decoding JSON:", err)
		return nil, err
	}

	for _, item := range response {
		fmt.Println()
		fmt.Println(item)
		fmt.Println()
		fmt.Println()
	}

	return nil, nil

	// searchReq := QdrantSearchRequest{
	// 	Query: queryEmbedding,
	// 	Limit: top,
	// }
	// body, err := json.Marshal(searchReq)
	// fmt.Println(body)
	// if err != nil {
	// 	return nil, err
	// }
	// url := fmt.Sprintf("%s/collections/%s/points/query", QdrantURL, CollectionName)
	// fmt.Println(url)
	// resp, err := http.Post(url, "application/json", bytes.NewBuffer(body))
	// if err != nil {
	// 	return nil, err
	// }
	// fmt.Println(resp.Body)
	// defer resp.Body.Close()
	// var res struct {
	// 	Result []QdrantSearchResult `json:"result"`
	// }
	// if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
	// 	return nil, err
	// }
	// return res.Result, nil
}

// ----- Summarization Service -----
// TO-DO: use ollama API (go)

// SummaryRequest is the payload to send to the summarization API.
type SummaryRequest struct {
	Prompt string `json:"prompt"`
	Model  string `json:"model"`
	Stream bool   `json:"stream"`
}

// SummaryResponse is the expected response structure.
type SummaryResponse struct {
	Summary string `json:"summary"`
}

// getSummary calls the Ollama summarization API with a given prompt.
func getSummary(prompt string) (string, error) {
	reqBody := SummaryRequest{
		Prompt: prompt,
		Model:  "qwen2.5:0.5b", // Adjust if necessary
		Stream: true,
	}
	body, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}
	resp, err := http.Post(OllamaSummaryURL, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	var sumResp SummaryResponse
	if err := json.NewDecoder(resp.Body).Decode(&sumResp); err != nil {
		return "", err
	}
	if sumResp.Summary == "" {
		return "", fmt.Errorf("empty summary returned")
	}
	return sumResp.Summary, nil
}

// ----- Embedding Service -----
// TO-DO: use ollama API (go)

// EmbeddingRequest is the payload to send to the embedding API.
type EmbeddingRequest struct {
	Prompt string `json:"prompt"`
	Model  string `json:"model"`
}

// EmbeddingResponse is the expected response structure.
type EmbeddingResponse struct {
	Embedding []float32 `json:"embedding"`
}

// getEmbedding calls the Ollama embedding API for the provided text.
func getEmbedding(text string) ([]float32, error) {
	reqBody := EmbeddingRequest{
		Prompt: text,
		Model:  "snowflake-arctic-embed:22m", // Adjust if necessary
	}
	body, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	resp, err := http.Post(OllamaEmbeddingURL, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var embResp EmbeddingResponse
	if err := json.NewDecoder(resp.Body).Decode(&embResp); err != nil {
		return nil, err
	}
	if len(embResp.Embedding) == 0 {
		return nil, fmt.Errorf("empty embedding returned")
	}
	return embResp.Embedding, nil
}

// ----- AI Agent -----

// generateSummary retrieves relevant document chunks based on the query and then summarizes them.
func generateSummary(query string) (string, error) {
	// Get the embedding for the query.
	queryEmbedding, err := getEmbedding(query)
	if err != nil {
		return "", fmt.Errorf("failed to embed query: %v", err)
	}
	// Query Qdrant to retrieve top relevant chunks.
	results, err := queryVectorDB(queryEmbedding, 5)
	if err != nil {
		return "", fmt.Errorf("vector DB query failed: %v", err)
	}
	if len(results) == 0 {
		return "No relevant content found to summarize.", nil
	}

	// Combine retrieved text from payloads.
	var texts []string
	for _, res := range results {
		if txt, ok := res.Payload["text"].(string); ok {
			texts = append(texts, txt)
		}
	}
	context := strings.Join(texts, "\n")
	// Build prompt for summarization.
	prompt := fmt.Sprintf("Summarize the following content concisely:\n\n%s\n\nUser Query: %s", context, query)
	summary, err := getSummary(prompt)
	if err != nil {
		return "", fmt.Errorf("error generating summary: %v", err)
	}
	return summary, nil
}

func main() {
	// Enter interactive loop.
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Enter queries to generate summaries (type 'exit' to quit):")
	for {
		fmt.Print("Query: ")
		query, _ := reader.ReadString('\n')
		query = strings.TrimSpace(query)
		if strings.ToLower(query) == "exit" {
			break
		}
		summary, err := generateSummary(query)
		if err != nil {
			fmt.Println("Error:", err)
		} else {
			fmt.Println("\nSummary:")
			fmt.Println(summary)
		}
		fmt.Println(strings.Repeat("-", 50))
	}
}
