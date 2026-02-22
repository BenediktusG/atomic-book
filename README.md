# 🎟️ AtomicBook: High-Concurrency Ticket Booking API

AtomicBook is a production-ready, highly concurrent backend REST API built in Go. It is designed to handle extreme traffic spikes (the "Thundering Herd" problem) and strictly prevent double-booking using advanced database locking mechanisms and caching strategies.

## 🏗️ Architecture & Tech Stack

- **Language**: Go (Golang)
- **Database**: PostgreSQL (Persistent storage & Transactional integrity)
- **Cache**: Redis (High-speed read layer)
- **Infrastructure**: Docker & Docker Compose (Multi-stage containerization)
- **CI/CD**: GitHub Actions (Automated build and test pipelines)
- **Load Testing**: K6 (Performance benchmarking)

## ✨ Key Engineering Features

- **Pessimistic Locking (Concurrency Control)**: Utilizes Postgres FOR UPDATE row-level locks within database transactions to ensure that even if 1,000 users try to book the exact same ticket at the exact same millisecond, the system will never double-book.
- **Redis Cache-Aside Pattern**: Protects the primary database from read-heavy traffic spikes by serving event details directly from Redis memory.
- **Graceful Shutdown**: Implements OS signal interception (SIGTERM/SIGINT) to safely drain active HTTP connections and close database connections without corrupting mid-flight user transactions.
- **Containerized Micro-Network**: Runs isolated within a Docker Compose private network utilizing internal DNS for secure service discovery between the API, Postgres, and Redis.

## 📊 Performance Benchmarks

- **Benchmarked using K6** simulating a massive traffic spike (500 concurrent virtual users ramping up and holding steady).
- **Target**: GET /event/{id}
- **Throughput**: ~14,200 Requests Per Second (RPS)
- **Success Rate**: 100% (0 dropped connections)
- **Bottleneck Resolved**: Offloading reads to Redis completely shielded the Postgres instance from CPU starvation.

## 🚀 Getting Started (Local Development)

Because the entire infrastructure is containerized, getting the project running on your local machine takes less than 2 minutes.

### 1. Prerequisites

```
Docker installed and running.

Git
```

### 2. Clone and Setup
```
git clone [https://github.com/BenediktusG/atomic-book.git](https://github.com/BenediktusG/atomic-book.git)
cd atomic-book
```

### 3. Environment Variables

Create a `.env` file in the root directory to configure the local containers:

```
# .env
DB_USER=postgres
DB_PASSWORD=secretpassword
DB_NAME=atomic_book
REDIS_PASSWORD=secretredis
```


### 4. Run the Cluster

Spin up the Go API, Postgres database, and Redis cache in detached mode:
```
docker-compose up --build -d
```

The API will now be live and listening at http://localhost:8080.

### 5. Verify the Infrastructure

Once the cluster is running, you can verify the services are communicating by hitting the cached endpoint:
```
curl http://localhost:8080/event/1
```

Note: The first request queries Postgres and caches the result in Redis. Subsequent requests are served directly from Redis memory.

### 6. Run the Load Tests (K6)

To prove the architecture's resilience, you can run the included load tests. Ensure you have K6 installed.

Test the "Thundering Herd" (Reads):
```
k6 run read_test.js
```


This ramps up to 500 concurrent users instantly refreshing the event page, testing the Redis cache performance.

Test the Transaction Locks (Writes):
```
k6 run write_test.js
```

This simulates massive concurrent booking attempts to verify the Postgres FOR UPDATE lock strictly prevents double-booking.

**Note**: Please ensure the event is exist before running the load tests. Feel free to change the url in the load test files to match your event id.

### 7. Stopping the Cluster

To gracefully shut down the environment and tear down the containers:

```
docker-compose down
```


(Check your terminal logs to see the API safely draining active HTTP requests before closing the database connections!)

## 🔌 API Endpoints

### 1. Get Event Details (Cached)

```
GET /event/{id}
```

Returns event information and available ticket count. Highly optimized for speed using the cache-aside pattern.

### 2. Book a Ticket (Transaction)

```
POST /book
```

Requires a JSON payload. Safely decrements the ticket count using row-level database locks.

```
{
  "event_id": 1,
  "user_id": "user_123",
  "tickets_requested": 2
}
```

Built with ❤️ to demonstrate modern backend engineering principles.