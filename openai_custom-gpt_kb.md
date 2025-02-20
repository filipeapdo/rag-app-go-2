<!-- ./README.md -->

# rag-app-go

In Summary, your project organization should embrace a microservices architecture with clearly delineated responsibilities:

- Ingestion Service: For file processing and embedding generation.
- Query/RAG Service: For retrieval and processing of queries.
- Agent Orchestration: To manage the pool of AI agents.
- Storage Service: For persisting embeddings in Qdrant.
- API Gateway: To secure and expose the services externally.

This modular structure allows you to develop, test, and deploy each part independently while ensuring that your system can scale horizontally as demand increases.

## Big Picture Architecture

Imagine your system as a collection of independent services that communicate over a secure internal network. Here’s a simplified diagram of how the services might interact:

```
                     +-----------------------+
                     |    API Gateway /      |
                     |   Orchestration Layer |
                     +-----------------------+
                              │
         ┌────────────────────┼───────────────────────┐
         │                    │                       │
  +----------------+  +-------------------+  +-------------------+
  | Ingestion      |  | Query & Retrieval |  | AI Agent Manager  |
  | Service        |  | Service (RAG)     |  | / Orchestration   |
  +----------------+  +-------------------+  +-------------------+
            │                 │
            └─────┬───────────┘
                  │
         +--------------------+
         |   Storage Service  |
         |    (Qdrant DB)     |
         +--------------------+
```

**Key Components**

- Ingestion Service:

  - Responsibility: Handle incoming files (PDFs, TXT, CSV) and extract their text.
  - Process: Calls your Ollama-based embedding module (e.g., using the “snowflake-arctic-embed” model) and then sends the resulting vectors to the storage layer.
  - Implementation: A Go microservice with TDD and clean, modular code.

- Query & Retrieval (RAG) Service:

  - Responsibility: Accept user queries, retrieve the most relevant documents from Qdrant, and relay these to downstream AI agents for processing.
  - Implementation: Again, built in Go for consistency, possibly with a REST and/or gRPC interface.

- AI Agent Manager / Orchestration:

  - Responsibility: Manage the pool of AI agents that perform automated tasks or further process queries.
  - Implementation: A controller that dispatches work based on the query results from the RAG pipeline.

- Storage Service (Qdrant):

  - Responsibility: Persist and index the vector embeddings and associated metadata.
  - Implementation: Already containerized via Docker Compose, with a persistent volume and internal networking.

- API Gateway / Orchestration Layer:
  - Responsibility: Expose the various services securely to external clients, handle authentication, request routing, and logging.
  - Implementation: Can be deployed as a separate container (or service) that ties the internal microservices together.

## Project Repository & Code Organization

**Monorepo with Microservices**
Organize your code in a single repository with clear subdirectories for each service. Current Organization:

```sh
rag-app-go-2/
├── bin/
│   └── ingestion
├── configs/ # Service configuration files
├── cmd/
│   ├── ingestion-service/ # Entry point for the Ingestion Service
│   │   └── main.go
│   ├── query-server/ # Entry point for the Query/RAG Service
│   └── agent-manager/ # Entry point for AI Agent Manager
├── deployments/ # Docker Compose and Kubernetes manifests
├── docker-compose.yml
├── go.mod
├── Makefile
├── pkg
│   ├── fileingestion
│   │   ├── fileingestion.go
│   │   └── fileingestion_test.go
│   ├── embedding/ # Wrapper around the Ollama API (for embeddings)
│   ├── storage/ # Integration with Qdrant
│   └── agents/ # AI agents orchestration logic
└── Makefile
└── README.md
```

**Key Best Practices**

- Microservice Boundaries: Each service is independently deployable and maintains its own domain logic.
- Containerization: Every service has its own Dockerfile. In development, Docker Compose ties them together. For production, plan on Kubernetes (or a similar orchestrator) to manage scaling and high availability.
- CI/CD & Testing: Implement automated tests (using TDD) for each service, and integrate CI/CD pipelines to build, test, and deploy your containers.
- Observability: Include logging, monitoring, and tracing in each service so you can track performance and debug issues quickly.

## Scaling & Future Expansion

- Local Testing & Production:

  - For development, use Docker Compose for a unified, containerized setup.
  - In production, migrate to an orchestrator like Kubernetes to handle scaling (horizontal pod autoscaling), rolling updates, and service discovery.

- Evolving Services:

  - Extend the Ingestion Module: Later, you can add API endpoints, web scrapers, or database integrations without modifying the core ingestion logic.
  - Modular AI Agents: As your needs grow, you can add more specialized AI agents that subscribe to the agent manager’s work queue.

- API Gateway:
  - Integrate an API gateway that supports routing, rate limiting, and security policies to manage external access and internal service communications.

<!-- ./docker-compose.yml -->

```yaml
networks:
  ai_network:
    driver: bridge

volumes:
  ollama_data:
  ollama_models:
  qdrant_storage:

services:
  ollama:
    image: ollama/ollama
    container_name: ollama
    ports:
      - "11434:11434" # Expose API externally
    restart: always
    networks:
      - ai_network
    volumes:
      - ollama_data:/root/.ollama
      - ollama_models:/root/.ollama/models # Persist models

  qdrant:
    image: qdrant/qdrant
    container_name: qdrant
    ports:
      - "6333:6333" # Expose API externally
    restart: always
    networks:
      - ai_network
    volumes:
      - qdrant_storage:/qdrant/storage

  # Future AI services (e.g., Vector DB) can be added here
```

