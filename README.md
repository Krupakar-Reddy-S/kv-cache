# Key-Value Cache Service

A high-performance in-memory key-value cache service implemented in Go using the Echo framework.

## Design Choices & Optimizations

- **Concurrency Management**:
  - Used Go's built-in `sync.RWMutex` for thread-safe operations
  - Implemented separate read/write locks to allow concurrent reads
  - Optimized for high-throughput concurrent access

- **Memory Management**:
  - Configurable maximum memory usage (default: 1.5GB for t3.small)
  - Memory-based eviction strategy removes oldest items when threshold is reached
  - Automatic cleanup routine prevents memory leaks

- **Time-Based Eviction**:
  - TTL (Time-To-Live) support with configurable item age
  - Lazy deletion during reads for expired items
  - Background cleanup process for expired entries

- **Performance Optimizations**:
  - Used Echo framework for its high-performance HTTP routing
  - Efficient map structure for O(1) lookups
  - Batch cleanup to reduce lock contention
  - Separate read/write locks for better concurrent performance

- **Resource Management**:
  - Configurable cleanup intervals
  - Memory usage monitoring and automatic eviction
  - Graceful shutdown support

## Building and Running

### Using Docker

1. Build the image:

```bash
docker build -t kv-cache .
```

2. Run the container:

```bash
docker run -p 7171:7171 kv-cache
```

### Local Development

1. Install Dependencies

```bash
go mod download
```

2. Run the application:

```bash
go run main.go
```

## API Usage

### PUT Operation

```bash
curl -X POST http://localhost:7171/put \
-H "Content-Type: application/json" \
-d '{"key": "example", "value": "test"}'
```

### GET Operation

```bash
curl http://localhost:7171/get?key=example
```

## API Response Format

### PUT Response
Success (HTTP 200):

```json
{
    "status": "OK",
    "message": "Key inserted/updated successfully."
}
```

### GET Response
Success (HTTP 200):

```json
{
    "status": "OK",
    "key": "example",
    "value": "test"
}
```

Not Found (HTTP 404):

```json
{
    "status": "ERROR",
    "message": "Key not found"
}
```


## Technical Specifications

### Cache Configuration
- Maximum memory usage: 1.5GB (configurable)
- TTL: 5 minutes (configurable)
- Cleanup interval: 1 minute (configurable)
- Maximum key and value size: 256 ASCII characters

### System Requirements
- Go 1.22 or higher
- Docker for containerized deployment
- Minimum 2GB RAM (AWS t3.small compatible)
- Port 7171 available

### Performance Characteristics
- O(1) average time complexity for operations
- Thread-safe for concurrent access
- Zero cache misses under normal memory conditions
- Automatic resource management
- Memory-efficient storage

## Testing
The service includes two test suites:
1. Basic functional tests (`test/basic_test.go`)
2. Advanced load tests (`test/advanced_test.go`) with:
   - Concurrent operations
   - Memory usage monitoring
   - Latency tracking
   - Cache hit/miss statistics