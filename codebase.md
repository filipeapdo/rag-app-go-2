# Project Organization

rag-app-go-2
├── cmd
│   └── main.go
├── go.mod
├── internal
│   └── qdrant
│   ├── client.go
│   └── collection.go
└── pkg
    └── logger
    └── logger.go

# Codebase

<!-- go.mod -->

```go
module github.com/filipeapdo/rag-app-go

go 1.24.0

require github.com/qdrant/go-client v1.13.0

require (
	golang.org/x/net v0.28.0 // indirect
	golang.org/x/sys v0.28.0 // indirect
	golang.org/x/text v0.21.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240827150818-7e3bb234dfed // indirect
	google.golang.org/grpc v1.66.0 // indirect
	google.golang.org/protobuf v1.34.2 // indirect
)
```

<!-- pkg/logger/logger.go -->

```go
package logger

import (
	"context"
	"log/slog"
	"os"
)

// InitLogger initializes the logger with the specified format.
// If format is "simple", it logs only the message; otherwise, it uses JSON format.
func InitLogger(format string) {
	var handler slog.Handler

	if format == "simple" {
		handler = &simpleHandler{}
	} else {
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			AddSource: true, // Include source file and line number
		})
	}

	logger := slog.New(handler)
	slog.SetDefault(logger)
}

// simpleHandler is a custom slog.Handler that logs only the message.
type simpleHandler struct{}

// Enabled reports whether the handler is enabled for the given level.
func (h *simpleHandler) Enabled(_ context.Context, _ slog.Level) bool {
	return true
}

// Handle logs only the message part of the record.
func (h *simpleHandler) Handle(_ context.Context, r slog.Record) error {
	_, err := os.Stdout.WriteString(r.Message + "\n")
	return err
}

// WithAttrs returns a new handler with the given attributes.
func (h *simpleHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return h // Attributes are ignored in this simple handler
}

// WithGroup returns a new handler with the given group name.
func (h *simpleHandler) WithGroup(name string) slog.Handler {
	return h // Groups are ignored in this simple handler
}
```

<!-- internal/qdrant/client.go -->

```go
package qdrant

import (
	"log/slog"

	"github.com/qdrant/go-client/qdrant"
)

func NewClient(qdrantHost string, qdrantGrpcPort int) (*qdrant.Client, error) {
	slog.Info("Initializing Qdrant client...")
	client, err := qdrant.NewClient(&qdrant.Config{
		Host: qdrantHost,
		Port: qdrantGrpcPort,
	})
	if err != nil {
		slog.Error("Failed to create Qdrant client", "error", err)
		return nil, err
	}

	slog.Info("✅ Qdrant client successfully initialized")
	return client, nil
}
```

<!-- internal/qdrant/collection.go -->

```go
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
```

package main

import (
"flag"
"fmt"
"log/slog"
"os"

    qdrantLib "github.com/qdrant/go-client/qdrant"

    "github.com/filipeapdo/rag-app-go/internal/qdrant"
    "github.com/filipeapdo/rag-app-go/pkg/logger"

)

var (
operation = flag.String("op", "", "Operation to perform: create, delete, list")
collection = flag.String("collection", "", "Name of the collection (required for create, delete)")
vectorSize = flag.Uint64("vectorSize", 1536, "Vector size for the collection (only for create)")
distance = flag.String("distance", "cosine", "Distance metric: cosine, euclidean, dot")
)

func init() {
flag.Usage = func() {
fmt.Fprintf(os.Stderr, "Usage\n")
flag.PrintDefaults()
}
}

<!-- cmd/main.go -->

```go
func main() {
	flag.Parse()

	logger.InitLogger("json")

	// Check for required flags
	if *operation == "" || (*operation != "list" && *collection == "") {
		fmt.Fprintln(os.Stderr, "Error: -op and -collection are required flags.")
		flag.Usage()
		os.Exit(1)
	}

	client, err := qdrant.NewClient("localhost", 6334)
	if err != nil {
		slog.Error("Application startup failed", "error", err)
		os.Exit(1)
	}

	var distanceMetric qdrantLib.Distance
	switch *distance {
	case "cosine":
		distanceMetric = qdrantLib.Distance_Cosine
	case "euclidean":
		distanceMetric = qdrantLib.Distance_Euclid
	case "dot":
		distanceMetric = qdrantLib.Distance_Dot
	default:
		slog.Error("Invalid distance metric. Use: cosine, euclidean, or dot")
		os.Exit(1)
	}

	switch *operation {
	case "create":
		if *collection == "" {
			slog.Error("Collection name is required for create operation")
			os.Exit(1)
		}
		config := qdrant.CollectionConfig{
			Name:       *collection,
			VectorSize: *vectorSize,
			Distance:   distanceMetric,
		}
		if err := qdrant.CreateCollection(client, config); err != nil {
			slog.Error("Failed to create collection", "error", err)
			os.Exit(1)
		}
	case "delete":
		if *collection == "" {
			slog.Error("Collection name is required for delete operation")
			os.Exit(1)
		}
		if err := qdrant.DeleteCollection(client, *collection); err != nil {
			slog.Error("Failed to delete collection", "error", err)
			os.Exit(1)
		}
	case "list":
		if err := qdrant.ListCollections(client); err != nil {
			slog.Error("Failed to list collections", "error", err)
			os.Exit(1)
		}
	default:
		slog.Error("Invalid operation. Use -op=create, -op=delete, or -op=list")
		os.Exit(1)
	}
}
```
