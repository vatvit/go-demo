---
code: GODEM-006
status: Proposed
dateCreated: 2026-01-09T20:03:11.183Z
type: Feature Enhancement
priority: High
dependsOn: GODEM-004, GODEM-005
---

# Add Kafka service with Go client and health check

## 1. Description

### Problem Statement
The application needs event streaming capability for asynchronous communication between services. Kafka provides reliable, scalable event streaming.

### Current State
- Docker Compose with app, MongoDB, Redis (after GODEM-004)
- No event streaming infrastructure

### Desired State
- Kafka and Zookeeper containers in Docker Compose
- Go Kafka client for producing/consuming events
- Health endpoint includes Kafka connectivity check

### Business/Technical Justification
- Event streaming enables decoupled, scalable architecture
- Kafka is industry standard for high-throughput messaging
- Essential for future microservices communication

## 2. Rationale

### Why This Change Is Necessary
- Real-time event processing requirements
- Async communication reduces coupling
- Enables event sourcing patterns

### What It Accomplishes
- Reliable event streaming infrastructure
- Go application can produce/consume events
- Health monitoring includes Kafka status

## 3. Solution Analysis

### Selected Approach

| Component | Choice | Rationale |
|-----------|--------|----------|
| Kafka | confluentinc/cp-kafka:7.5.0 | Well-maintained, includes Kraft mode |
| Zookeeper | confluentinc/cp-zookeeper:7.5.0 | Required for Kafka coordination |
| Go Client | github.com/IBM/sarama | Mature, full-featured, well-maintained |
| Alternative | github.com/segmentio/kafka-go | Simpler API, also good option |

### Alternatives Considered

| Option | Pros | Cons | Decision |
|--------|------|------|----------|
| RabbitMQ | Simpler, good for queues | Not event streaming | Rejected |
| NATS | Lightweight | Less ecosystem | Rejected |
| Redpanda | Kafka-compatible, simpler | Less common | Consider later |

## 4. Implementation Specification

### Docker Compose Additions

```yaml
services:
  # ... existing services ...

  zookeeper:
    image: confluentinc/cp-zookeeper:7.5.0
    environment:
      ZOOKEEPER_CLIENT_PORT: 2181
      ZOOKEEPER_TICK_TIME: 2000
    volumes:
      - zookeeper_data:/var/lib/zookeeper/data

  kafka:
    image: confluentinc/cp-kafka:7.5.0
    depends_on:
      - zookeeper
    ports:
      - "9092:9092"
    environment:
      KAFKA_BROKER_ID: 1
      KAFKA_ZOOKEEPER_CONNECT: zookeeper:2181
      KAFKA_ADVERTISED_LISTENERS: PLAINTEXT://kafka:29092,PLAINTEXT_HOST://localhost:9092
      KAFKA_LISTENER_SECURITY_PROTOCOL_MAP: PLAINTEXT:PLAINTEXT,PLAINTEXT_HOST:PLAINTEXT
      KAFKA_INTER_BROKER_LISTENER_NAME: PLAINTEXT
      KAFKA_OFFSETS_TOPIC_REPLICATION_FACTOR: 1
    volumes:
      - kafka_data:/var/lib/kafka/data

volumes:
  # ... existing volumes ...
  zookeeper_data:
  kafka_data:
```

### New Files

```
go-demo/
└── internal/
    └── kafka/
        ├── kafka.go         # Client wrapper
        ├── producer.go      # Event producer
        └── consumer.go      # Event consumer
```

### Kafka Client (`internal/kafka/kafka.go`)

```go
type Client struct {
    brokers  []string
    producer sarama.SyncProducer
}

func New(brokers []string) (*Client, error)
func (c *Client) Close() error
func (c *Client) Ping(ctx context.Context) error
func (c *Client) Producer() *Producer
func (c *Client) NewConsumer(groupID string) (*Consumer, error)
```

### Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `KAFKA_BROKERS` | `localhost:9092` | Comma-separated broker list |

### Health Endpoint Update

```json
{
  "status": "ok",
  "services": {
    "mongodb": "connected",
    "redis": "connected",
    "kafka": "connected"
  }
}
```

### Updated Makefile Targets

| Target | Description |
|--------|-------------|
| `make kafka-topics` | Create default topics |
| `make kafka-console` | Open Kafka console consumer |

## 5. Acceptance Criteria

- [ ] Kafka and Zookeeper containers start with `docker-compose up`
- [ ] Go app connects to Kafka successfully
- [ ] Can produce messages to a topic
- [ ] Can consume messages from a topic
- [ ] `/health` reports Kafka connectivity status
- [ ] Data persists in volumes
- [ ] BDD tests for Kafka health check
- [ ] Hot-reload still works