package storage

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

const QdrantURL = "http://localhost:6333/collections/doc_chunks/points"

type Point struct {
	ID      string            `json:"id"`
	Vector  []float32         `json:"vector"`
	Payload map[string]string `json:"payload"`
}

type UpsertRequest struct {
	Points []Point `json:"points"`
}

// UpsertVector stores an embedding in Qdrant with document reference
func UpsertVector(docID, chunkID string, vector []float32, text string) error {
	point := Point{
		ID:     chunkID,
		Vector: vector,
		Payload: map[string]string{
			"doc_id": docID,
			"text":   text,
		},
	}

	reqBody, err := json.Marshal(UpsertRequest{Points: []Point{point}})
	if err != nil {
		return err
	}

	resp, err := http.Post(QdrantURL, "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.New("failed to upsert vector")
	}

	fmt.Printf("Successfully stored chunk %s for document %s\n", chunkID, docID)
	return nil
}
