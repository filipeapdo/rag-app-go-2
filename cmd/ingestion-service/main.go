package main

import (
	"fmt"
	"log"
	"os"

	"github.com/filipeapdo/rag-app-go/pkg/pipeline"
)

func main() {
	// Read list of document paths from CLI
	if len(os.Args) < 2 {
		fmt.Println("Usage: ingestion-service <file1> <file2> ...")
		os.Exit(1)
	}

	docPaths := os.Args[1:]

	// Process all documents
	err := pipeline.ProcessDocuments(docPaths)
	if err != nil {
		log.Fatalf("Error processing documents: %v", err)
	}

	fmt.Println("✅ All documents processed successfully!")
}
