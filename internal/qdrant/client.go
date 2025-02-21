package qdrant

import (
	"log/slog"

	"github.com/qdrant/go-client/qdrant"
)

func NewClient(qdrantHost string, qdrantGrpcPort int) (*qdrant.Client, error) {
	// slog.Info("Initializing Qdrant client...")
	client, err := qdrant.NewClient(&qdrant.Config{
		Host: qdrantHost,
		Port: qdrantGrpcPort,
	})
	if err != nil {
		slog.Error("Failed to create Qdrant client", "error", err)
		return nil, err
	}

	// slog.Info("✅ Qdrant client successfully initialized")
	return client, nil
}
