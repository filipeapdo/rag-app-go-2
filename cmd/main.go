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
	operation  = flag.String("op", "", "Operation to perform: create, delete, list")
	collection = flag.String("collection", "", "Name of the collection (required for create, delete)")
	vectorSize = flag.Uint64("vectorSize", 1536, "Vector size for the collection (only for create)")
	distance   = flag.String("distance", "cosine", "Distance metric: cosine, euclidean, dot")
)

func init() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage\n")
		flag.PrintDefaults()
	}
}

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
