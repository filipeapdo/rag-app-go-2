# rag-app-go

## Big Picture Architecture

A system as a collection of independent services that communicate over a secure internal network.

**Key Components**

- Ingestion Service:

  - Responsibility: Handle incoming files (PDFs, TXT, CSV) and extract their text.
  - Process: Calls Ollama-based, OpenaAI-based embedding module (e.g., using the “snowflake-arctic-embed” model) and then sends the resulting vectors to the storage layer.
  - Implementation: A Go microservice with TDD and clean, modular code.

- Query & Retrieval (RAG) Service:

  - Responsibility: Accept user queries, retrieve the most relevant documents from a vector database, and relay these to downstream AI agents for processing.
  - Implementation: Again, built in Go for consistency, possibly with a REST and/or gRPC interface.

- AI Agent Manager / Orchestration:

  - Responsibility: Manage the pool of AI agents that perform automated tasks or further process queries.
  - Implementation: A controller that dispatches work based on the query results from the RAG pipeline.

- Storage Service:

  - Responsibility: Persist and index the vector embeddings and associated metadata.
  - Implementation: Already containerized via Docker Compose, with a persistent volume and internal networking.

- API Gateway / Orchestration Layer:
  - Responsibility: Expose the various services securely to external clients, handle authentication, request routing, and logging.
  - Implementation: Can be deployed as a separate container (or service) that ties the internal microservices together.

## Scaling & Future Expansion

- Local Testing & Production:

  - For development, use Docker Compose for a unified, containerized setup.
  - In production, migrate to an orchestrator like Kubernetes to handle scaling (horizontal pod autoscaling), rolling updates, and service discovery.

- Evolving Services:

  - Extend the Ingestion Module: Later, you can add API endpoints, web scrapers, or database integrations without modifying the core ingestion logic.
  - Modular AI Agents: As your needs grow, you can add more specialized AI agents that subscribe to the agent manager’s work queue.

- API Gateway:
  - Integrate an API gateway that supports routing, rate limiting, and security policies to manage external access and internal service communications.
