# Building a Real-Time Data Pipeline with Debezium, Kafka, and Elasticsearch

This project demonstrates how to build a robust real-time data pipeline that captures changes from PostgreSQL and streams them to Elasticsearch using Debezium, Kafka, and Go.

## Architecture Overview

The application follows a Domain-Driven Design (DDD) architecture with clean separation of concerns:

```
PostgreSQL (Upsert) --> Debezium --> Kafka --> Kafka Connect --> Elasticsearch (GET)
         ▲                                                        |
         |                                                        v
      Go App ------------------------------------------------------

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

- **Go**: Core application language (v1.24.1+)
- **Fiber**: Fast HTTP web framework
- **GORM**: Powerful ORM for database operations
- **Viper**: Configuration management
- **PostgreSQL**: Source database
- **Debezium**: Change Data Capture platform
- **Kafka**: Event streaming platform
- **Elasticsearch**: Search and analytics engine
- **Docker**: Containerization

## Detailed Setup Guide

### Prerequisites

- Go 1.24.1 or higher
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

## Docker Compose Setup

Our `docker-compose.yml` file sets up the entire infrastructure needed for the CDC pipeline:

```yaml
version: '3'
services:
  # PostgreSQL database with Debezium support
  postgres:
    image: debezium/postgres:13
    ports:
      - "5432:5432"
    environment:
      - POSTGRES_USER=postgres
      - POSTGRES_PASSWORD=postgres
      - POSTGRES_DB=inventory
    volumes:
      - ./init.sql:/docker-entrypoint-initdb.d/init.sql

  # Zookeeper for Kafka coordination
  zookeeper:
    image: debezium/zookeeper:1.9
    ports:
      - "2181:2181"
      - "2888:2888"
      - "3888:3888"

  # Kafka message broker
  kafka:
    image: debezium/kafka:1.9
    ports:
      - "9092:9092"
    environment:
      - ZOOKEEPER_CONNECT=zookeeper:2181
    depends_on:
      - zookeeper

  # Kafka Connect with Debezium connectors
  connect:
    build:
      context: .
      dockerfile: Dockerfile
    ports:
      - "8083:8083"
    environment:
      - BOOTSTRAP_SERVERS=kafka:9092
      - GROUP_ID=1
      - CONFIG_STORAGE_TOPIC=connect_configs
      - OFFSET_STORAGE_TOPIC=connect_offsets
      - STATUS_STORAGE_TOPIC=connect_statuses
    depends_on:
      - kafka
      - postgres

  # Elasticsearch for data indexing and search
  elasticsearch:
    image: docker.elastic.co/elasticsearch/elasticsearch:7.17.0
    ports:
      - "9200:9200"
      - "9300:9300"
    environment:
      - discovery.type=single-node
      - xpack.security.enabled=false

  # Our Go application
  app:
    build:
      context: ./app
      dockerfile: Dockerfile
    ports:
      - "8080:8080"
    environment:
      - POSTGRES_HOST=postgres
      - POSTGRES_PORT=5432
      - POSTGRES_USER=postgres
      - POSTGRES_PASSWORD=postgres
      - POSTGRES_DB=inventory
      - ELASTICSEARCH_URL=http://elasticsearch:9200
    depends_on:
      - postgres
      - elasticsearch
      - connect
```

### Connector Configuration

#### PostgreSQL Source Connector

The PostgreSQL source connector (`postgres-source-all.json`) captures changes from our database:

```json
{
  "name": "postgres-source-all",
  "config": {
    "connector.class": "io.debezium.connector.postgresql.PostgresConnector",
    "database.hostname": "postgres",
    "database.port": "5432",
    "database.user": "postgres",
    "database.password": "postgres",
    "database.dbname": "inventory",
    "database.server.name": "dbserver1",
    "table.include.list": "public.orders",
    "plugin.name": "pgoutput",
    "transforms": "unwrap",
    "transforms.unwrap.type": "io.debezium.transforms.ExtractNewRecordState",
    "transforms.unwrap.drop.tombstones": "false",
    "transforms.unwrap.delete.handling.mode": "rewrite",
    "transforms.unwrap.add.fields": "op,table,lsn,source.ts_ms"
  }
}
```

#### Elasticsearch Sink Connector

The Elasticsearch sink connector (`elastic-sink-all.json`) streams the data to Elasticsearch:

```json
{
  "name": "elastic-sink-all",
  "config": {
    "connector.class": "io.confluent.connect.elasticsearch.ElasticsearchSinkConnector",
    "tasks.max": "1",
    "topics": "dbserver1.public.orders",
    "connection.url": "http://elasticsearch:9200",
    "transforms": "KeyToValue",
    "transforms.KeyToValue.type": "org.apache.kafka.connect.transforms.ValueToKey",
    "transforms.KeyToValue.fields": "id",
    "key.ignore": "false",
    "type.name": "_doc",
    "behavior.on.null.values": "delete",
    "schema.ignore": "true",
    "key.converter": "org.apache.kafka.connect.json.JsonConverter",
    "key.converter.schemas.enable": "false",
    "value.converter": "org.apache.kafka.connect.json.JsonConverter",
    "value.converter.schemas.enable": "false"
  }
}
```

## Automated Setup with steps.sh

We provide a convenient shell script (`steps.sh`) that automates the entire setup process. Here's what it does:

```bash
#!/bin/bash