<!-- ./Makefile -->

```Makefile
# Makefile for the Go project

# PHONY targets ensure that make always runs these commands
.PHONY: all build test fmt lint clean

# 'all' runs formatting, linting, tests, and then builds the binary.
all: fmt lint test build

# Build the binary from the main package.
build:
	@echo "Building the binary..."
	go build -o bin/ingestion-server ./cmd/ingestion-server

# Run all tests with verbose output.
test:
	@echo "Running tests..."
	go test -v ./...

# Format the code using go fmt.
fmt:
	@echo "Formatting code..."
	go fmt ./...

# Lint the code.
# Uncomment the lint line if you have golangci-lint installed.
lint:
	@echo "Running lint checks..."
	# golangci-lint run

# Clean up build artifacts.
clean:
	@echo "Cleaning up..."
	rm -rf bin
```

<!-- go.mod -->

```go
module github.com/filipeapdo/rag-app-go

go 1.24.0
```

<!-- ./pkg/fileingestion/fileingestion.go -->

```go
package fileingestion

import (
	"encoding/csv"
	"errors"
	"os"
	"path/filepath"
	"strings"
)

func ExtractText(filePath string) (string, error) {
	ext := strings.ToLower(filepath.Ext(filePath))
	switch ext {
	case ".txt":
		return readTextFile(filePath)
	case ".csv":
		return readCSVFile(filePath)
	default:
		return "", errors.New("unsupported file type")
	}
}

func readTextFile(filePath string) (string, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func readCSVFile(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return "", err
	}

	var sb strings.Builder
	for _, row := range records {
		sb.WriteString(strings.Join(row, ", "))
		sb.WriteString("\n")
	}

	return sb.String(), nil
}
```

<!-- ./pkg/fileingestion/fileingestion_test.go -->

```go
package fileingestion

import (
	"os"
	"strings"
	"testing"
)

func TestExtractText(t *testing.T) {
	// Test case: Unsupported file type
	t.Run("Unsupported file type", func(t *testing.T) {
		// Create a temporary UNSUPPORTED file.
		content := "dummy"
		tmpFile, err := os.CreateTemp("", "test-*.md")
		if err != nil {
			t.Fatal(err)
		}
		defer os.Remove(tmpFile.Name())
		os.WriteFile(tmpFile.Name(), []byte(content), 0644)

		_, err = ExtractText(tmpFile.Name())
		if err == nil {
			t.Error("Expected error for UNSUPPORTED file type, got nil")
		}
	})

	// Test case: Valid ".txt" extraction
	t.Run("Valid \".txt\" text extraction", func(t *testing.T) {
		// Create a temporary TXT file.
		content := "Hello, this is a test file."
		tmpFile, err := os.CreateTemp("", "test-*.txt")
		if err != nil {
			t.Fatal(err)
		}
		defer os.Remove(tmpFile.Name())
		os.WriteFile(tmpFile.Name(), []byte(content), 0644)

		text, err := ExtractText(tmpFile.Name())
		if err != nil {
			t.Fatalf("ExtractText returned error: %v", err)
		}
		contentCheck := "Hello, this is a test file."
		if strings.TrimSpace(text) != contentCheck {
			t.Errorf("Expected |%s|, got |%s|", contentCheck, text)
		}
	})

	// Test case: Valid ".csv" extraction
	t.Run("Valid \".csv\" extraction", func(t *testing.T) {
		// Create a temporary CSV file
		csvContent := "name,age\nFilipe,37\nDani,32"
		tmpFile, err := os.CreateTemp("", "test-*.csv")
		if err != nil {
			t.Fatal(err)
		}
		defer os.Remove(tmpFile.Name())
		os.WriteFile(tmpFile.Name(), []byte(csvContent), 0644)

		text, err := ExtractText(tmpFile.Name())
		if err != nil {
			t.Fatalf("ExtractText returned error: %v", err)
		}
		contentCheck := "Dani"
		if !strings.Contains(text, contentCheck) {
			t.Errorf("Expected extracted text to contain |%s|, extracted text is |%s|", contentCheck, text)
		}
	})
}
```

<!-- ./cmd/ingestion-service/main.go -->

```go
package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	fileingestion "github.com/filipeapdo/rag-app-go/pkg/fileingestion"
)

func main() {
	// Parse command-line flags
	filePath := flag.String("file", "", "Path to the file to be ingested")
	flag.Parse()

	if *filePath == "" {
		fmt.Println("Usage: ingestion-service -file=<path-to-file>")
		os.Exit(1)
	}

	// Process the file
	extractedText, err := fileingestion.ExtractText(*filePath)
	if err != nil {
		log.Fatalf("Error processing file: %v", err)
	}

	fmt.Println("File content extracted successfuly!")
	fmt.Print(extractedText)
	fmt.Println("---")
}
```
