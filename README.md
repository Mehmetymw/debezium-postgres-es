# Building a Real-Time Data Pipeline with Debezium, Kafka, and Elasticsearch

This project demonstrates how to build a robust real-time data pipeline that captures changes from PostgreSQL and streams them to Elasticsearch using Debezium, Kafka, and Go.

## Architecture Overview

The application follows a Domain-Driven Design (DDD) architecture with clean separation of concerns:

```
┌─────────────────┐     ┌─────────────┐     ┌─────────────────┐
│                 │     │             │     │                 │
│   PostgreSQL    │────▶│    Kafka    │────▶│  Elasticsearch  │
│                 │     │             │     │                 │
└─────────────────┘     └─────────────┘     └─────────────────┘
         │                                          ▲
         │                                          │
         │                                          │
         ▼                                          │
┌─────────────────────────────────────────────────────────────┐
│                                                             │
│                      Go Application                         │
│                                                             │
│  ┌─────────────┐    ┌─────────────┐    ┌─────────────┐     │
│  │             │    │             │    │             │     │
│  │   Domain    │◀──▶│ Application │◀──▶│ Interfaces  │     │
│  │             │    │             │    │             │     │
│  └─────────────┘    └─────────────┘    └─────────────┘     │
│         ▲                  ▲                  ▲            │
│         │                  │                  │            │
│         ▼                  ▼                  ▼            │
│  ┌─────────────┐    ┌─────────────┐    ┌─────────────┐     │
│  │             │    │             │    │             │     │
│  │Infrastructure│   │   Config    │    │    API      │     │
│  │             │    │             │    │             │     │
│  └─────────────┘    └─────────────┘    └─────────────┘     │
│                                                             │
└─────────────────────────────────────────────────────────────┘
```

## Key Features

- **Change Data Capture (CDC)**: Capture database changes in real-time using Debezium
- **Event Streaming**: Process and route events through Kafka
- **Search and Analytics**: Index data in Elasticsearch for powerful search capabilities
- **RESTful API**: Expose data through a clean API built with Fiber
- **Domain-Driven Design**: Well-structured codebase following DDD principles
- **Repository Pattern**: Clean separation between domain and data access layers
- **Configuration Management**: Flexible configuration using Viper

## Technology Stack

- **Go**: Core application language (v1.21+)
- **Fiber**: Fast HTTP web framework
- **GORM**: Powerful ORM for database operations
- **Viper**: Configuration management
- **PostgreSQL**: Source database
- **Debezium**: Change Data Capture platform
- **Kafka**: Event streaming platform
- **Elasticsearch**: Search and analytics engine
- **Docker**: Containerization

## Getting Started

### Prerequisites

- Go 1.21 or higher
- Docker and Docker Compose
- Git

### Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/mehmetymw/debezium-postgres-es.git
   cd debezium-postgres-es
   ```

2. Install dependencies:
   ```bash
   go mod tidy
   ```

3. Configure the application:
   Create a `config.yaml` file in the root directory or set environment variables.

### Running the Application

#### Using Go

```bash
go run main.go
```

#### Using Docker

```bash
docker-compose up -d
```

## API Endpoints

The application exposes the following RESTful API endpoints:

- `GET /api/orders` - Get all orders
- `GET /api/orders/:id` - Get a specific order
- `POST /api/orders` - Create a new order
- `PUT /api/orders/:id` - Update an existing order
- `DELETE /api/orders/:id` - Delete an order
- `GET /api/orders/status/:status` - Get orders by status
- `GET /health` - Health check endpoint

## Project Structure

The project follows a clean architecture with DDD principles:

```
app/
├── domain/                 # Domain layer (entities, repositories interfaces)
│   ├── entity/             # Domain entities
│   └── repository/         # Repository interfaces
├── application/            # Application layer (services, use cases)
│   └── service/            # Application services
├── infrastructure/         # Infrastructure layer (DB, external services)
│   └── persistence/        # Database related code
│       ├── models/         # Database models
│       ├── repository/     # Repository implementations
│       └── migrations/     # Database migrations
├── interfaces/             # Interface layer (API, CLI)
│   └── api/                # API related code
│       ├── handlers/       # HTTP handlers
│       └── routes/         # Route definitions
├── config/                 # Configuration
├── main.go                 # Application entry point
└── docker-compose.yml      # Docker Compose configuration
```

## How It Works

1. **Change Capture**: Debezium monitors PostgreSQL's transaction log for changes
2. **Event Publishing**: Changes are published to Kafka topics
3. **Event Consumption**: The application consumes events from Kafka
4. **Data Indexing**: Events are transformed and indexed in Elasticsearch
5. **Data Access**: The API provides access to the data from both PostgreSQL and Elasticsearch

## Configuration

The application uses Viper for configuration management. You can configure it using:

- Environment variables
- Configuration files (YAML, JSON, etc.)
- Command-line flags

Example configuration:

```yaml
postgres:
  host: localhost
  port: 5432
  user: postgres
  password: postgres
  dbname: inventory

elasticsearch:
  url: http://localhost:9200
  username: ""
  password: ""

server:
  port: 8080
```