# Step 1: Start all services with Docker Compose
docker-compose down && docker-compose up -d --build

# Step 2: Wait for services to start
sleep 10

# Step 3: Create orders table in PostgreSQL
docker exec -i debezium-postgres-1 psql -U postgres inventory -c "CREATE TABLE IF NOT EXISTS public.orders (
    id VARCHAR(255) PRIMARY KEY, 
    orderId VARCHAR(255), 
    customerId VARCHAR(255), 
    status VARCHAR(50)
);"

# Step 4: Insert sample data
docker exec -i debezium-postgres-1 psql -U postgres inventory -c "
INSERT INTO public.orders (id, orderId, customerId, status) VALUES 
('1', '101', '501', 'NEW'),
('2', '102', '502', 'PROCESSING'),
('3', '103', '503', 'COMPLETED'),
('4', '104', '504', 'SHIPPED'),
('5', '105', '505', 'DELIVERED'),
('6', '106', '506', 'CANCELLED'),
('7', '107', '507', 'RETURNED'),
('8', '108', '508', 'PENDING'),
('9', '109', '509', 'ON_HOLD'),
('10', '110', '510', 'BACKORDERED');"

# Step 5: Delete connectors if they exist
curl -X DELETE http://localhost:8083/connectors/postgres-source-all || echo "Connector not found"
curl -X DELETE http://localhost:8083/connectors/elastic-sink-all || echo "Connector not found"

# Step 6: Configure connectors
curl -X POST -H "Content-Type: application/json" --data @postgres-source-all.json http://localhost:8083/connectors
curl -X POST -H "Content-Type: application/json" --data @elastic-sink-all.json http://localhost:8083/connectors

# Step 7: Check Elasticsearch data
echo "Data in Elasticsearch after initial sync:"
curl -X GET "http://localhost:9200/dbserver1.public.orders/_search?pretty"

# Step 8: Check connector status
echo "PostgreSQL connector status:"
curl -X GET http://localhost:8083/connectors/postgres-source-all/status

echo "Elasticsearch connector status:"
curl -X GET http://localhost:8083/connectors/elastic-sink-all/status
```

### How to Use steps.sh

1. Make the script executable:
   ```bash
   chmod +x steps.sh
   ```

2. Run the script:
   ```bash
   ./steps.sh
   ```

The script performs the following actions:
1. Starts all services defined in docker-compose.yml
2. Creates the orders table in PostgreSQL
3. Inserts sample data into the orders table
4. Configures the Debezium PostgreSQL source connector
5. Configures the Elasticsearch sink connector
6. Verifies the data in Elasticsearch
7. Checks the status of both connectors

## Running the Application

### Using Go

```bash
go run main.go
```

### Using Docker

```bash
docker-compose up -d
```

## Testing the Pipeline

Once the pipeline is running, you can test it by making changes to the PostgreSQL database and observing how they propagate to Elasticsearch:

1. Insert a new order:
```bash
docker exec -i debezium-postgres-1 psql -U postgres inventory -c "
INSERT INTO public.orders (id, orderId, customerId, status) VALUES ('11', '111', '511', 'NEW');"
```

2. Update an existing order:
```bash
docker exec -i debezium-postgres-1 psql -U postgres inventory -c "
UPDATE public.orders SET status = 'SHIPPED' WHERE id = '2';"
```

3. Delete an order:
```bash
docker exec -i debezium-postgres-1 psql -U postgres inventory -c "
DELETE FROM public.orders WHERE id = '10';"
```

4. Check Elasticsearch to see the changes:
```bash
curl -X GET "http://localhost:9200/dbserver1.public.orders/_search?pretty"
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

## Troubleshooting

### Common Issues

1. **Connectors not starting**: Check Kafka Connect logs:
   ```bash
   docker logs debezium-connect-1
   ```

2. **Data not appearing in Elasticsearch**: Verify the connector status:
   ```bash
   curl -X GET "http://localhost:8083/connectors/elastic-sink-all/status"
   ```

3. **PostgreSQL connection issues**: Ensure PostgreSQL is running:
   ```bash
   docker ps | grep postgres
   ```

4. **Elasticsearch not indexing data**: Check Elasticsearch logs:
   ```bash
   docker logs debezium-elasticsearch-1
   ```

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

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the LICENSE file for details. 
