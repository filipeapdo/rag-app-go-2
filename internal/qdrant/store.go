package qdrant

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
	"github.com/qdrant/go-client/qdrant"
)

// Metadata structure for document storage
type Metadata struct {
	Department   string `json:"department"`
	DocumentType string `json:"document_type"`
	ReferenceID  string `json:"reference_id"`
	CreatedAt    string `json:"created_at"`
}

// StoreVectorWithMetadata stores a vector with metadata and the actual text chunk
func StoreVectorWithMetadata(client *qdrant.Client, collection string, embedding []float32, chunkText string, metadata Metadata) error {
	ctx := context.Background()

	// Insert new record
	res, err := client.Upsert(ctx, &qdrant.UpsertPoints{
		CollectionName: collection,
		Points: []*qdrant.PointStruct{
			{
				Id:      qdrant.NewID(uuid.New().String()),
				Vectors: qdrant.NewVectors(embedding...), // Embedding data
				Payload: qdrant.NewValueMap(map[string]any{
					"department":    metadata.Department,
					"document_type": metadata.DocumentType,
					"reference_id":  metadata.ReferenceID,
					"created_at":    metadata.CreatedAt,
					"text":          chunkText, // Storing chunk text
				}),
			},
		},
	})
	if err != nil {
		slog.Error("Failed to store vector in Qdrant", "error", err)
		return err
	}
	fmt.Println(res)

	// slog.Info("✅ Successfully stored latest document", "reference_id", metadata.ReferenceID)
	return nil
}
