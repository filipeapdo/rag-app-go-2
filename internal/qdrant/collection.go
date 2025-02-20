package qdrant

import (
	"context"
	"log/slog"
	"time"

	"github.com/qdrant/go-client/qdrant"
)

// CollectionConfig holds the parameters for a collection
type CollectionConfig struct {
	Name       string
	VectorSize uint64
	Distance   qdrant.Distance
}

// CreateCollection creates a new Qdrant collection
func CreateCollection(client *qdrant.Client, config CollectionConfig) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	slog.Info("Creating collection in Qdrant...", "collection", config.Name)
	err := client.CreateCollection(ctx, &qdrant.CreateCollection{
		CollectionName: config.Name,
		VectorsConfig: qdrant.NewVectorsConfig(&qdrant.VectorParams{
			Size:     config.VectorSize,
			Distance: config.Distance,
		}),
	})
	if err != nil {
		slog.Error("Failed to create collection", "error", err)
		return err
	}

	slog.Info("✅ Collection created successfully", "collection", config.Name)
	return nil
}

// DeleteCollection removes a collection from Qdrant
func DeleteCollection(client *qdrant.Client, collection string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	slog.Info("Deleting collection in Qdrant...", "collection", collection)
	err := client.DeleteCollection(ctx, collection)
	if err != nil {
		slog.Error("Failed to delete collection", "error", err)
		return err
	}

	slog.Info("✅ Collection deleted successfully", "collection", collection)
	return nil
}

// ListCollections retrieves all collections in Qdrant
func ListCollections(client *qdrant.Client) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	slog.Info("Listing collections in Qdrant...")
	collections, err := client.ListCollections(ctx)
	if err != nil {
		slog.Error("Failed to list collections", "error", err)
		return err
	}

	if len(collections) == 0 {
		slog.Info("No Collections found")
		return nil
	}

	for _, collection := range collections {
		slog.Info("Collection found", "name", collection)
	}
	slog.Info("✅ All Collections listed successfully")
	return nil
}
