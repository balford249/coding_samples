# File Evaluator

A distributed file evaluation system built with Go that processes file evaluation events from Kafka and persists results to PostgreSQL.

Note there are no unit tests, more time was spent upskilling on Docker so that changes could be validated by rebuilding images and passing in test 
GET and POST requests. See another go project sample to view unit test compentency. 

## Architecture

### Overview
The system follows a **registry pattern** with **dependency injection**, allowing multiple evaluator types to process events independently:

```
API Server ──POST──> /evaluate ──Kafka──> KafkaRunner ──> Registry ──> Evaluators ──> PostgreSQL
                         │                                                    ▲
                         └────────────GET /evaluate?id=X───────────────────┘
```

### Components

#### 1. **API Layer** (`api/`)

**Main Server** (`api/main.go`)
- HTTP server on configurable port
- Initializes Kafka producer and database connection
- Routes requests to handlers

**HttpHandler** (`api/handlers/handlers.go`)
- Processes `/evaluate` endpoint (POST & GET)
- Validates payloads and query parameters
- Orchestrates database and Kafka producer

**KafkaProducer** (`api/kafka/producer.go`)
- Wraps Confluent Kafka producer
- Sends evaluation events to Kafka topic
- Handles delivery confirmations with 5s flush timeout

#### 2. **Consumer Layer** (`consumer/`)
- `KafkaRunner`: Orchestrates Kafka message consumption
- `ConsumerType`: Interface for pluggable evaluator implementations
- Subscribes to Kafka topics and dispatches events to registered evaluators

#### 3. **Evaluators** (`evaluators/types/`)

**BaseFileEvaluator**
- Base struct with shared functionality
- Handles event deserialization from Kafka messages
- Provides `insertResult()` helper to persist evaluation results

**FileExistsConsumer**
- Evaluates if a file exists on the filesystem
- Uses `os.Stat()` to check file presence
- Returns `passed`/`failed` status

**IsTxtFileConsumer**
- Evaluates if a file is a text file
- (Implementation details in `isTxtFile.go`)

#### 4. **Database Layer** (`db/`)

**Store**
- PostgreSQL wrapper around `*sql.DB`
- **InitDB()**: Connects to PostgreSQL with retry logic (exponential backoff: 1s, 2s, 4s)
- **CreateNewEvent()**: Inserts new evaluation event
- **GetEvent()**: Retrieves event with all evaluator statuses (defaults to `pending`)
- **InsertResult()**: Stores evaluation results

**Database Connection**
```
user=appuser
dbname=eval
host=postgres
password=password
sslmode=disable
```

#### 5. **Registry Pattern** (`evaluators/main.go`)

Centralizes evaluator initialization with shared database store:

```go
func initRegistry(store *db.Store) map[string]consumer.ConsumerType {
    return map[string]consumer.ConsumerType{
        "FileExists": types.NewFileExistsConsumer(store),
        "IsTxtFile":  types.NewIsTxtFileConsumer(store),
    }
}
```

## API Endpoints

### POST /evaluate
Creates a new file evaluation request and publishes it to Kafka.

**Request Body:**
```json
{
  "path": "/path/to/file.txt"
}
```

**Response (201):**
```json
{
  "id": 12345
}
```

**Error Responses:**
- `400 Bad Request`: Missing or invalid `path` field
- `500 Internal Server Error`: Database or Kafka failure

**Flow:**
1. Validates `path` field
2. Creates event in database (`CreateNewEvent()`)
3. Publishes event to Kafka with event ID
4. Returns event ID to client

### GET /evaluate?id=X
Retrieves the evaluation status for a given event ID.

**Query Parameters:**
- `id` (required): Event ID (integer)

**Response (200):**
```json
[
  {
    "type": "FileExists",
    "status": "passed"
  },
  {
    "type": "IsTxtFile",
    "status": "pending"
  }
]
```

**Error Responses:**
- `400 Bad Request`: Missing or non-integer `id` parameter
- `500 Internal Server Error`: Event not found or database error

**Flow:**
1. Validates `id` query parameter
2. Fetches event from database (`GetEvent()`)
3. Returns array of evaluator statuses (defaults to `pending` if not yet evaluated)

## Data Model

### FileEvalEvent (Kafka Message)
```json
{
  "id": 12345,
  "path": "/path/to/file.txt"
}
```

### EvaluationResult (Database)
```json
{
  "type": "FileExists|IsTxtFile",
  "status": "passed|failed|pending"
}
```

## Configuration

### API Server (`config.json`)
```json
{
  "broker": "kafka:9092",
  "topic": "file-eval-events",
  "apiPort": 8080
}
```

### Consumer (`evaluators/config.json`)
```json
{
  "broker": "kafka:9092",
  "topic": "file-eval-events",
  "group": "file-evaluator-group",
  "type": "FileExists"
}
```

### Usage
```bash
# Start API server
cd api && go run main.go

# Start evaluator(s)
cd evaluators && go run main.go --config config.json
```

## Design Patterns

1. **Dependency Injection**: Database store and Kafka producer passed to handlers via constructors
2. **Registry Pattern**: Evaluator lookup by type string
3. **Interface-based Design**: `ConsumerType` interface for pluggability
4. **Base Class Pattern**: `BaseFileEvaluator` provides shared behavior
5. **Producer-Consumer Pattern**: Decoupled API and evaluation via Kafka

## Request Flow

```
Client → POST /evaluate → API Handler
              ↓
         CreateNewEvent() → DB
              ↓
         ProduceEvent() → Kafka Topic
              ↓
         KafkaRunner → Evaluator (FileExists/IsTxtFile)
              ↓
         InsertResult() → DB
              ↓
Client ← GET /evaluate?id=X ← GetEvent() ← DB
```

## Dependencies

- `confluent-kafka-go`: Kafka producer and consumer
- `github.com/lib/pq`: PostgreSQL driver (implied)
